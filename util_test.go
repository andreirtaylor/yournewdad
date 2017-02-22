package kaa

import (
	"testing"
)

func TestGetPointInDirection(t *testing.T) {
	data, err := NewMoveRequest(`{
		"you":"dfda0e37-be0c-4ea6-a1b3-09bb6799c06a",
		"width":10,
		"height":10,
		"turn":80,"snakes":[
			{"taunt":"battlesnake-go!",
			"name":"641321b4-48e4-420b-9358-72947fc21dfb (10x10)",
			"id":"dfda0e37-be0c-4ea6-a1b3-09bb6799c06a",
			"health_points":100,"coords":[ [2,3], [2,3]]
			}],
		"food":[[2,9],[5,8],[6,0],[1,1],[1,4],[3,4],[7,3]]}`)

	if err != nil {
		t.Errorf("%v", err)
	}

	p, err := GetPointInDirection(&Point{X: 2, Y: 4}, UP, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	if p != nil {
		t.Errorf("P should be nili %v", err)
	}

	p, err = GetPointInDirection(&Point{X: 9, Y: 4}, RIGHT, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	if p != nil {
		t.Errorf("P should be nili %v", err)
	}

	p, err = GetPointInDirection(&Point{X: 9, Y: 9}, DOWN, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	if p != nil {
		t.Errorf("P should be nili %v", err)
	}
}
