package kaa

import (
	"container/heap"
	"fmt"
	"math"
)

func keepFMTforGraph() {
	fmt.Printf("%v")
}

func fullStatsPnt(pos *Point, data *MoveRequest) *MetaDataDirec {
	return quickStats(pos, data, math.MaxInt64, false)
}

func fullStatsMe(pos *Point, data *MoveRequest) *MetaDataDirec {
	return quickStats(pos, data, math.MaxInt64, true)
}

// me is a boolean that basicly
// a function that is to be used on other points
func quickStats(pos *Point, data *MoveRequest, depth int, me bool) *MetaDataDirec {
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

	t, err := getTail(data.MyIndex, data)
	if err != nil {
		return nil
	}

	currDepth := 1

	accumulator := &MetaDataDirec{}
	accumulator.FoodHash = make(map[string]*FoodData)
	accumulator.MoveHash = make(map[string]*MinMaxData)

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
			if accumulator.ClosestFood == nil && me {
				accumulator.ClosestFood = &p
			}
			accumulator.Food += 1
			foodptr := &FoodData{moves: currDepth - 2, pnt: &p}
			accumulator.FoodHash[foodptr.pnt.String()] = foodptr
			if me {
				accumulator.sortedFood = append(accumulator.sortedFood, foodptr)
			}
		}
		// add 1 to the moves in this direction in this generation
		accumulator.Moves += 1
		if me {
			FindMinSnakePointInSurroundingArea(&p, data, ksd)
			if p.isNeighbour(t) && currDepth > 2 {
				accumulator.SeeTail = true
			}
		}

		// the snake head shouldnt be in the hash
		if currDepth > 2 {
			accumulator.MoveHash[p.String()] = &MinMaxData{moves: currDepth}
		}
	}

	accumulator.KeySnakeData = ksd
	return accumulator
}
