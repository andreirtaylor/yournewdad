package kaa

import (
	"math"
)

// A file for all of the filtering of moves

// not necessairily the best move but the move that we are going with
func bestMoves(metaD map[string]*MetaData) ([]string, error) {
	moves, err := FilterPossibleMoves(metaD)
	if err != nil {
		return nil, err
	}
	moves = FilterMovesVsSpace(metaD, moves)
	moves = ClosestFoodDirections(metaD, moves)
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
func FilterMovesVsSpace(metaD map[string]*MetaData, moves []string) []string {
	ret := []string{}
	for _, direc := range moves {
		if metaD[direc].MovesVsSpace > -2 {
			//fmt.Printf("%v\n", ret)
			ret = append(ret, direc)
		}
	}
	if len(ret) == 0 {
		max := math.MinInt64
		for _, direc := range moves {
			if metaD[direc].MovesVsSpace < max {
				ret = []string{direc}
			} else if metaD[direc].MovesVsSpace < max {
				ret = append(ret, direc)
			}
		}
	}
	return ret
}

func FilterPossibleMoves(metaD map[string]*MetaData) ([]string, error) {
	directions := []string{UP, DOWN, LEFT, RIGHT}
	ret := []string{}
	for _, direc := range directions {
		if len(metaD[direc].MovesAway) == 0 {
			ret = append(ret, direc)
		}
	}
	return directions, nil
}
