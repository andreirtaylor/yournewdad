package kaa

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
)

// This file contains all of the functions which build
// the metadata in each direction

func getStaticData(data *MoveRequest, direc string) ([]*StaticData, error) {
	head, err := getMyHead(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get head of your snake"))
	}
	p, err := GetPointInDirection(&head, direc, data)
	if err != nil {
		return nil, err
	}
	return graphSearch(p, data, direc), nil
}

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

// returns an array of static data, Each value represents the
// things available in that number of moves.
//	i.e. the first value in the array are the things that"
//	     one move away, the second is 2 moves
// Will search from the point pos to the maximum depth provided
// a depth of any positive integer will max out at that integer and a depth of
// any negative integer will allow any negative number
func graphSearch(pos *Point, data *MoveRequest, currentDirec string) []*StaticData {
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
	totalMoves := 0

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

		// push all directions on priority queue
		pushOntoPQ(p.Up(data), seen, &pq, item.priority)
		pushOntoPQ(p.Down(data), seen, &pq, item.priority)
		pushOntoPQ(p.Left(data), seen, &pq, item.priority)
		pushOntoPQ(p.Right(data), seen, &pq, item.priority)

		if data.FoodMap[p.String()] {
			//fmt.Printf("food\n")
			accumulator.Food += 1
		}
		// add 1 to the moves in this direction in this generation
		accumulator.Moves += 1

		FindMinSnakePointInArea(&p, data, currentDirec)

		// add 1 to the total moves
		totalMoves += 1

	}
	data.Direcs[currentDirec].TotalMoves = totalMoves

	ret = append(ret, accumulator)
	// cut off extra accumulated value
	return ret[1:]
}

func ClosestFoodDirections(data *MoveRequest, moves []string) []string {
	directions := []string{}
	min := math.MaxInt64
	metaD := data.Direcs
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

func bestMove(data *MoveRequest) (string, error) {
	moves, err := bestMoves(data)
	if err != nil {
		return "", err
	}
	if len(moves) == 0 {
		return "", errors.New("Unable to give you a good Move")
	}

	move, err := FilterMinimizeSpace(data, moves)
	if err != nil {
		return "", err
	}
	return move, nil
}

func GetMovesVsSpace(data *MoveRequest, direc string) int {
	last, err := data.Direcs[direc].moveMax()
	if err != nil {
		return 0
	}
	excessMoves := last.Moves - last.Food - data.Direcs[direc].minKeySnakePart().lengthLeft
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

func GenerateMetaData(data *MoveRequest) error {
	data.init()
	data.Direcs = make(MoveMetaData)
	data.Direcs[UP] = &MetaDataDirec{}
	data.Direcs[DOWN] = &MetaDataDirec{}
	data.Direcs[LEFT] = &MetaDataDirec{}
	data.Direcs[RIGHT] = &MetaDataDirec{}

	for direc, direcMD := range data.Direcs {
		sd, err := getStaticData(data, direc)
		if err != nil {
			return err
		}

		direcMD.MovesAway = sd
		direcMD.ClosestFood = ClosestFood(sd)
		direcMD.MovesVsSpace = GetMovesVsSpace(data, direc)

	}
	return nil
}
