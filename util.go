package main

import (
	"bytes"
	"errors"
	"fmt"
)

func swap(arr []*Snake, a, b int) {
	arr[b], arr[a] = arr[a], arr[b]
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MinMax(data *MoveRequest, direc string) {
	// generated the hazards without the hazards around the other snakes

	data.GenHazards(data, false)
	myHead := data.Snakes[data.MyIndex].Head()
	if direc != "" {
		myHeadtmp, err := GetPointInDirection(myHead, direc, data)
		if err != nil {
			return
		}
		myHead = myHeadtmp
		if myHead == nil {
			return
		}
		if myHead != nil {
			data.Hazards[myHead.String()] = true
		}
	}

	ret := quickStats2(data, direc)
	//stats := fullStatsMe(myHead, data)
	if direc == "" && len(ret.sortedFood) > 0 {
		data.DistToFood = ret.sortedFood[0].moves
	}
	if direc != "" {
		data.Direcs[direc].ClosestFood = ret.ClosestFood
		data.Direcs[direc].Food = ret.Food
		data.Direcs[direc].Moves = ret.Moves
		data.Direcs[direc].SeeTail = ret.SeeTail
		data.Direcs[direc].KeySnakeData = ret.KeySnakeData
		data.Direcs[direc].FoodHash = ret.FoodHash
		data.Direcs[direc].sortedFood = ret.sortedFood
		data.Direcs[direc].MoveHash = ret.MoveHash
	}

	if direc != "" {
		data.Direcs[direc].minMaxArr = ret.minMaxArr
	} else if direc == "" {
		data.minMaxArr = ret.minMaxArr
	}

}

func GenMinMaxStats(arr MMArray) MinMaxMetaData {
	ret := MinMaxMetaData{}
	ret.movesHash = make(map[string]int)
	ret.tiesHash = make(map[string][]int)
	ret.snakes = make(map[int]MinMaxSnakeMD)
	for i := range arr {
		for j := range arr[i] {
			p := &Point{X: i, Y: j}
			ids := arr[i][j].snakeIds

			if arr[i][j].tie {
				ret.tiesHash[p.String()] = ids
			}

			for _, id := range ids {
				s, ok := ret.snakes[id]
				if !ok {
					ret.snakes[id] = MinMaxSnakeMD{}
				}
				if arr[i][j].tie {
					s.ties++
				} else {
					s.moves++
					ret.movesHash[p.String()] = id
				}
				ret.snakes[id] = s
			}
		}
	}
	return ret
}

func stringAllMinMAX(data *MoveRequest) string {
	var buffer bytes.Buffer
	buffer.WriteString("\n board\n ")
	buffer.WriteString(data.minMaxArr.String())
	for _, direc := range []string{UP, RIGHT, DOWN, LEFT} {
		if data.Direcs[direc].minMaxArr != nil {
			buffer.WriteString(fmt.Sprintf("%v\n", direc))
			buffer.WriteString(data.Direcs[direc].minMaxArr.String())
		}
	}
	return buffer.String()
}

func findGuaranteedClosestFood(data *MoveRequest, direc string) *FoodData {
	for _, food := range data.Direcs[direc].sortedFood {
		for _, id := range data.minMaxArr[food.pnt.Y][food.pnt.X].snakeIds {
			if id == data.MyIndex {
				return food
			}
		}
	}
	return nil
}

func getTail(ind int, data *MoveRequest) (*Point, error) {
	if (ind < 0) || (ind >= len(data.Snakes)) {
		return nil, errors.New("Index out of bounds")
	}
	snake := data.Snakes[ind]
	return &(snake.Coords[len(snake.Coords)-1]), nil

}

func IsSnakeHead(p *Point, data *MoveRequest) bool {
	if p != nil && data.SnakeHeads[p.String()] {
		return true
	}
	return false
}

func getTaunt(turn int) string {
	if turn < 30 {
		return "This dad likes what he sees"
	} else if turn < 60 {
		return "My god you've grown"
	} else if turn < 90 {
		return "Let me get my glasses"
	}
	return "I need to go to bed"
}

// get the position of all neighbouring snake tiles and
// return the snake data corresponding to the last piece
// of snake that you see
// if there are no snakes around you return nil
func FindMinSnakePointInSurroundingArea(p *Point, data *MoveRequest, KeySnakeData map[int]*SnakeData) {
	pts := []*Point{
		p.UpHazard(data),
		p.DownHazard(data),
		p.LeftHazard(data),
		p.RightHazard(data)}

	for _, pt := range pts {
		if pt != nil {
			sd := data.SnakeHash[pt.String()]
			if sd != nil {
				if KeySnakeData[sd.id] == nil ||
					sd.lengthLeft < KeySnakeData[sd.id].lengthLeft {
					KeySnakeData[sd.id] = sd
				}
			}
		}
	}
}

// returns the number of valid neighbours to a point p
func GetNumNeighbours(data *MoveRequest, p *Point) (int, error) {
	if p == nil {
		return 0, nil
	}
	neighbours := 0
	for _, d := range []string{UP, DOWN, LEFT, RIGHT} {
		neighbour, err := GetPointInDirection(p, d, data)
		if err != nil {
			return 0, err
		}
		//fmt.Printf("In Loop neighbour %v, %v\n", p, d)
		if neighbour != nil {
			neighbours += 1
		}
	}
	//fmt.Printf("getting neighbours %v, %v\n", direc, neighbours)
	return neighbours, nil
}

// returns a point representing traven in the direction direc
// i.e. if you pass in direc "up" it will give you the point
// that is above p
// will only return points that are valid moves
func GetPointInDirection(p *Point, direc string, data *MoveRequest) (*Point, error) {
	if p == nil {
		return nil, nil
	}
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

func GetPointInDirectionHazards(p *Point, direc string, data *MoveRequest) (*Point, error) {
	if p == nil {
		return nil, nil
	}
	switch direc {
	case UP:
		return p.UpHazard(data), nil
	case DOWN:
		return p.DownHazard(data), nil
	case LEFT:
		return p.LeftHazard(data), nil
	case RIGHT:
		return p.RightHazard(data), nil
	}
	return nil, errors.New(fmt.Sprintf("could not find direction %v", direc))
}

func toStringPointer(str string) *string {
	return &str
}

func getMyHead(data *MoveRequest) (*Point, error) {
	for _, snake := range data.Snakes {
		if snake.Id == data.You && len(data.You) > 0 {
			return snake.Head(), nil
		}
	}
	return &Point{}, errors.New("Could not get head")
}

func getMyTail(data *MoveRequest) (*Point, error) {
	return getTail(data.MyIndex, data)
}
