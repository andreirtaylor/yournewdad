package kaa

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Hazard int

const (
	UP    = "up"
	DOWN  = "down"
	LEFT  = "left"
	RIGHT = "right"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (point *Point) UnmarshalJSON(data []byte) error {
	var coords []int
	json.Unmarshal(data, &coords)
	if len(coords) != 2 {
		return errors.New("Bad set of coordinates: " + string(data))
	}
	*point = Point{X: coords[0], Y: coords[1]}
	return nil
}

func (point *Point) String() string {
	return fmt.Sprintf("{%d,%d}", point.X, point.Y)
}

func (point *Point) getUp(data *MoveRequest, hazards bool) *Point {
	if point.Y == 0 {
		return nil
	}
	ret := &Point{point.X, point.Y - 1}
	if data.Hazards[ret.String()] && !hazards {
		return nil
	}
	return ret
}

// hazards tells if you want to return hazards
// do not use these functions use the functions below
func (point *Point) getDown(data *MoveRequest, hazards bool) *Point {
	if point.Y == data.Height-1 {
		return nil
	}
	ret := &Point{point.X, point.Y + 1}
	if data.Hazards[ret.String()] && !hazards {
		return nil
	}
	return ret
}

func (point *Point) getLeft(data *MoveRequest, hazards bool) *Point {
	if point.X == 0 {
		return nil
	}
	ret := &Point{point.X - 1, point.Y}
	if data.Hazards[ret.String()] && !hazards {
		return nil
	}
	return ret
}

func (point *Point) getRight(data *MoveRequest, hazards bool) *Point {
	if point.X == data.Width-1 {
		return nil
	}
	ret := &Point{point.X + 1, point.Y}
	if data.Hazards[ret.String()] && !hazards {
		return nil
	}
	return ret
}

// directional functions return a new point or nil if the point is out of the
// board
func (point *Point) Up(data *MoveRequest) *Point {
	return point.getUp(data, false)
}
func (point *Point) Down(data *MoveRequest) *Point {
	return point.getDown(data, false)
}

func (point *Point) Left(data *MoveRequest) *Point {
	return point.getLeft(data, false)
}
func (point *Point) Right(data *MoveRequest) *Point {
	if point.X == data.Width-1 {
		return nil
	}
	ret := &Point{point.X + 1, point.Y}
	if data.Hazards[ret.String()] {
		return nil
	}
	return ret
}

// Hazard functions return the points of Hazards as well as
// valid moves
// returns nil if the point is a wall
func (point *Point) DownHazard(data *MoveRequest) *Point {
	return point.getDown(data, false)
}

func (point *Point) UpHazard(data *MoveRequest) *Point {
	return point.getUp(data, true)
}

func (point *Point) LeftHazard(data *MoveRequest) *Point {
	return point.getLeft(data, true)
}

func (point *Point) RightHazard(data *MoveRequest) *Point {
	return point.getRight(data, true)
}
