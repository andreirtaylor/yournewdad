package kaa

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"os"
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

func getMoves(data *MoveRequest, direc string) (int, error) {
	head, err := getMyHead(data)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Unable to get head of your snake"))
	}

	seen := make(map[string]bool)
	switch direc {
	case UP:
		return possibleMoves(head.Up(data), data, seen), nil
	case DOWN:
		return possibleMoves(head.Down(data), data, seen), nil
	case LEFT:
		return possibleMoves(head.Left(data), data, seen), nil
	case RIGHT:
		return possibleMoves(head.Right(data), data, seen), nil
	}
	return 0, errors.New(fmt.Sprintf("invalid direction", direc))
}

func possibleMoves(pos *Point, data *MoveRequest, seen map[string]bool) int {
	if pos == nil || seen[pos.String()] {
		return 0
	}
	seen[pos.String()] = true

	ret := 1
	ret += possibleMoves(pos.Up(data), data, seen)
	ret += possibleMoves(pos.Down(data), data, seen)
	ret += possibleMoves(pos.Left(data), data, seen)
	ret += possibleMoves(pos.Right(data), data, seen)

	return ret
}

func GenerateMetaData(data *MoveRequest) (map[string]*MetaData, error) {
	var err error
	metad := make(map[string]*MetaData)
	metad["up"] = &MetaData{}
	metad["down"] = &MetaData{}
	metad["right"] = &MetaData{}
	metad["left"] = &MetaData{}

	for direc, direcMD := range metad {
		direcMD.Moves, err = getMoves(data, direc)
		if err != nil {
			return metad, err
		}
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
