package kaa

import (
	"fmt"
	"math"
)

func keepFMTforGraph() {
	fmt.Printf("%v")
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
	t, err := getTail(data.MyIndex, data)
	if err != nil {
		return nil
	}

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

	accumulator.FoodHash = make(map[string]*FoodData)
	accumulator.MoveHash = make(map[string]*MinMaxData)
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

			//fmt.Printf("%v", p)
			if data.FoodMap[p.String()] {
				//fmt.Printf("food\n")
				if accumulator.ClosestFood == nil {
					accumulator.ClosestFood = p
				}
				accumulator.Food += 1
				foodptr := &FoodData{moves: item.moves, pnt: p}
				accumulator.FoodHash[foodptr.pnt.String()] = foodptr
				if me {
					accumulator.sortedFood = append(accumulator.sortedFood, foodptr)
				}
			}
			// add 1 to the moves in this direction in this generation
			if me {
				FindMinSnakePointInSurroundingArea(p, data, ksd)
				if p.isNeighbour(t) && item.moves > 1 {
					accumulator.SeeTail = true
				}
			}
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
