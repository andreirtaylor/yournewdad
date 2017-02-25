package kaa

import (
	"reflect"
	"testing"
)

func TestMetaDataOnlyOneSnake(t *testing.T) {
	data, err := NewMoveRequest(gameString3)

	if err != nil {
		t.Logf("error: %v", err)
	}

	head, err := getMyHead(data)
	if err != nil {
		t.Errorf("Getting Head %v", err)
	}

	if !reflect.DeepEqual(head, Point{X: 3, Y: 9}) {
		t.Errorf("Expected %v to be %v", head, Point{X: 1, Y: 3})
	}

}

func TestClosestFoodNoFood(t *testing.T) {
	data, err := NewMoveRequest(gameString1)

	if err != nil {
		t.Logf("error: %v", err)
	}
	direcs := data.Direcs

	expected := 0
	if direcs[LEFT].ClosestFood != expected {
		t.Errorf(
			"closest food LEFT is %v moves away, got %v",
			expected, direcs[LEFT].ClosestFood)
	}

	if direcs[RIGHT].ClosestFood != expected {
		t.Errorf(
			"closest food RIGHT is %v moves away, got %v",
			expected, direcs[RIGHT].ClosestFood)
	}

	if direcs[UP].ClosestFood != expected {
		t.Errorf(
			"closest food UP is %v moves away, got %v",
			expected, direcs[UP].ClosestFood)
	}

	if direcs[DOWN].ClosestFood != expected {
		t.Errorf(
			"closest food DOWN is %v moves away, got %v",
			expected, direcs[DOWN].ClosestFood)
	}

}

func TestClosestFoodWithFood(t *testing.T) {
	data, err := NewMoveRequest(gameString5)

	if err != nil {
		t.Logf("error: %v", err)
		return
	}

	direcs := data.Direcs
	expected := 3
	if direcs[LEFT].ClosestFood != expected {
		t.Errorf(
			"closest food LEFT is %v moves away, got %v",
			expected, direcs[LEFT].ClosestFood)
	}

	expected = 0
	if direcs[RIGHT].ClosestFood != expected {
		t.Errorf(
			"closest food RIGHT is %v moves away, got %v",
			expected, direcs[RIGHT].ClosestFood)
	}

	expected = 1
	if direcs[UP].ClosestFood != expected {
		t.Errorf(
			"closest food UP is %v moves away, got %v",
			expected, direcs[UP].ClosestFood)
	}

	expected = 5
	if direcs[DOWN].ClosestFood != expected {
		t.Errorf(
			"closest food DOWN is %v moves away, got %v",
			expected, direcs[DOWN].ClosestFood)
	}

}

func TestSmallSpaceWithFood(t *testing.T) {
	// in this test moving down will result in certain death
	// if the snake wants to move down it is wrong!
	req := &MoveRequest{GameId: "d3684dd4-975b-4449-91ea-7051ea3f47da", Height: 8, Width: 8, Turn: 17, Food: []Point{Point{X: 0, Y: 3}, Point{X: 6, Y: 0}, Point{X: 1, Y: 3}, Point{X: 0, Y: 7}, Point{X: 7, Y: 7}, Point{X: 5, Y: 0}, Point{X: 7, Y: 5}, Point{X: 7, Y: 3}, Point{X: 6, Y: 6}, Point{X: 5, Y: 4}}, Snakes: []Snake{Snake{Coords: []Point{Point{X: 0, Y: 6}, Point{X: 1, Y: 6}, Point{X: 2, Y: 6}, Point{X: 2, Y: 5}, Point{X: 2, Y: 4}, Point{X: 2, Y: 3}, Point{X: 2, Y: 2}, Point{X: 1, Y: 2}, Point{X: 1, Y: 2}}, HealthPoints: 100, Id: "be171270-8030-412d-81d7-72e2e1e97895", Name: "d3684dd4-975b-4449-91ea-7051ea3f47da (8x8)", Taunt: "be171270-8030-412d-81d7-72e2e1e97895"}, Snake{Coords: []Point{Point{X: 5, Y: 7}, Point{X: 4, Y: 7}, Point{X: 3, Y: 7}, Point{X: 3, Y: 6}, Point{X: 3, Y: 5}, Point{X: 3, Y: 4}, Point{X: 3, Y: 3}, Point{X: 3, Y: 2}, Point{X: 3, Y: 1}, Point{X: 3, Y: 0}}, HealthPoints: 99, Id: "f220e2b6-7e02-4857-97e8-a5831d79ba78", Name: "d3684dd4-975b-4449-91ea-7051ea3f47da (8x8)", Taunt: "f220e2b6-7e02-4857-97e8-a5831d79ba78"}}, You: "be171270-8030-412d-81d7-72e2e1e97895"}

	err := GenerateMetaData(req)

	//fmt.Printf("%#v", req.Snakes[0])

	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	moves, err := bestMoves(req)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best moves %v", err)
	}

	if len(moves) > 1 {
		t.Errorf("There is only one good move")
	}

	if moves[0] == DOWN {
		t.Errorf("There is only one good move, and its UP, you gave me: %v", moves[0])
	}
}

func TestEfficientSpace(t *testing.T) {
	data, err := NewMoveRequest(gameString1)
	if err != nil {
		t.Logf("error: %v", err)
		return
	}

	moves, err := bestMoves(data)

	if err != nil {
		t.Errorf("%v", err)
	}

	if moves[0] != DOWN {
		t.Errorf("Expected Move to be left got, %v", moves[0])
	}

	totalMoves := data.Direcs[DOWN].TotalMoves
	if totalMoves != 30 {
		t.Errorf("Expected 30 total moves got, %v", totalMoves)
	}

}
