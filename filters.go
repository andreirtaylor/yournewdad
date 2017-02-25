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
	FilterClosestFoodDirections,
}

var SPACE_SAVING_FUNCS = []func(*MoveRequest, []string) []string{
	FilterPossibleMoves,
	FilterMovesVsSpace,
	FilterMinimizeSpace,
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
func FilterClosestFoodDirections(data *MoveRequest, moves []string) []string {
	directions := []string{}
	min := math.MaxInt64
	metaD := data.Direcs
	for _, direc := range moves {
		if data.Direcs[direc].ClosestFood == 0 {
			continue
		}
		if metaD[direc].ClosestFood < min {
			directions = []string{}
			directions = append(directions, direc)
			min = metaD[direc].ClosestFood
		} else if metaD[direc].ClosestFood == min {
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
		p, err := GetPointInDirection(&head, direc, data)
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
		if data.Direcs[direc].MovesVsSpace > -2 {
			//fmt.Printf("%v\n", ret)
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		max := math.MinInt64
		for _, direc := range moves {
			if data.Direcs[direc].MovesVsSpace > max {
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
