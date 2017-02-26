package kaa

import (
	"container/heap"
	"fmt"
	"math"
)

func keepFMTforGraph() {
	fmt.Printf("%v")
}

func fullStats(pos *Point, data *MoveRequest) *StaticData {
	return quickStats(pos, data, math.MaxInt64)
}

// a function that is to be used on other points
func quickStats(pos *Point, data *MoveRequest, depth int) *StaticData {
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	seen := make(map[string]bool)
	ksd := make(map[int]*SnakeData)
	// make the first item priority 1 so that if statement
	// at the end of the loop is executed
	pushOntoPQ(pos, seen, &pq, 1)
	if pq.Len() == 0 {
		return nil
	}

	currDepth := 1

	accumulator := &StaticData{}
	accumulator.FoodHash = make(map[string]*FoodData)

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if item.priority > currDepth {
			currDepth = item.priority
			if currDepth > depth {
				break
			}
		}
		p := item.value

		// push all directions on priority queue
		pushOntoPQ(p.Up(data), seen, &pq, item.priority)
		pushOntoPQ(p.Down(data), seen, &pq, item.priority)
		pushOntoPQ(p.Left(data), seen, &pq, item.priority)
		pushOntoPQ(p.Right(data), seen, &pq, item.priority)

		//fmt.Printf("%v", p)
		if data.FoodMap[p.String()] {
			//fmt.Printf("food\n")
			if accumulator.ClosestFood == nil {
				accumulator.ClosestFood = &p
			}
			accumulator.Food += 1
			foodptr := &FoodData{moves: currDepth - 2, pnt: &p}
			accumulator.FoodHash[foodptr.pnt.String()] = foodptr
		}
		// add 1 to the moves in this direction in this generation
		accumulator.Moves += 1
		FindMinSnakePointInSurroundingArea(&p, data, ksd)
	}

	accumulator.KeySnakeData = ksd
	return accumulator
}
