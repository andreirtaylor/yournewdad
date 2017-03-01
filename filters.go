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
	FilterMinMax,
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
	FilterMinMax,
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

func FilterTail(data *MoveRequest, moves []string) []string {
	ret := []string{}
	for _, direc := range moves {
		if data.Direcs[direc].SeeTail {
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		return moves
	}
	return ret
}

func FilterMinMax(data *MoveRequest, moves []string) []string {
	ret := []string{}
	currStats := GenMinMaxStats(data.minMaxArr)

	for _, move := range moves {
		nextStats := GenMinMaxStats(data.Direcs[move].minMaxArr)
		for key, val := range nextStats.snakes {
			if key != data.MyIndex {
				nextMoves := float64(val.moves)
				currMoves := float64(currStats.snakes[key].moves)
				if 1-nextMoves/currMoves >= 0.3 {
					ret = append(ret, move)
				}
			}
		}
	}
	if len(ret) == 0 {
		return moves
	}
	return ret
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
		p, err := GetPointInDirectionHazards(head, direc, data)
		if err != nil {
			return []string{}
		}
		if p != nil && data.KillZones[p.String()] {
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

		p2 := data.Direcs[direc].KeySnakeData.minKeySnakePart().pnt
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
		if data.Direcs[direc].Moves > 0 {
			head := data.Snakes[data.MyIndex].Head()
			p, _ := GetPointInDirectionHazards(head, direc, data)

			if p != nil && !data.Hazards[p.String()] {
				ret = append(ret, direc)
			}
		}
	}
	if len(ret) == 0 {
		for _, direc := range directions {
			if data.Direcs[direc].Moves > 0 {
				ret = append(ret, direc)
			}
		}
	}
	return ret
}
