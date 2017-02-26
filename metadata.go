package kaa

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
)

// This file contains all of the functions which build
// the metadata in each direction

func getStaticData(data *MoveRequest, direc string) (*StaticData, error) {
	head, err := getMyHead(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to get head of your snake"))
	}
	p, err := GetPointInDirection(head, direc, data)
	if err != nil {
		return nil, err
	}

	// dont do this ye it doesnt work
	//MoveSnakeForward(data.MyIndex, data, direc)
	ret := graphSearch(p, data, direc)
	//MoveSnakeBackward(data.MyIndex, data)

	return ret, nil
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
func graphSearch(pos *Point, data *MoveRequest, currentDirec string) *StaticData {
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	if data.Direcs[currentDirec].KeySnakeData == nil {
		data.Direcs[currentDirec].KeySnakeData = make(map[int]*SnakeData)
	}

	seen := make(map[string]bool)
	// make the first item priority 1 so that if statement
	// at the end of the loop is executed
	pushOntoPQ(pos, seen, &pq, 1)
	if pq.Len() == 0 {
		return nil
	}

	priority := 1

	totalMoves := 0

	accumulator := &StaticData{}
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if item.priority > priority {
			priority = item.priority
		}
		p := item.value

		// push all directions on priority queue
		pushOntoPQ(p.Up(data), seen, &pq, item.priority)
		pushOntoPQ(p.Down(data), seen, &pq, item.priority)
		pushOntoPQ(p.Left(data), seen, &pq, item.priority)
		pushOntoPQ(p.Right(data), seen, &pq, item.priority)

		if data.FoodMap[p.String()] {
			//fmt.Printf("food\n")
			if data.Direcs[currentDirec].ClosestFood == 0 {
				data.Direcs[currentDirec].ClosestFood = priority - 1
			}
			accumulator.Food += 1
		}
		// add 1 to the moves in this direction in this generation
		accumulator.Moves += 1

		FindMinSnakePointInSurroundingArea(&p, data, data.Direcs[currentDirec].KeySnakeData)

		// add 1 to the total moves
		totalMoves += 1

	}
	data.Direcs[currentDirec].TotalMoves = totalMoves

	return accumulator
}

func GetMovesVsSpace(data *MoveRequest, direc string) int {
	if data.Direcs[direc] == nil {
		return 0
	}
	totalMoves := data.Direcs[direc].TotalMoves
	totalFood := data.Direcs[direc].TotalFood
	excessMoves := totalMoves - totalFood - data.Direcs[direc].minKeySnakePart().lengthLeft

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

	tightSpace := true
	for direc, direcMD := range data.Direcs {
		sd, err := getStaticData(data, direc)
		if err != nil {
			return err
		}

		//fmt.Printf("%#v\n%#v\n", sd, direc)
		if sd != nil {
			direcMD.MovesVsSpace = GetMovesVsSpace(data, direc)
			direcMD.TotalMoves = sd.Moves
			direcMD.TotalFood = sd.Food
			if direcMD.MovesVsSpace > 20 {
				tightSpace = false
			}
		}
	}
	data.tightSpace = tightSpace
	return nil
}
