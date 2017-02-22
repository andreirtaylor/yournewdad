package kaa

import (
	"fmt"
	"math"
)

func keepFMTFilters() {
	fmt.Printf("")
}

// A file for all of the filtering of moves

// not necessairily the best move but the move that we are going with
func bestMoves(data *MoveRequest) ([]string, error) {
	moves, err := FilterPossibleMoves(data)
	if err != nil {
		return nil, err
	}
	moves = FilterMovesVsSpace(data, moves)
	moves = ClosestFoodDirections(data, moves)
	return moves, nil
}

func FilterMinimizeSpace(data *MoveRequest, moves []string) (string, error) {
	min := math.MaxInt64
	ret := ""
	head, err := getMyHead(data)
	if err != nil {
		return "", err
	}
	for _, direc := range moves {
		p, err := GetPointInDirection(&head, direc, data)
		if err != nil {
			return "", err
		}

		neighbours, err := GetNumNeighbours(data, p)
		if err != nil {
			return "", err
		}
		if neighbours < min {
			ret = direc
			min = neighbours
		}
	}
	return ret, nil
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
			if data.Direcs[direc].MovesVsSpace < max {
				ret = []string{direc}
			} else if data.Direcs[direc].MovesVsSpace < max {
				ret = append(ret, direc)
			}
		}
	}
	return ret
}

func FilterPossibleMoves(data *MoveRequest) ([]string, error) {
	directions := []string{UP, DOWN, LEFT, RIGHT}
	ret := []string{}
	for _, direc := range directions {
		if len(data.Direcs[direc].MovesAway) != 0 {
			ret = append(ret, direc)
		}
	}
	return ret, nil
}
