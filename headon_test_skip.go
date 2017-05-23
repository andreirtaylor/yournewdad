package main

import (
	"reflect"
	"testing"
)

func Test_HeadOn(t *testing.T) {
	data, err := NewMoveRequest(gameString21)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !headOn(data, 0) {
		t.Errorf("You are in a bad spot with that snake")
	}
}

func Test_HeadOnMove(t *testing.T) {
	data, err := NewMoveRequest(gameString21)

	if err != nil {
		t.Errorf("%v", err)
	}

	directions, err := bestMoves(data)
	if err != nil {
		t.Errorf("%v", err)
	}

	left := []string{RIGHT}

	// sort both of the strings so that deep equal will be able to see them

	if !reflect.DeepEqual(directions, left) {
		t.Errorf("expected left got, %v", directions)
	}
}
