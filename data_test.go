package main

import (
	"testing"
)

func Test_NoFood(t *testing.T) {
	data, err := NewMoveRequest(gameString8)
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.NoFood() != true {
		t.Errorf("There is no food in your area")
	}
}

func Test_NoFoodWithFood(t *testing.T) {
	data, err := NewMoveRequest(gameString7)
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.NoFood() != false {
		t.Errorf("There is food in your area")
	}
}
