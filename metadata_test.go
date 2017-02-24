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
	data, err := NewMoveRequest(`{"you":"82557bbc-5ff2-4e51-8133-f6875d4f8d71","width":10,"turn":233,"snakes":[{"taunt":"battlesnake-go!","name":"7eef72e9-72fc-4c27-a387-898384639f46 (10x10)","id":"82557bbc-5ff2-4e51-8133-f6875d4f8d71","health_points":100,"coords":[[1,3],[0,3],[0,4],[0,5],[0,6],[0,7],[1,7],[2,7],[3,7],[3,8],[3,9],[4,9],[4,8],[4,7],[4,6],[4,5],[5,5],[5,4],[6,4],[7,4],[7,3],[6,3],[5,3],[4,3],[4,4],[3,4],[3,3],[3,2],[4,2],[5,2],[5,1],[4,1],[3,1],[2,1],[2,0],[3,0],[4,0],[5,0],[6,0],[7,0],[8,0],[9,0],[9,1],[9,2],[9,2]]}],"height":10,"game_id":"7eef72e9-72fc-4c27-a387-898384639f46","food":[[6,2],[7,5],[2,3]],"dead_snakes":[]}`)

	if err != nil {
		t.Logf("error: %v", err)
	}

	head, err := getMyHead(data)
	if err != nil {
		t.Errorf("Getting Head %v", err)
	}

	if !reflect.DeepEqual(head, Point{X: 1, Y: 3}) {
		t.Errorf("Expected %v to be %v", head, Point{X: 1, Y: 3})
	}

	// all moves are possible except moving onto yourself

	//	moves := data.Width*data.Height - len(data.Snakes[0].Coords)

	//	for direc, dirData := range data.Direcs {
	//		// all moves are possible except for moving onto yourself
	//		moveMax, err := dirData.moveMax()
	//		if err != nil {
	//			t.Errorf("getting moveMax %v", err)
	//			continue
	//		}
	//		if direc != LEFT {
	//			if moveMax.Moves != moves[direc] {
	//				t.Errorf("expected %v to be %v", moveMax.Moves, moves)
	//			}
	//		} else {
	//			if moveMax != nil {
	//				t.Errorf("Moving left moves you onto your body, it is not a valid move")
	//			}
	//		}
	//	}
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

	directions := FilterPossibleMoves(req, []string{UP, DOWN, LEFT, RIGHT})

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

	if moves[0] == DOWN {
		t.Errorf("There is only one good move, and its UP, you gave me: %v", moves[0])
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

func TestEfficientSpace(t *testing.T) {
	data, err := NewMoveRequest(`{"you":"0623b12a-411b-4674-a115-591063ef92d3","width":10,"turn":124,"snakes":[{"taunt":"battlesnake-go!","name":"7eef72e9-72fc-4c27-a387-898384639f46 (10x10)","id":"0623b12a-411b-4674-a115-591063ef92d3","health_points":96,"coords":[[9,1],[9,0],[8,0],[8,1],[8,2],[7,2],[7,3],[7,4],[7,5],[7,6],[6,6],[6,7],[5,7],[4,7],[3,7],[2,7],[1,7],[1,8],[0,8],[0,7],[0,6],[1,6],[2,6],[3,6],[4,6],[5,6],[5,5],[5,4],[5,3],[5,2],[6,2],[6,1]]}],"height":10,"game_id":"7eef72e9-72fc-4c27-a387-898384639f46","food":[[0,0],[1,3],[4,0]],"dead_snakes":[]}`)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	moves, err := bestMoves(data)

	if err != nil {
		t.Errorf("%v", err)
	}

	if moves[0] != DOWN {
		t.Errorf("Expected Move to be left got, %v", moves[0])
	}

	totalMoves := data.Direcs[DOWN].TotalMoves
	if totalMoves != 30 {
		t.Errorf("Expected 30 total moves got, %v", totalMoves)
	}

}
