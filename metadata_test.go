package kaa

import (
	"fmt"
	"math"
	"reflect"
	"sort"
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
		Height: 20,
		Width:  20,
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
			},
		},
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
		moveMax := dirData.moveMax()
		if moveMax.Snakes != 0 {
			t.Errorf("Expected %v to be %v", moveMax.Snakes, 0)
		}

		// all moves are possible except for moving onto yourself
		if direc != LEFT {
			if moveMax.Moves != moves {
				t.Errorf("expected %v to be %v", moveMax.Moves, moves)
			}
		} else {
			if moveMax.Moves != 0 {
				t.Errorf("expected %v to be %v", moveMax.Moves, 0)
			}
		}
	}
}

func TestMetaDataWithMoves(t *testing.T) {
	req := &MoveRequest{
		Height: 20,
		Width:  20,
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
			},
		},
	}

	// need to make the hazards manually
	req.init()

	data, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	// all moves are possible except moving onto yourself
	if len(data[UP].MovesAway) == len(moves_to_depth) {
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
		all_food := data[UP].moveMax().Food
		if all_food != len(req.Food) {
			t.Errorf("Total food should be %v, got %v", len(req.Food), all_food)
		}
	} else {
		t.Errorf("Moves away should  be length %d got %d", len(moves_to_depth), len(data[UP].MovesAway))
	}
}

func TestClosestFood(t *testing.T) {
	req := &MoveRequest{
		Height: 20,
		Width:  20,
		Food: []Point{
			Point{X: 11, Y: 9},
			Point{X: 11, Y: 10},
			Point{X: 14, Y: 8},
			Point{X: 12, Y: 11},
		},
		Snakes: []Snake{
			Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 14, Y: 13},
					// snake layout
					//   | x |
				},
			},
		},
	}

	// need to make the hazards manually
	req.init()

	data, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	if data[LEFT].ClosestFood != 3 {
		t.Errorf("expected the closest food to the left to be 3, got %v", data[LEFT].ClosestFood)
	}

	if data[RIGHT].ClosestFood != 5 {
		t.Errorf("expected the closest food to the left to be 5, got %v", data[RIGHT].ClosestFood)
	}

	if data[UP].ClosestFood != 3 {
		t.Errorf("expected the closest food to the left to be 3, got %v", data[UP].ClosestFood)
	}

	if data[DOWN].ClosestFood != math.MaxInt64 {
		t.Errorf("expected the closest food to the left to be max int, got %v", data[DOWN].ClosestFood)
	}

	directions := FilterPossibleMoves(data)
	all := []string{LEFT, DOWN, UP, RIGHT}
	sort.Strings(all)
	sort.Strings(directions)

	if !reflect.DeepEqual(directions, all) {
		t.Errorf("expected all directions got %v", directions)
	}

	foodDirections := ClosestFoodDirections(data, directions)
	expectedFoodDirections := []string{LEFT, UP}
	sort.Strings(foodDirections)
	sort.Strings(expectedFoodDirections)

	if !reflect.DeepEqual(foodDirections, expectedFoodDirections) {
		t.Errorf("expected %v directions got %v", expectedFoodDirections, foodDirections)
	}

}

func TestRandom(t *testing.T) {
	req := &MoveRequest{
		GameId: "d0bd244e-91da-4e63-86e6-ea575376c3be",
		Height: 20,
		Width:  20,
		Turn:   20,
		Food: []Point{
			Point{X: 19, Y: 15}},
		Snakes: []Snake{Snake{
			Coords: []Point{
				Point{X: 7, Y: 12},
				Point{X: 7, Y: 13},
				Point{X: 7, Y: 14},
			},
			HealthPoints: 80,
			Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
			Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
			Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You:     "6db6f851-635b-4534-b882-6f219e0a1f6a",
		Hazards: map[string]bool{"{7,12}": true, "{7,13}": true, "{7,14}": true},
		FoodMap: map[string]bool(nil)}
	req.init()

	data, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	moves := bestMoves(data)
	fmt.Println(moves)
}
