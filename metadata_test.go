package kaa

import (
	"fmt"
	"reflect"
	"testing"
)

func keepFMT() {
	fmt.Printf("")
}
func TestPoint(t *testing.T) {
	var point Point

	// test string
	point = Point{X: 5, Y: 13}
	if point.String() != "{5,13}" {
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

func TestMetaDataOnlyOneSnake(t *testing.T) {
	req := &MoveRequest{
		GameId: "d0bd244e-91da-4e63-86e6-ea575376c3be",
		Height: 20,
		Width:  20,
		Turn:   4,
		Food: []Point{
			Point{X: 5, Y: 13},
		},
		Snakes: []Snake{
			Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 13, Y: 12},
					Point{X: 13, Y: 13},
					// snake layout
					//   | x | x |
					//   | x |   |
				},
				HealthPoints: 96,
				Id:           "639fb7cd-2590-4418-abcc-3da577559fc6",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "639fb7cd-2590-4418-abcc-3da577559fc6",
			},
		},
		You: "639fb7cd-2590-4418-abcc-3da577559fc6",
	}

	// need to make the hazards manually
	req.GenHazards()

	data, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	head, _ := getMyHead(req)
	if !reflect.DeepEqual(head, Point{X: 14, Y: 12}) {
		t.Errorf("Expected %v to be %v", head, Point{X: 14, Y: 12})
	}

	// all moves are possible except moving onto yourself

	moves := req.Width*req.Height - len(req.Snakes[0].Coords)
	for direc, dirData := range data {
		if dirData.Snakes != 0 {
			t.Errorf("Expected %v to be %v", dirData.Snakes, 0)
		}

		// all moves are possible except for moving onto yourself
		if direc != LEFT {
			if dirData.Moves != moves {
				t.Errorf("expected %v to be %v", dirData.Moves, moves)
			}
		} else {
			if dirData.Moves != 0 {
				t.Errorf("expected %v to be %v", dirData.Moves, 0)
			}
		}
	}
}

func TestMetaDataWithMoves(t *testing.T) {
	req := &MoveRequest{
		GameId: "d0bd244e-91da-4e63-86e6-ea575376c3be",
		Height: 20,
		Width:  20,
		Turn:   4,
		Food: []Point{
			Point{X: 14, Y: 9},
			Point{X: 11, Y: 10},
			Point{X: 11, Y: 9},
		},
		Snakes: []Snake{
			Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 13, Y: 12},
					Point{X: 13, Y: 13},
					// snake layout
					//   | x | x |
					//   | x |   |
				},
				HealthPoints: 96,
				Id:           "639fb7cd-2590-4418-abcc-3da577559fc6",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "639fb7cd-2590-4418-abcc-3da577559fc6",
			},
		},
		You: "639fb7cd-2590-4418-abcc-3da577559fc6",
	}

	// need to make the hazards manually
	req.init()

	data, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	// all moves are possible except moving onto yourself
	if len(data[UP].MovesAway) == 3 {
		movesAway := data[UP].MovesAway

		m_1 := movesAway[MOVE_ONE]
		if m_1.Food != 0 {
			t.Errorf("move 1 should have 0 food, got ", m_1.Food)
		}
		m_3 := movesAway[MOVE_THREE]
		if m_3.Food != 1 {
			t.Errorf("move 3 should have 1 food, got %v", m_3.Food)
		}
		m_5 := movesAway[MOVE_FIVE]
		if m_5.Food != 2 {
			t.Errorf("move 5 should have 2 food, got %v", m_5.Food)
		}
		all_food := data[UP].Food
		if all_food != len(req.Food) {
			t.Errorf("Total food should be %v, got %v", len(req.Food), all_food)
		}
	} else {
		t.Errorf("Moves away should not be nil")
	}
}
