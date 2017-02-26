package kaa

import (
	"reflect"
	"testing"
)

func Test_FindNoGuaranteedClosestFood(t *testing.T) {
	data, err := NewMoveRequest(gameString9)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	p := findGuaranteedClosestFood(data, UP)
	if p != nil {
		t.Errorf("you are not closest to any food")
	}
}

func Test_FindGuaranteedClosestFood(t *testing.T) {
	data, err := NewMoveRequest(gameString10)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	p := findGuaranteedClosestFood(data, DOWN)
	if p == nil {
		t.Errorf("you are closest to some food")
	}

	//t.Logf("%#v", p)
	if !reflect.DeepEqual(p.pnt, &Point{X: 12, Y: 18}) {
		t.Errorf("expected something else, %v", p)
	}
}

func Test_MoveSnakeOntoFood(t *testing.T) {
	data, err := NewMoveRequest(gameString4)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	tail, err := getTail(0, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = MoveSnakeForward(0, data, LEFT)
	if err != nil {
		t.Errorf("%v", err)
	}

	p := &Point{X: 0, Y: 15}
	// the head moves onto food meaning the tail stays where it is
	if !data.Hazards[p.String()] || !data.Hazards[tail.String()] {
		t.Errorf("If the head moves onto a food the tail should not move")
	}

	err = MoveSnakeBackward(0, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	// the head moves back and the tail remains in the same spot
	if data.Hazards[p.String()] || !data.Hazards[tail.String()] {
		t.Errorf("If the head moves back from a food the tail should stay in the same spot")
	}

}

func Test_MoveSnakeOntoBlankSpace(t *testing.T) {
	data, err := NewMoveRequest(gameString4)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	tail, err := getTail(0, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.Hazards[tail.String()] == false {
		t.Errorf("The tail should be a hazard")
	}

	p := &Point{X: 1, Y: 16}
	if data.Hazards[p.String()] == true {
		t.Errorf("The position of the new head should not be a hazard")
	}

	err = MoveSnakeForward(0, data, DOWN)
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.Hazards[p.String()] == false {
		t.Errorf("The new position of the head should be a hazzard")
	}

	if data.Hazards[tail.String()] == true {
		t.Errorf("The tail should no longet be a hazard")
	}
	err = MoveSnakeBackward(0, data)
	if err != nil {
		t.Errorf("%v", err)
	}

	oldHead := &Point{X: 1, Y: 16}
	if data.Hazards[oldHead.String()] == true {
		t.Errorf("The old head should not be a hazard")
	}

	if data.Hazards[tail.String()] == true {
		t.Errorf("The old tail should come back")
	}

}
func Test_getMyHead(t *testing.T) {
	data, err := NewMoveRequest(gameString3)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	head, err := getMyHead(data)
	if err != nil {
		t.Errorf("Getting Head %v", err)
	}

	if !reflect.DeepEqual(head, &Point{X: 3, Y: 9}) {
		t.Errorf("Expected %v to be %v", head, Point{X: 1, Y: 3})
	}

}

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
	data, err := NewMoveRequest(gameString1)
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
	if !reflect.DeepEqual(tail, &expected) {
		t.Errorf("Expected %v to be %v", tail, expected)
	}

	head, err := getMyHead(data)
	if err != nil {
		t.Errorf("getting NumNeighbours up,  %v", err)
	}
	if !reflect.DeepEqual(head, &Point{X: 9, Y: 1}) {
		t.Errorf("head should be %v got %v", &Point{X: 9, Y: 1}, head)
	}

}

func Test_NumberofNeighbours(t *testing.T) {
	data, err := NewMoveRequest(gameString1)

	n, err := GetNumNeighbours(data, &Point{X: 0, Y: 0})
	if err != nil {
		t.Errorf("%v", err)
	}
	if n != 2 {
		t.Errorf("Expected 2 neighbours got %v", n)
	}

	n, err = GetNumNeighbours(data, &Point{X: 6, Y: 5})
	if err != nil {
		t.Errorf("%v", err)
	}
	if n != 1 {
		t.Errorf("Expected 2 neighbours got %v", n)
	}
}
