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

func AppendIfMissing(slice []int, i int) ([]int, bool) {
	for _, ele := range slice {
		if ele == i {
			return slice, false
		}
	}
	slice = append(slice, i)
	return slice, true
}

// me is a boolean that basicly
// a function that is to be used on other points
func quickStats2(data *MoveRequest, direc string) *MetaDataDirec {
	// generated the hazards without the hazards around the other snakes
	data.GenHazards(data, false)
	myHead := data.Snakes[data.MyIndex].Head()
	if direc != "" {
		myHeadtmp, err := GetPointInDirection(myHead, direc, data)
		if err != nil {
			return nil
		}
		myHead = myHeadtmp
		if myHead == nil {
			return nil
		}
		if myHead != nil {
			data.Hazards[myHead.String()] = true
		}
	}

	q := Queue{}

	ksd := make(map[int]*SnakeData)
	// make the first item priority 1 so that if statement
	// at the end of the loop is executed
	for i, snake := range data.Snakes {
		head := &snake.Coords[0]
		moves := -1
		if i == data.MyIndex {
			moves = 0
			head = myHead
		}
		mmd := &MinMaxData{
			moves:    moves,
			snakeIds: []int{i},
			tie:      false,
			pnt:      head,
		}
		//ret[head.y][head.X] = mmd
		q.Push(mmd)
	}

	//t, err := getTail(data.MyIndex, data)
	//if err != nil {
	//	return nil
	//}

	accumulator := &MetaDataDirec{}
	//accumulator.FoodHash = make(map[string]*FoodData)
	accumulator.MoveHash = make(map[string]*MinMaxData)
	accumulator.minMaxArr = make(MMArray, data.Height)
	for i := range accumulator.minMaxArr {
		accumulator.minMaxArr[i] = make([]MinMaxData, data.Width)
		for j := range accumulator.minMaxArr[i] {
			accumulator.minMaxArr[i][j].moves = math.MaxInt64
		}

	}
	for q.Len() > 0 {
		item := q.Pop()

		p := item.pnt

		if p == nil {
			continue
		}

		boardState := accumulator.minMaxArr[item.pnt.Y][item.pnt.X]
		//fmt.Printf("%v %v %v\n ", p, boardState.moves, item.moves)
		if boardState.moves == item.moves {
			for _, id := range item.snakeIds {
				if snakeIds, ok := AppendIfMissing(accumulator.minMaxArr[item.pnt.Y][item.pnt.X].snakeIds, id); ok {
					accumulator.minMaxArr[item.pnt.Y][item.pnt.X].snakeIds = snakeIds
					accumulator.minMaxArr[item.pnt.Y][item.pnt.X].tie = true
				}

			}
			continue
		}
		if len(boardState.snakeIds) != 0 {
			continue
		}
		q.Push(
			&MinMaxData{
				moves:    item.moves + 1,
				snakeIds: item.snakeIds,
				pnt:      p.Up(data),
			})

		q.Push(
			&MinMaxData{
				moves:    item.moves + 1,
				snakeIds: item.snakeIds,
				pnt:      p.Right(data),
			})
		q.Push(
			&MinMaxData{
				moves:    item.moves + 1,
				snakeIds: item.snakeIds,
				pnt:      p.Left(data),
			})
		q.Push(
			&MinMaxData{
				moves:    item.moves + 1,
				snakeIds: item.snakeIds,
				pnt:      p.Down(data),
			})
		me := false
		for _, sid := range item.snakeIds {
			if sid == data.MyIndex {
				me = true
			}
		}
		if me {
			accumulator.MoveHash[item.pnt.String()] = item
			accumulator.Moves += 1
		}

		if (!me && item.moves >= 0) || (me && item.moves > 0) {
			accumulator.minMaxArr[item.pnt.Y][item.pnt.X].moves = item.moves
			accumulator.minMaxArr[item.pnt.Y][item.pnt.X].pnt = item.pnt
			accumulator.minMaxArr[item.pnt.Y][item.pnt.X].snakeIds = item.snakeIds
		}

	}

	accumulator.KeySnakeData = ksd
	return accumulator
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
