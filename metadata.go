package kaa

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// remember to defer db.close
func getDB(req *http.Request) (*sql.DB, error) {
	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := mustGetenv("CLOUDSQL_PASSWORD")

	sqlStr := fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName)

	if appengine.IsDevAppServer() {
		sqlStr = "root@/Kaa" // dev server has no password baby
	}

	db, err := sql.Open(
		"mysql",
		sqlStr)

	return db, err
}

func saveGame(g *GameStartRequest, req *http.Request) {
	ctx := appengine.NewContext(req)

	db, err := getDB(req)
	if err != nil {
		log.Errorf(ctx, "Could not get DB %v", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(
		`INSERT INTO Games(GameId, Width, Height)
			VALUES(?,?,?)`)

	if err != nil {
		log.Errorf(ctx, "Unable to prepare game saving statement: %v", err)
		return
	}
	_, err = stmt.Exec(g.GameId, g.Width, g.Height)
	if err != nil {
		log.Errorf(ctx, "Error executing game save statement: %v", err)
		return
	}
}

func SaveMove(bo *MoveRequest, req *http.Request) {
	ctx := appengine.NewContext(req)

	db, err := getDB(req)
	if err != nil {
		log.Errorf(ctx, "Could not get DB %v", err)
		return
	}
	defer db.Close()

	var id int
	err = db.QueryRow(
		`SELECT id FROM Games
		WHERE gameid = ?`, bo.GameId).Scan(&id)

	if err != nil {
		log.Errorf(ctx, "Could not retrieve game Id %v", err)
		return
	}

	stmt, err := db.Prepare(
		`INSERT INTO MoveReq(g_id, turn)
			VALUES(?,?)`)

	if err != nil {
		log.Errorf(ctx, "Unable to prepare move saving statement: %v", err)
		return
	}

	res, err := stmt.Exec(id, bo.Turn)
	if err != nil {
		log.Errorf(ctx, "Error executing move save statement: %v", err)
		return
	}

	m_id, err := res.LastInsertId()
	if err != nil {
		log.Errorf(ctx, "Error getting last move statement: %v", err)
		return
	}

	for _, snake := range bo.Snakes {
		// store snakes
		stmt, err := db.Prepare(
			`INSERT INTO Snakes(name_, health, m_id, len)
				VALUES(?,?,?,?)`)

		if err != nil {
			log.Errorf(ctx, "Unable to prepare move saving statement: %v", err)
			return
		}

		_, err = stmt.Exec(snake.Name[:15], snake.HealthPoints, m_id, len(snake.Coords))
		if err != nil {
			log.Errorf(ctx, "Error executing move save statement: %v", err)
			return
		}

		for _, coords := range snake.Coords {
			stmt, err := db.Prepare(
				`INSERT INTO Body(x,y,m_id)
				VALUES(?,?,?)`)

			if err != nil {
				log.Errorf(ctx, "preparing coordinate save:\n%v", err)
				return
			}

			_, err = stmt.Exec(coords.X, coords.Y, m_id)
			if err != nil {
				log.Errorf(ctx, "executing coordinate save:\n%v", err)
				return
			}
		}
	}

	for _, food := range bo.Food {
		stmt, err := db.Prepare(
			`INSERT INTO Food(x,y,m_id)
			VALUES(?,?,?)`)

		if err != nil {
			log.Errorf(ctx, "preparing food save:\n%v", err)
			return
		}

		_, err = stmt.Exec(food.X, food.Y, m_id)
		if err != nil {
			log.Errorf(ctx, "executing food save:\n%v", err)
			return
		}
	}

}

func getMyHead(data *MoveRequest) (Point, error) {
	for _, snake := range data.Snakes {
		if snake.Id == data.You {
			return snake.Head(), nil
		}
	}
	return Point{}, errors.New("Could not get head")
}

func getStaticData(data *MoveRequest, direc string) ([]*StaticData, error) {
	head, err := getMyHead(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get head of your snake"))
	}

	switch direc {
	case UP:
		return graphSearch(head.Up(data), data), nil
	case DOWN:
		return graphSearch(head.Down(data), data), nil
	case LEFT:
		return graphSearch(head.Left(data), data), nil
	case RIGHT:
		return graphSearch(head.Right(data), data), nil
	}
	return nil, errors.New(fmt.Sprintf("invalid direction", direc))
}

func graphSearchRec(pos *Point, data *MoveRequest, seen map[string]bool, depth, max int) *StaticData {
	if depth == max || pos == nil || seen[pos.String()] {
		return &StaticData{}
	}
	seen[pos.String()] = true

	ret := &StaticData{}
	ret_up := graphSearchRec(pos.Up(data), data, seen, depth+1, max)
	ret_down := graphSearchRec(pos.Down(data), data, seen, depth+1, max)
	ret_left := graphSearchRec(pos.Left(data), data, seen, depth+1, max)
	ret_right := graphSearchRec(pos.Right(data), data, seen, depth+1, max)

	ret.Moves = ret_up.Moves + ret_down.Moves + ret_left.Moves + ret_right.Moves + 1

	ret.Food = ret_up.Food + ret_down.Food + ret_left.Food + ret_right.Food
	if data.FoodMap[pos.String()] {
		ret.Food += 1
	}
	return ret
}

// returns an array of static data, the final static data is
// the maximum depth and the other depths, are defined in moves_to_depth
// in data.go. Will search from the point pos to the maximum depth provided
// a depth of any positive integer will max out at that integer and a depth of
// any negative integer will allow any negative number
func graphSearch(pos *Point, data *MoveRequest) []*StaticData {
	ret := []*StaticData{}

	for _, depth := range moves_to_depth {
		seen := make(map[string]bool)
		sd := graphSearchRec(pos, data, seen, 0, depth)
		ret = append(ret, sd)
	}
	return ret
}

func ClosestFood(data []*StaticData) int {
	for i, staticData := range data {
		//fmt.Printf("direction : %v\ndata for move %v: %#v\n", direc, i, staticData)
		if staticData.Food > 0 {
			return i + 1
		}
	}
	return math.MaxInt64
}

func FilterPossibleMoves(metaD map[string]*MetaData) []string {
	directions := []string{UP, DOWN, LEFT, RIGHT}
	ret := []string{}
	for _, direc := range directions {
		if metaD[direc].moveMax().Moves > 0 {
			ret = append(ret, direc)
		}
	}
	return directions
}

func ClosestFoodDirections(metaD map[string]*MetaData, moves []string) []string {
	directions := []string{}
	min := math.MaxInt64
	for _, direc := range moves {
		if metaD[direc].ClosestFood < min {
			directions = []string{}
			directions = append(directions, direc)
			min = metaD[direc].ClosestFood
		} else if metaD[direc].ClosestFood == min {
			directions = append(directions, direc)
		}
	}
	return directions
}

// not necessairily the best move but the move that we are going with
func bestMoves(metaD map[string]*MetaData) []string {
	moves := FilterPossibleMoves(metaD)
	rand.Seed(time.Now().Unix()) // initialize global pseudorandom generator
	moves = ClosestFoodDirections(metaD, moves)
	return moves
}

func bestMove(metaD map[string]*MetaData) string {
	moves := bestMoves(metaD)
	return moves[rand.Intn(len(moves))]
}

func GenerateMetaData(data *MoveRequest) (map[string]*MetaData, error) {
	metad := make(map[string]*MetaData)
	metad["up"] = &MetaData{}
	metad["down"] = &MetaData{}
	metad["right"] = &MetaData{}
	metad["left"] = &MetaData{}

	for direc, direcMD := range metad {
		sd, err := getStaticData(data, direc)
		if err != nil {
			return metad, err
		}

		direcMD.MovesAway = sd
		direcMD.ClosestFood = ClosestFood(sd)
	}
	return metad, nil
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		//log.Errorf("%s environment variable not set.", k)
	}
	return v
}
