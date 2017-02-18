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
			Point{X: 5, Y: 8},
		},
		Snakes: []*Snake{
			&Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 13, Y: 12},
					Point{X: 13, Y: 13},
					// snake layout
					//   | x | x |
					//   | x |   |
				},
				HealthPoints: 80,
				Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You: "6db6f851-635b-4534-b882-6f219e0a1f6a",
	}

	// need to make the hazards manually
	req.init()

	data, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	head, err := getMyHead(req)
	if err != nil {
		t.Errorf("Getting Head %v", err)
	}

	if !reflect.DeepEqual(head, Point{X: 14, Y: 12}) {
		t.Errorf("Expected %v to be %v", head, Point{X: 14, Y: 12})
	}

	// all moves are possible except moving onto yourself

	moves := req.Width*req.Height - len(req.Snakes[0].Coords)

	for direc, dirData := range data.MD {
		// all moves are possible except for moving onto yourself
		moveMax, err := dirData.moveMax()
		if direc != LEFT {

			if moveMax.Moves != moves {
				t.Errorf("expected %v to be %v", moveMax.Moves, moves)
			}
			if err != nil {
				t.Errorf("getting moveMax %v", err)
				continue
			}
			if moveMax.Snakes != 0 {
				t.Errorf("Expected %v to be %v", moveMax.Snakes, 0)
			}

		} else {
			if moveMax != nil {
				t.Errorf("Moving left moves you onto your body, it is not a valid move")
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
		Snakes: []*Snake{
			&Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 13, Y: 12},
					Point{X: 13, Y: 13},
					// snake layout
					//   | x | x |
					//   | x |   |
				},
				HealthPoints: 80,
				Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You: "6db6f851-635b-4534-b882-6f219e0a1f6a",
	}

	// need to make the hazards manually
	req.init()

	req, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	data := req.MD
	movesAway := data[UP].MovesAway
	//for _, x := range movesAway {
	//	fmt.Printf("%v\n", x)
	//}

	// movesAway is indexed from 0
	m_1 := movesAway[0]
	if m_1.Food != 0 {
		t.Errorf("move 1 should have 0 food, got ", m_1.Food)
	}
	m_3 := movesAway[2]
	if m_3.Food != 1 {
		t.Errorf("move 3 should have 1 food, got %v", m_3.Food)
	}
	m_5 := movesAway[4]
	if m_5.Food != 2 {
		t.Errorf("move 5 should have 2 food, got %v", m_5.Food)
	}
	maxMove, err := data[UP].moveMax()
	if err != nil {
		t.Errorf("getting moveMax %v", err)
	}
	all_food := maxMove.Food
	if all_food != len(req.Food) {
		t.Errorf("Total food should be %v, got %v", len(req.Food), all_food)
	}

	head, err := getMyHead(req)
	if err != nil {
		t.Errorf("getting NumNeighbours up,  %v", err)
	}

	num_neighbours, err := GetNumNeighbours(req, head.Up(req))
	if err != nil {
		t.Errorf("getting NumNeighbours up,  %v", err)
	}

	if num_neighbours != 3 {
		t.Errorf("expected 3 neighbours got, %v", err)
	}

	num_neighbours, err = GetNumNeighbours(req, head.Down(req))
	if num_neighbours != 2 {
		t.Errorf("expected 2 neighbours got, %v", err)
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
		Snakes: []*Snake{
			&Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 14, Y: 13},
					// snake layout
					//   | x |
				},
				HealthPoints: 80,
				Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You: "6db6f851-635b-4534-b882-6f219e0a1f6a",
	}

	// need to make the hazards manually
	req.init()

	req, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	data := req.MD

	if data[LEFT].ClosestFood != 3 {
		t.Errorf(
			"closest food to the left is 3 moves away, got %v",
			data[LEFT].ClosestFood)
	}

	if data[RIGHT].ClosestFood != 5 {
		t.Errorf(
			"closest food to the right should be 5 moves away got %v",
			data[RIGHT].ClosestFood)
	}

	if data[UP].ClosestFood != 3 {
		t.Errorf(
			"expected the closest food up to be 3 moves away, got %v",
			data[UP].ClosestFood)
	}

	if data[DOWN].ClosestFood != math.MaxInt64 {
		t.Errorf(
			"Going down is invalid should be max int, got %v",
			data[DOWN].ClosestFood)
	}

	directions, err := FilterPossibleMoves(data)
	if err != nil {
		t.Errorf("Unexpected error in filtering possible directions %v", err)
	}

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

func TestSmallSpaceWithFood(t *testing.T) {
	// in this test moving down will result in certain death
	// if the snake wants to move down it is wrong!
	req := &MoveRequest{GameId: "d3684dd4-975b-4449-91ea-7051ea3f47da", Height: 8, Width: 8, Turn: 17, Food: []Point{Point{X: 0, Y: 3}, Point{X: 6, Y: 0}, Point{X: 1, Y: 3}, Point{X: 0, Y: 7}, Point{X: 7, Y: 7}, Point{X: 5, Y: 0}, Point{X: 7, Y: 5}, Point{X: 7, Y: 3}, Point{X: 6, Y: 6}, Point{X: 5, Y: 4}}, Snakes: []*Snake{&Snake{Coords: []Point{Point{X: 0, Y: 6}, Point{X: 1, Y: 6}, Point{X: 2, Y: 6}, Point{X: 2, Y: 5}, Point{X: 2, Y: 4}, Point{X: 2, Y: 3}, Point{X: 2, Y: 2}, Point{X: 1, Y: 2}, Point{X: 1, Y: 2}}, HealthPoints: 100, Id: "be171270-8030-412d-81d7-72e2e1e97895", Name: "d3684dd4-975b-4449-91ea-7051ea3f47da (8x8)", Taunt: "be171270-8030-412d-81d7-72e2e1e97895"}, &Snake{Coords: []Point{Point{X: 5, Y: 7}, Point{X: 4, Y: 7}, Point{X: 3, Y: 7}, Point{X: 3, Y: 6}, Point{X: 3, Y: 5}, Point{X: 3, Y: 4}, Point{X: 3, Y: 3}, Point{X: 3, Y: 2}, Point{X: 3, Y: 1}, Point{X: 3, Y: 0}}, HealthPoints: 99, Id: "f220e2b6-7e02-4857-97e8-a5831d79ba78", Name: "d3684dd4-975b-4449-91ea-7051ea3f47da (8x8)", Taunt: "f220e2b6-7e02-4857-97e8-a5831d79ba78"}}, You: "be171270-8030-412d-81d7-72e2e1e97895"}

	req.init()

	req, err := GenerateMetaData(req)

	fmt.Printf("%#v", req.Snakes[0])

	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	data := req.MD
	moves, err := bestMoves(data)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best moves %v", err)
	}

	if len(moves) > 1 {
		t.Errorf("There is only one good move")
	}

	move, err := bestMove(req)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best move %v", err)
	}
	if move == DOWN {
		t.Errorf("There is only one good move, and its UP, you gave me: %v", move)
	}
}

// not really a test, this is for checking out if the snake works as expected
func TestRandom(t *testing.T) {
	req := &MoveRequest{
		GameId: "d0bd244e-91da-4e63-86e6-ea575376c3be",
		Height: 5,
		Width:  5,
		Turn:   20,
		Food: []Point{
			Point{X: 0, Y: 0},
			Point{X: 4, Y: 4},
			Point{X: 3, Y: 2},
		},
		Snakes: []*Snake{&Snake{
			Coords: []Point{
				Point{X: 1, Y: 0},
				Point{X: 1, Y: 1},
				Point{X: 1, Y: 2},
				Point{X: 1, Y: 3},
				Point{X: 1, Y: 4},
			},
			HealthPoints: 80,
			Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
			Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
			Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You: "6db6f851-635b-4534-b882-6f219e0a1f6a",
	}
	req.init()

	req, err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	data := req.MD
	moves, err := bestMoves(data)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best move %v", err)
	}
	if len(moves) > 1 {
		t.Errorf("The best move is left")
	}
}
