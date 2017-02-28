package kaa

import (
	"reflect"
	"testing"
)

func TestClosestFoodNoFood(t *testing.T) {
	data, err := NewMoveRequest(gameString1)

	if err != nil {
		t.Logf("error: %v", err)
	}
	for direc, direcData := range data.Direcs {
		if len(direcData.sortedFood) != 0 {
			t.Errorf("No food in direction %v", direc)
		}
	}

}

func TestClosestFoodWithFood(t *testing.T) {
	data, err := NewMoveRequest(gameString5)

	if err != nil {
		t.Logf("error: %v", err)
		return
	}

	direcs := data.Direcs
	expected := &Point{X: 2, Y: 8}
	if !reflect.DeepEqual(direcs[LEFT].sortedFood[0].pnt, expected) {
		t.Errorf("closest food LEFT is %v, got %v",
			expected, direcs[LEFT].sortedFood[0].pnt)
	}

	if len(direcs[RIGHT].sortedFood) != 0 {
		t.Errorf("there is no food to the right thats your body")
	}

	expected = &Point{X: 5, Y: 7}
	if !reflect.DeepEqual(direcs[UP].sortedFood[0].pnt, expected) {
		t.Errorf(
			"closest food UP is %v moves away, got %v",
			expected, direcs[UP].sortedFood[0].pnt)
	}

	expected = &Point{X: 5, Y: 7}
	if !reflect.DeepEqual(direcs[UP].sortedFood[0].pnt, expected) {
		t.Errorf(
			"closest food DOWN is %v moves away, got %v",
			expected, direcs[UP].sortedFood[0].pnt)
	}

}

func Test_EfficientSpace(t *testing.T) {
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

	totalMoves := data.Direcs[DOWN].Moves
	if totalMoves != 30 {
		t.Errorf("Expected 30 total moves got, %v", totalMoves)
	}

}

func Test_MovesVsSpace(t *testing.T) {
	data, err := NewMoveRequest(gameString7)
	if err != nil {
		t.Logf("error: %v", err)
		return
	}

	mvs := data.Direcs[LEFT].MovesVsSpace
	if mvs != 11 {
		t.Errorf("Expected 30 total moves got, %v", mvs)
	}

}

func Test_MovesVsSpace2(t *testing.T) {
	data, err := NewMoveRequest(gameString11)
	if err != nil {
		t.Logf("error: %v", err)
		return
	}

	mvs := data.Direcs[LEFT].MovesVsSpace
	expected := 1
	if mvs != expected {
		t.Errorf("Expected %v total moves got, %v", expected, mvs)
	}

	mvs = data.Direcs[UP].MovesVsSpace
	expected = 1
	if mvs != expected {
		t.Errorf("Expected %v total moves got, %v", expected, mvs)
	}

}
