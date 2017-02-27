package kaa

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
)

// used to keep the fmt function around for filters
func keepFMTForFilters() {
	fmt.Printf("")
}

var GROW_FUNCS = []func(*MoveRequest, []string) []string{
	FilterPossibleMoves,
	FilterMovesVsSpace,
	FilterKillArea,
	FilterTail,
	FilterClosestFoodDirections,
	FilterMinimizeSpace,
}

var SPACE_SAVING_FUNCS = []func(*MoveRequest, []string) []string{
	FilterPossibleMoves,
	FilterTail,
	FilterKillArea,
	FilterMovesVsSpace,
	FilterMinimizeSpace,
	FilterKeyArea,
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// A file for all of the filtering of moves

// not necessairily the best move but the move that we are going with
func bestMoves(data *MoveRequest) ([]string, error) {
	moves := []string{UP, DOWN, LEFT, RIGHT}

	funcArray := GROW_FUNCS

	if data.MetaData.tightSpace || data.NoFood() {
		funcArray = SPACE_SAVING_FUNCS
	}

	for _, filt := range funcArray {
		moves = filt(data, moves)
		if len(moves) == 0 {
			return []string{}, errors.New(
				fmt.Sprintf(
					"0 results returned from %v", GetFunctionName(filt)))
		}
	}
	return moves, nil
}

func GetPossibleDeath(data *MoveRequest, direc string, turns int) int {
	baseData := []*StaticData{}
	for _, snake := range data.Snakes {
		sd := fullStats(snake.HeadPoint, data)
		baseData = append(baseData, sd)
	}

	err := MoveSnakeForward(data.MyIndex, data, direc)
	if err != nil {
		return 0
	}
	newData := []*StaticData{}
	for _, snake := range data.Snakes {
		sd := fullStats(snake.HeadPoint, data)
		newData = append(newData, sd)
	}

	for i := range newData {
		if newData[i].Moves < baseData[i].Moves {
			err = MoveSnakeBackward(data.MyIndex, data)
			if err != nil {
				return 0
			}
			return 1

		}
	}

	err = MoveSnakeForward(data.MyIndex, data, direc)
	if err != nil {
		return 0
	}
	newData = []*StaticData{}
	for _, snake := range data.Snakes {
		sd := fullStats(snake.HeadPoint, data)
		newData = append(newData, sd)
	}
	for i := range newData {
		if newData[i].Moves < baseData[i].Moves {
			err = MoveSnakeBackward(data.MyIndex, data)
			err = MoveSnakeBackward(data.MyIndex, data)
			if err != nil {
				return 0
			}
			return 1

		}
	}

	err = MoveSnakeBackward(data.MyIndex, data)
	err = MoveSnakeBackward(data.MyIndex, data)
	if err != nil {
		return 0
	}
	return 0
}

func FilterTail(data *MoveRequest, moves []string) []string {
	ret := []string{}
	for _, direc := range moves {
		if data.Direcs[direc].myTail {
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		return moves
	}
	return ret
}

func FilterAggression(data *MoveRequest, moves []string) []string {
	data.GenHazards(data, false)
	//for _, move := range moves {
	deaths := GetPossibleDeath(data, DOWN, 2)
	if deaths > 0 {
		return []string{DOWN}
	}
	//}
	data.GenHazards(data, true)
	return moves
}

func FilterKillArea(data *MoveRequest, moves []string) []string {
	ret := []string{}
	head, err := getMyHead(data)
	if err != nil {
		return []string{}
	}
	for _, direc := range moves {
		// we know this is a valid move because all moves are filterd to be vaild
		// this is the location you are moving to
		p, err := GetPointInDirection(head, direc, data)
		if err != nil {
			return []string{}
		}
		if data.KillZones[p.String()] {
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		return moves
	}
	return ret

}

func FilterKeyArea(data *MoveRequest, moves []string) []string {
	ret := []string{}
	head, err := getMyHead(data)
	if err != nil {
		return []string{}
	}
	for _, direc := range moves {
		// we know this is a valid move because all moves are filterd to be vaild
		// this is the location you are moving to
		p, err := GetPointInDirection(head, direc, data)
		if err != nil {
			return []string{}
		}
		p2 := data.Direcs[direc].minKeySnakePart().pnt
		distFromHead := head.Dist(p2)
		distFromPnt := p.Dist(p2)

		// prefer to move in the opposite direction
		if distFromPnt.X > distFromHead.X || distFromPnt.Y > distFromHead.Y {
			ret = append(ret, direc)
		} else if distFromHead.X > distFromHead.Y {
			if direc == UP || direc == DOWN {
				ret = append(ret, direc)
			}
		} else if distFromHead.X < distFromHead.Y {
			if direc == RIGHT || direc == LEFT {
				ret = append(ret, direc)
			}
		}

	}
	if len(ret) == 0 {
		return moves
	}
	return ret
}

func FilterClosestFoodDirections(data *MoveRequest, moves []string) []string {
	directions := []string{}
	min := math.MaxInt64
	for _, direc := range moves {
		food := findGuaranteedClosestFood(data, direc)
		if food == nil {
			continue
		}
		if food.moves < min {
			directions = []string{}
			directions = append(directions, direc)
			min = food.moves
		} else if food.moves == min {
			directions = append(directions, direc)
		}
	}
	if len(directions) == 0 {
		return moves
	}
	return directions
}

func FilterMinimizeSpace(data *MoveRequest, moves []string) []string {
	ret := []string{}
	head, err := getMyHead(data)
	if err != nil {
		return []string{}
	}

	for _, direc := range moves {
		p, err := GetPointInDirection(head, direc, data)
		if err != nil {
			return []string{}
		}

		// if there are multiple moves to consider
		neighbours, err := GetNumNeighbours(data, p)
		if err != nil {
			return []string{}
		}
		if neighbours <= 2 {
			ret = append(ret, direc)
		}
	}
	// dont return a empty array
	if len(ret) == 0 {
		return moves
	}
	return ret
}

// Filters out moves that will put you into tight places.
func FilterMovesVsSpace(data *MoveRequest, moves []string) []string {
	ret := []string{}
	for _, direc := range moves {
		if data.Direcs[direc].MovesVsSpace > len(data.Food) {
			//fmt.Printf("%v\n", ret)
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		max := math.MinInt64
		for _, direc := range moves {
			if data.Direcs[direc].MovesVsSpace == max {
				ret = append(ret, direc)
			} else if data.Direcs[direc].MovesVsSpace > max {
				ret = []string{direc}
				max = data.Direcs[direc].MovesVsSpace
			}

		}
	}
	return ret
}

func FilterPossibleMoves(data *MoveRequest, directions []string) []string {
	ret := []string{}
	for _, direc := range directions {
		if data.Direcs[direc].TotalMoves > 0 {
			ret = append(ret, direc)
		}
	}
	return ret
}
