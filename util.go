package kaa

import (
	"errors"
	"fmt"
)

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

func toStringPointer(str string) *string {
	return &str
}

func getMyHead(data *MoveRequest) (Point, error) {
	for _, snake := range data.Snakes {
		if snake.Id == data.You && len(data.You) > 0 {
			return snake.Head(), nil
		}
	}
	return Point{}, errors.New("Could not get head")
}
