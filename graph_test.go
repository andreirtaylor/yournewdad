package kaa

import (
	"testing"
)

func Test_FullStats(t *testing.T) {
	data, err := NewMoveRequest(gameString9)

	if err != nil {
		t.Errorf("%v", err)
	}

	for _, snake := range data.Snakes {
		head := snake.Head()
		stats := fullStats(head, data)
		t.Logf("%v", stats)
		stats = quickStats(head, data, 5)
		t.Logf("%v", stats)
	}
}
