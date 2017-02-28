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
	return fmt.Sprintf("%d,%d", point.X, point.Y)
}

func (p1 *Point) isNeighbour(p2 *Point) bool {
	if p1 == nil || p2 == nil {
		return false
	}
	d := p1.Dist(p2)
	if d.X == 0 && d.Y == 1 {
		return true
	}
	if d.X == 1 && d.Y == 0 {
		return true
	}
	return false
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

// returns a map that has strings correspondings to the valid neighbours
func (p *Point) GetValidNeighboursMap(data *MoveRequest) map[string]bool {
	pts := []*Point{
		p.Up(data),
		p.Down(data),
		p.Left(data),
		p.Right(data)}

	m := make(map[string]bool)
	for _, pt := range pts {
		if pt != nil {
			m[pt.String()] = true
		}
	}
	return m
}

// returns which direction (as a string)
// the point is in
func (p1 *Point) WhichDirectionIs(p2 *Point) []string {
	if p2 == nil {
		return nil
	}
	var ret = []string{}

	if p1.X > p2.X {
		ret = append(ret, LEFT)
	} else if p1.X < p2.X {
		ret = append(ret, RIGHT)
	}

	if p1.Y < p2.Y {
		ret = append(ret, DOWN)
	} else if p1.Y > p2.Y {
		ret = append(ret, UP)
	}
	return ret
}

// returns the absolute distance from point p1 to point p2
// This distance is returned as a point in both the X and Y directions
func (p1 *Point) Dist(p2 *Point) *Point {
	if p2 == nil {
		return nil
	}
	ret := &Point{}
	if p1.X > p2.X {
		ret.X = p1.X - p2.X
	} else if p1.X < p2.X {
		ret.X = p2.X - p1.X
	}

	if p1.Y > p2.Y {
		ret.Y = p1.Y - p2.Y
	} else if p1.Y < p2.Y {
		ret.Y = p2.Y - p1.Y
	}
	return ret
}
