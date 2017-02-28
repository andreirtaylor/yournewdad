package kaa

import (
	"reflect"
	"testing"
)

func Test_FullStats(t *testing.T) {
	data, err := NewMoveRequest(gameString9)

	if err != nil {
		t.Errorf("%v", err)
	}

	for _, snake := range data.Snakes {
		head := snake.Head()
		stats := fullStatsPnt(head, data)
		if stats.Food != 6 {
			t.Errorf("All 6 pieces of food are accessable by all snakes got %v", stats.Food)
		}
		if stats.Moves != 339 {
			t.Errorf("Moves should be 337 got %v", stats.Moves)
		}
	}

	stats := quickStats(&Point{X: 13, Y: 2}, data, 5, false)
	if stats.Food != 0 {
		t.Errorf("There is no food within 5 moves got %v", stats.Food)
	}
	if stats.ClosestFood != nil {
		t.Errorf("closest food should be nil within 5 moves got %v", stats.ClosestFood)
	}

	stats = quickStats(&Point{X: 13, Y: 2}, data, 10, false)
	if stats.Food != 1 {
		t.Errorf("There is 1 food within 10 moves got %v", stats.Food)
	}
	if reflect.DeepEqual(stats.ClosestFood, &Point{X: 0, Y: 8}) {
		t.Errorf("closest food should be nil within 5 moves got %v", stats.ClosestFood)
	}
}
