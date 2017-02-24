package kaa

import (
	"reflect"
	"sort"
	"testing"
)

func TestFilterDFSMoveMax(t *testing.T) {
	_, err := NewMoveRequest(gameString3)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func Test_FilterPossibleMoves(t *testing.T) {
	data, err := NewMoveRequest(gameString3)

	if err != nil {
		t.Errorf("%v", err)
	}

	directions := FilterPossibleMoves(data, []string{UP, DOWN, LEFT, RIGHT})

	notRightOrDown := []string{LEFT, UP}

	// sort both of the strings so that deep equal will be able to see them
	sort.Strings(notRightOrDown)
	sort.Strings(directions)

	if !reflect.DeepEqual(directions, notRightOrDown) {
		t.Errorf("expected all directions except down, got %v", directions)
	}

}

func Test_ClosestFood(t *testing.T) {
	data, err := NewMoveRequest(gameString3)
	if err != nil {
		t.Errorf("%v", err)
	}

	directions := []string{LEFT, UP}
	foodDirections := ClosestFoodDirections(data, directions)
	expectedFoodDirections := []string{LEFT, UP}
	sort.Strings(foodDirections)
	sort.Strings(expectedFoodDirections)

	if !reflect.DeepEqual(foodDirections, expectedFoodDirections) {
		t.Errorf("expected %v directions got %v", expectedFoodDirections, foodDirections)
	}

}
