package kaa

import (
	"reflect"
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

func TestSetMinSnakePointInArea(t *testing.T) {
	data, err := NewMoveRequest(`{"you":"0623b12a-411b-4674-a115-591063ef92d3","width":10,"turn":124,"snakes":[{"taunt":"battlesnake-go!","name":"7eef72e9-72fc-4c27-a387-898384639f46 (10x10)","id":"0623b12a-411b-4674-a115-591063ef92d3","health_points":96,"coords":[[9,1],[9,0],[8,0],[8,1],[8,2],[7,2],[7,3],[7,4],[7,5],[7,6],[6,6],[6,7],[5,7],[4,7],[3,7],[2,7],[1,7],[1,8],[0,8],[0,7],[0,6],[1,6],[2,6],[3,6],[4,6],[5,6],[5,5],[5,4],[5,3],[5,2],[6,2],[6,1]]}],"height":10,"game_id":"7eef72e9-72fc-4c27-a387-898384639f46","food":[[0,0],[1,3],[4,0]],"dead_snakes":[]}`)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	p := data.Direcs[DOWN].minKeySnakePart().pnt
	expected := Point{X: 0, Y: 8}
	if !reflect.DeepEqual(p, &expected) {
		t.Errorf("Expected %v to be %v", p, expected)
	}

	expected = Point{X: 6, Y: 1}
	tail, _ := getMyTail(data)
	if !reflect.DeepEqual(&tail, &expected) {
		t.Errorf("Expected %v to be %v", tail, expected)
	}
}
