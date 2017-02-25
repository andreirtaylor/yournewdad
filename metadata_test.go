package kaa

import (
	"testing"
)

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

func TestWeirdDeath(t *testing.T) {
	data, err := NewMoveRequest(` {"you":"fe292f8e-d74f-46f8-a85e-b5688d0e8ca0","width":20,"turn":497,"snakes":[{"taunt":"Dad 2.0 Ready","name":"Your New Dad","id":"fe292f8e-d74f-46f8-a85e-b5688d0e8ca0","health_points":70,"coords":[[13,15],[14,15],[14,16],[13,16],[12,16],[11,16],[10,16],[9,16],[9,15],[9,14],[9,13],[10,13],[11,13],[12,13],[12,12],[12,11],[12,10],[12,9],[12,8],[12,7],[13,7],[14,7],[15,7],[15,8],[16,8],[16,7],[17,7],[17,8],[18,8],[18,7],[18,6],[17,6],[16,6],[15,6],[14,6],[13,6],[12,6],[11,6],[11,7],[11,8],[11,9],[10,9],[9,9],[9,10],[9,11],[10,11],[11,11],[11,12],[10,12],[9,12],[8,12],[8,13],[8,14],[8,15],[8,16],[8,17],[9,17],[10,17],[11,17],[12,17],[13,17],[14,17],[15,17],[16,17],[17,17],[18,17],[18,16],[18,15],[17,15],[16,15],[15,15],[15,14],[16,14],[16,13],[16,12],[16,11],[16,10],[16,9]]}],"height":20,"game_id":"e5fe7eb6-6420-4499-a387-d8537b40c6d1","food":[[6,18],[7,4],[6,2],[6,10],[4,14],[16,3]],"dead_snakes":[]}
`)
	if err != nil {
		t.Errorf("%v", err)
	}

	moves, err := bestMoves(data)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best moves %v", err)
	}
	t.Logf("%v\n", moves)
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
		return
	}

	if moves[0] != DOWN {
		t.Errorf("Expected Move to be left got, %v", moves[0])
	}

	totalMoves := data.Direcs[DOWN].TotalMoves
	if totalMoves != 30 {
		t.Errorf("Expected 30 total moves got, %v", totalMoves)
	}

}
