package kaa

import (
	"errors"
	"fmt"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// only handles valid moves right now
func MoveSnakeForward(ind int, data *MoveRequest, direc string) error {
	if (ind < 0) || (ind >= len(data.Snakes)) {
		return errors.New("Index out of bounds")
	}
	head := data.Snakes[ind].Head()
	fmt.Printf(" head is %v\n", head)

	p, err := GetPointInDirection(head, direc, data)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("Invalid move")
	}
	data.Hazards[p.String()] = true

	data.Snakes[ind].HeadPoint = p

	if !data.FoodMap[p.String()] {
		t, err := getTail(ind, data)
		if err != nil {
			return err
		}
		data.Hazards[t.String()] = false
		data.Snakes[ind].TailStack.Push(t)
	}
	// append the coords to the front of the snake
	data.Snakes[ind].Coords = append([]Point{Point{X: p.X, Y: p.Y}}, (data.Snakes[ind].Coords)...)
	return nil
}

func MoveSnakeBackward(ind int, data *MoveRequest) error {
	if (ind < 0) || (ind >= len(data.Snakes)) {
		return errors.New("Index out of bounds")
	}
	// assumes the snakes are all more than length 1
	p := data.Snakes[ind].Head()

	data.Hazards[p.String()] = false
	data.Snakes[ind].Coords = data.Snakes[ind].Coords[1:]
	//fmt.Printf("%v\n", data.Snakes[ind].Coords)
	if !data.FoodMap[p.String()] {
		t := data.Snakes[ind].TailStack.Pop()
		data.Hazards[t.String()] = false
		data.Snakes[ind].Coords = append(data.Snakes[ind].Coords, *t)
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
