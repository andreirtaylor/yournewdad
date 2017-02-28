package kaa

import (
	"reflect"
	"testing"
)

func Test_Dist(t *testing.T) {
	p := &Point{X: 5, Y: 13}

	output := p.Dist(nil)
	if output != nil {
		t.Errorf("Should not return anything when there is not another point")
	}

	p = &Point{X: 5, Y: 13}
	p2 := &Point{X: 0, Y: 13}

	output = p.Dist(p2)
	expected := &Point{X: 5, Y: 0}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 0, Y: 10}

	output = p.Dist(p2)
	expected = &Point{X: 5, Y: 3}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 5, Y: 10}

	output = p.Dist(p2)
	expected = &Point{X: 0, Y: 3}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 8, Y: 10}

	output = p.Dist(p2)
	expected = &Point{X: 3, Y: 3}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 8, Y: 13}

	output = p.Dist(p2)
	expected = &Point{X: 3, Y: 0}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}
	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 8, Y: 15}

	output = p.Dist(p2)
	expected = &Point{X: 3, Y: 2}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 5, Y: 15}

	output = p.Dist(p2)
	expected = &Point{X: 0, Y: 2}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 2, Y: 15}

	output = p.Dist(p2)
	expected = &Point{X: 3, Y: 2}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}
}

func Test_WhichDirectionIs(t *testing.T) {
	p := &Point{X: 5, Y: 13}

	output := p.WhichDirectionIs(nil)
	if output != nil {
		t.Errorf("Should not return anything when there is not another point")
	}

	p = &Point{X: 5, Y: 13}
	p2 := &Point{X: 0, Y: 13}

	output = p.WhichDirectionIs(p2)
	expected := []string{LEFT}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 0, Y: 10}

	output = p.WhichDirectionIs(p2)
	expected = []string{LEFT, UP}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 5, Y: 10}

	output = p.WhichDirectionIs(p2)
	expected = []string{UP}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 8, Y: 10}

	output = p.WhichDirectionIs(p2)
	expected = []string{RIGHT, UP}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 8, Y: 13}

	output = p.WhichDirectionIs(p2)
	expected = []string{RIGHT}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}
	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 8, Y: 15}

	output = p.WhichDirectionIs(p2)
	expected = []string{RIGHT, DOWN}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 5, Y: 15}

	output = p.WhichDirectionIs(p2)
	expected = []string{DOWN}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}

	p = &Point{X: 5, Y: 13}
	p2 = &Point{X: 2, Y: 15}

	output = p.WhichDirectionIs(p2)
	expected = []string{LEFT, DOWN}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("Expected %v to be %v", output, expected)
	}
}

func TestPoint(t *testing.T) {
	var point Point

	// test string
	point = Point{X: 5, Y: 13}
	if point.String() != "5,13" {
		t.Errorf("expected point to print %v got %v", "{5,13}", point)
	}
	data := &MoveRequest{
		Height: 20,
		Width:  20,
	}

	// test x limit
	point = Point{X: 0, Y: 13}
	if point.Left(data) != nil {
		t.Errorf("point %v should not be able to go left", point)
	}

	// test x limit
	point = Point{X: 19, Y: 13}
	if point.Right(data) != nil {
		t.Errorf("point %v should not be able to go right", point)
	}

	// test x limit
	point = Point{X: 19, Y: 19}
	if point.Down(data) != nil {
		t.Errorf("point %v should not be able to go down", point)
	}

	// test x limit
	point = Point{X: 19, Y: 0}
	if point.Up(data) != nil {
		t.Errorf("point %v should not be able to go up", point)
	}

}
