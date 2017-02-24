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
	ClosestFoodDirections,
}

var SPACE_SAVING_FUNCS = []func(*MoveRequest, []string) []string{
	FilterPossibleMoves,
	FilterMovesVsSpace,
	FilterMinimizeSpace,
}

func GetFunctionName(i interface{}) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if len(funcName) < 29 {
		return funcName
	}

	return funcName[29:]
}

// A file for all of the filtering of moves

// not necessairily the best move but the move that we are going with
func bestMoves(data *MoveRequest) ([]string, error) {
	moves := []string{UP, DOWN, LEFT, RIGHT}

	funcArray := GROW_FUNCS

	if data.MetaData.tightSpace {
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

func FilterMinimizeSpace(data *MoveRequest, moves []string) []string {
	min := math.MaxInt64
	ret := ""
	head, err := getMyHead(data)
	if err != nil {
		return []string{}
	}
	for _, direc := range moves {
		p, err := GetPointInDirection(&head, direc, data)
		if err != nil {
			return []string{}
		}

		neighbours, err := GetNumNeighbours(data, p)
		if err != nil {
			return []string{}
		}
		if neighbours < min {
			ret = direc
			min = neighbours
		}
	}
	return []string{ret}
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
