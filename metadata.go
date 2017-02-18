package kaa

import (
	"container/heap"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"os"
)

func getMyHead(data *MoveRequest) (Point, error) {
	for _, snake := range data.Snakes {
		if snake.Id == data.You && len(data.You) > 0 {
			return snake.Head(), nil
		}
	}
	return Point{}, errors.New("Could not get head")
}

func GetPointInDirection(p Point, direc string, data *MoveRequest) (*Point, error) {
	switch direc {
	case UP:
		return p.Up(data), nil
	case DOWN:
		return p.Down(data), nil
	case LEFT:
		return p.Left(data), nil
	case RIGHT:
		return p.Right(data), nil
	}
	return nil, errors.New(fmt.Sprintf("could not find direction %v", direc))
}

func getStaticData(data *MoveRequest, direc string) ([]*StaticData, error) {
	head, err := getMyHead(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get head of your snake"))
	}
	p, err := GetPointInDirection(head, direc, data)
	if err != nil {
		return nil, err
	}
	return graphSearch(p, data), nil
}

// returns an array of static data, the final static data is
// the maximum depth and the other depths, are defined in moves_to_depth
// in data.go. Will search from the point pos to the maximum depth provided
// a depth of any positive integer will max out at that integer and a depth of
// any negative integer will allow any negative number

func pushOntoPQ(
	p *Point,
	seen map[string]bool,
	pq *PriorityQueue,
	priority int) {

	if p != nil {
		if !seen[p.String()] {
			seen[p.String()] = true

			heap.Push(pq, &Item{
				value:    *p,
				priority: priority + 1,
			})
		}
	}
}

func graphSearch(pos *Point, data *MoveRequest) []*StaticData {
	ret := []*StaticData{}
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	priority := 1
	seen := make(map[string]bool)

	// make the first item priority 1 so that if statement
	// at the end of the loop is executed
	pushOntoPQ(pos, seen, &pq, 1)

	accumulator := &StaticData{}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if item.priority > priority {
			//	fmt.Printf("%v\n", item)
			//	for _, x := range ret {
			//		fmt.Printf("%v", x)
			//	}
			priority = item.priority
			ret = append(ret, accumulator)
			// copy accumulator
			tmp := *accumulator
			accumulator = &tmp
		}
		p := item.value
		//fmt.Printf("%v", p)

		pushOntoPQ(p.Up(data), seen, &pq, item.priority)
		pushOntoPQ(p.Down(data), seen, &pq, item.priority)
		pushOntoPQ(p.Left(data), seen, &pq, item.priority)
		pushOntoPQ(p.Right(data), seen, &pq, item.priority)

		if data.FoodMap[p.String()] {
			//fmt.Printf("food\n")
			accumulator.Food += 1
		}

		accumulator.Moves += 1

	}
	ret = append(ret, accumulator)
	// cut off extra accumulated value
	return ret[1:]
}

func FilterPossibleMoves(metaD map[string]*MetaData) []string {
	directions := []string{UP, DOWN, LEFT, RIGHT}
	ret := []string{}
	for _, direc := range directions {
		if len(metaD[direc].MovesAway) == 0 {
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

func FilterMovesVsSpace(metaD map[string]*MetaData, moves []string) []string {
	ret := []string{}
	for _, direc := range moves {
		if metaD[direc].MovesVsSpace > -2 {
			//fmt.Printf("%v\n", ret)
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		max := math.MinInt64
		for _, direc := range moves {
			if metaD[direc].MovesVsSpace < max {
				ret = []string{direc}
			} else if metaD[direc].MovesVsSpace < max {
				ret = append(ret, direc)
			}
		}
	}
	return ret
}

// not necessairily the best move but the move that we are going with
func bestMoves(metaD map[string]*MetaData) ([]string, error) {
	moves := FilterPossibleMoves(metaD)
	moves = FilterMovesVsSpace(metaD, moves)
	moves = ClosestFoodDirections(metaD, moves)
	return moves, nil
}

func numNeighbours(data *MoveRequest, direc string) (int, error) {
	head, err := getMyHead(data)
	if err != nil {
		return 0, err
	}
	_, err = GetPointInDirection(head, direc, data)
	if err != nil {
		return 0, err
	}
	return 0, nil
}
func FilterMinimizeSpace(data *MoveRequest, moves []string) string {

	return ""
}

func bestMove(data *MoveRequest) (string, error) {
	moves, err := bestMoves(data.MD)
	if err != nil {
		return "", err
	}
	if len(moves) == 0 {
		return "", errors.New("Unable to give you a good Move")
	}

	move := FilterMinimizeSpace(data, moves)
	return move, nil
}

func GetMovesVsSpace(data *MoveRequest, direc string) int {
	last, err := data.MD[direc].moveMax()
	if err != nil {
		return 0
	}
	excessMoves := last.Moves - data.MyLength - last.Food
	if excessMoves > 0 {
		// do something clever here to account for food
	}

	return excessMoves
}

func ClosestFood(data []*StaticData) int {
	for i, staticData := range data {
		//fmt.Printf("direction : %v\ndata for move %v: %#v\n", direc, i, staticData)
		if staticData.Food > 0 {
			// the first entry is always empty
			return i + 1
		}
	}
	return math.MaxInt64
}

func GenerateMetaData(data *MoveRequest) (*MoveRequest, error) {
	metaD := make(MoveMetaData)
	metaD["up"] = &MetaData{}
	metaD["down"] = &MetaData{}
	metaD["right"] = &MetaData{}
	metaD["left"] = &MetaData{}
	data.MD = metaD

	for direc, direcMD := range metaD {
		sd, err := getStaticData(data, direc)
		if err != nil {
			return data, err
		}

		direcMD.MovesAway = sd
		direcMD.ClosestFood = ClosestFood(sd)
		direcMD.MovesVsSpace = GetMovesVsSpace(data, direc)
	}
	return data, nil
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		//log.Errorf("%s environment variable not set.", k)
	}
	return v
}
