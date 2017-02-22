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
				HealthPoints: 80,
				Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You: "6db6f851-635b-4534-b882-6f219e0a1f6a",
	}

	// need to make the hazards manually

	err := GenerateMetaData(req)
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

	for direc, dirData := range req.Direcs {
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
				HealthPoints: 80,
				Id:           "6db6f851-635b-4534-b882-6f219e0a1f6a",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "6db6f851-635b-4534-b882-6f219e0a1f6a"},
		},
		You: "6db6f851-635b-4534-b882-6f219e0a1f6a",
	}

	err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}

	direcs := req.Direcs
	movesAway := direcs[UP].MovesAway
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
	maxMove, err := direcs[UP].moveMax()
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
		Snakes: []Snake{
			Snake{
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
	err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	direcs := req.Direcs

	if direcs[LEFT].ClosestFood != 3 {
		t.Errorf(
			"closest food to the left is 3 moves away, got %v",
			direcs[LEFT].ClosestFood)
	}

	if direcs[RIGHT].ClosestFood != 5 {
		t.Errorf(
			"closest food to the right should be 5 moves away got %v",
			direcs[RIGHT].ClosestFood)
	}

	if direcs[UP].ClosestFood != 3 {
		t.Errorf(
			"expected the closest food up to be 3 moves away, got %v",
			direcs[UP].ClosestFood)
	}

	if direcs[DOWN].ClosestFood != math.MaxInt64 {
		t.Errorf(
			"Going down is invalid should be max int, got %v",
			direcs[DOWN].ClosestFood)
	}

	directions, err := FilterPossibleMoves(req)
	if err != nil {
		t.Errorf("Unexpected error in filtering possible directions %v", err)
	}

	all := []string{LEFT, UP, RIGHT}
	sort.Strings(all)
	sort.Strings(directions)

	if !reflect.DeepEqual(directions, all) {
		t.Errorf("expected all directions except down, got %v", directions)
	}

	foodDirections := ClosestFoodDirections(req, directions)
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
	req := &MoveRequest{GameId: "d3684dd4-975b-4449-91ea-7051ea3f47da", Height: 8, Width: 8, Turn: 17, Food: []Point{Point{X: 0, Y: 3}, Point{X: 6, Y: 0}, Point{X: 1, Y: 3}, Point{X: 0, Y: 7}, Point{X: 7, Y: 7}, Point{X: 5, Y: 0}, Point{X: 7, Y: 5}, Point{X: 7, Y: 3}, Point{X: 6, Y: 6}, Point{X: 5, Y: 4}}, Snakes: []Snake{Snake{Coords: []Point{Point{X: 0, Y: 6}, Point{X: 1, Y: 6}, Point{X: 2, Y: 6}, Point{X: 2, Y: 5}, Point{X: 2, Y: 4}, Point{X: 2, Y: 3}, Point{X: 2, Y: 2}, Point{X: 1, Y: 2}, Point{X: 1, Y: 2}}, HealthPoints: 100, Id: "be171270-8030-412d-81d7-72e2e1e97895", Name: "d3684dd4-975b-4449-91ea-7051ea3f47da (8x8)", Taunt: "be171270-8030-412d-81d7-72e2e1e97895"}, Snake{Coords: []Point{Point{X: 5, Y: 7}, Point{X: 4, Y: 7}, Point{X: 3, Y: 7}, Point{X: 3, Y: 6}, Point{X: 3, Y: 5}, Point{X: 3, Y: 4}, Point{X: 3, Y: 3}, Point{X: 3, Y: 2}, Point{X: 3, Y: 1}, Point{X: 3, Y: 0}}, HealthPoints: 99, Id: "f220e2b6-7e02-4857-97e8-a5831d79ba78", Name: "d3684dd4-975b-4449-91ea-7051ea3f47da (8x8)", Taunt: "f220e2b6-7e02-4857-97e8-a5831d79ba78"}}, You: "be171270-8030-412d-81d7-72e2e1e97895"}

	err := GenerateMetaData(req)

	//fmt.Printf("%#v", req.Snakes[0])

	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	moves, err := bestMoves(req)
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
		Snakes: []Snake{Snake{
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

	err := GenerateMetaData(req)
	if err != nil {
		t.Errorf("Unexpected Errror %v", err)
	}
	moves, err := bestMoves(req)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best move %v", err)
	}
	if len(moves) > 1 {
		t.Errorf("The best move is left")
	}
}

func TestStrangeDeath(t *testing.T) {
	data, err := NewMoveRequest(`{"you":"89d8690a-d3ee-49ea-a272-968c7bf467d2","width":20,"turn":642,"snakes":[{"taunt":"battlesnake-go!","name":"06e65a2e-dc3b-4d19-80c3-aa436a0f7ca0 (20x20)","id":"89d8690a-d3ee-49ea-a272-968c7bf467d2","health_points":100,"coords":[[15,7],[15,6],[15,5],[15,4],[15,3],[16,3],[16,2],[15,2],[15,1],[14,1],[13,1],[13,0],[12,0],[11,0],[10,0],[9,0],[9,1],[10,1],[11,1],[12,1],[12,2],[13,2],[14,2],[14,3],[13,3],[12,3],[12,4],[12,5],[12,6],[12,7],[11,7],[11,8],[10,8],[9,8],[9,9],[9,10],[9,11],[9,12],[10,12],[11,12],[12,12],[13,12],[14,12],[15,12],[16,12],[17,12],[17,11],[17,10],[17,9],[18,9],[18,10],[18,11],[18,12],[18,13],[18,14],[18,15],[18,16],[18,17],[17,17],[16,17],[16,18],[15,18],[14,18],[13,18],[12,18],[11,18],[10,18],[9,18],[8,18],[7,18],[6,18],[5,18],[5,17],[5,16],[4,16],[4,17],[4,18],[3,18],[3,17],[3,16],[3,15],[3,14],[3,13],[2,13],[1,13],[1,12],[1,11],[1,10],[1,9],[1,8],[0,8],[0,7],[1,7],[2,7],[2,8],[2,9],[2,10],[2,11],[2,12],[3,12],[4,12],[5,12],[6,12],[6,11],[6,10],[6,9],[6,8],[5,8],[4,8],[4,7],[4,6],[4,5],[4,4],[4,3],[4,2],[3,2],[2,2],[2,1],[2,1]]}],"height":20,"game_id":"06e65a2e-dc3b-4d19-80c3-aa436a0f7ca0","food":[[2,5],[3,3],[8,3],[2,14],[19,18],[2,19],[13,19],[4,11],[1,0],[7,6]],"dead_snakes":[]}`)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	moves, err := bestMoves(data)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best move %v", err)
	}

	for _, move := range moves {
		if move == UP {
			t.Errorf("You cant move onto yourself")
		}
	}

}

func TestEfficientSpace(t *testing.T) {
	data, err := NewMoveRequest(`{"you":"dfda0e37-be0c-4ea6-a1b3-09bb6799c06a","width":10,"turn":80,"snakes":[{"taunt":"battlesnake-go!","name":"641321b4-48e4-420b-9358-72947fc21dfb (10x10)","id":"dfda0e37-be0c-4ea6-a1b3-09bb6799c06a","health_points":100,"coords":[[9,1],[9,2],[9,3],[9,4],[9,5],[9,6],[9,7],[8,7],[7,7],[6,7],[5,7],[4,7],[4,6],[3,6],[2,6],[1,6],[1,7],[1,8],[0,8],[0,7],[0,6],[0,5],[0,4],[0,3],[1,3],[1,2],[0,2],[0,1],[0,0],[1,0],[2,0],[3,0],[4,0],[5,0],[5,1],[5,2],[5,3],[5,3]]}],"height":10,"game_id":"641321b4-48e4-420b-9358-72947fc21dfb","food":[[2,9],[5,8],[6,0],[1,1],[1,4],[3,4],[7,3],[8,3],[4,9],[9,0]],"dead_snakes":[]}
`)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	move, err := bestMove(data)
	if err != nil {
		t.Errorf("Unexpected Errror whlie getting best move %v", err)
	}

	if move != LEFT {
		t.Errorf("Expected Move to be left got, %v", move)
	}

}
