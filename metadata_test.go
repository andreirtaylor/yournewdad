package kaa

import (
	"reflect"
	"sort"
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
