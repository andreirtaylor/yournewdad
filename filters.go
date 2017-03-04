package main

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
	FilterTieAreas,
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

var AGGRESSION = []func(*MoveRequest, []string) []string{
	FilterPossibleMoves,
	FilterMovesVsSpace,
	FilterMinMax,
	FilterTieAreas,
	FilterKillArea,
	FilterTail,
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func imAgressive(data *MoveRequest) bool {
	if !data.MetaData.tightSpace &&
		data.Snakes[data.MyIndex].HealthPoints > 80 &&
		data.MyLength > 2*data.Width+len(data.Food) {
		return true
	}
	return false
}

// A file for all of the filtering of moves

// not necessairily the best move but the move that we are going with
func bestMoves(data *MoveRequest) ([]string, error) {
	moves := []string{UP, DOWN, LEFT, RIGHT}

	funcArray := GROW_FUNCS

	if data.MetaData.tightSpace || data.NoFood() {
		funcArray = SPACE_SAVING_FUNCS
	}

	if imAgressive(data) {
		funcArray = AGGRESSION
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

func FilterTieAreas(data *MoveRequest, moves []string) []string {
	ret := []string{}

	for _, move := range moves {
		head := data.Snakes[data.MyIndex].Head()

		p, _ := GetPointInDirection(head, move, data)

		if p == nil {
			continue
		}
		if data.minMaxArr[p.Y][p.X].tie {
			ret = append(ret, move)
		}
	}
	if len(ret) == 0 {
		return moves
	}
	return ret
}

func FilterMinMax(data *MoveRequest, moves []string) []string {
	ret := []string{}

	lossThreshold := 0.3
	if imAgressive(data) {
		lossThreshold = 0.1
	}
	currStats := data.MinMaxMD

	for _, move := range moves {
		nextStats := data.Direcs[move].MinMaxMD
		for key, val := range nextStats.snakes {
			if key != data.MyIndex {
				nextMoves := float64(val.moves)
				currMoves := float64(currStats.snakes[key].moves)
				if 1-nextMoves/currMoves >= lossThreshold {
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

func FilterExpandArea(data *MoveRequest, moves []string) []string {
	ret := []string{}
	currStats := GenMinMaxStats(data.minMaxArr)

	for _, move := range moves {
		nextStats := GenMinMaxStats(data.Direcs[move].minMaxArr)
		// if I increase my area by going this way go that way
		if len(currStats.movesHash) < len(nextStats.movesHash) {
			ret = append(ret, move)
		}
	}

	if len(ret) == 0 {
		return moves
	}
	return ret
	return moves
}
func FilterKeyArea(data *MoveRequest, moves []string) []string {
	ret := []string{}
	head, err := getMyHead(data)
	if err != nil {
		return []string{}
	}
	maxDist := 0
	for _, direc := range moves {
		// we know this is a valid move because all moves are filterd to be vaild
		// this is the location you are moving to
		p, err := GetPointInDirection(head, direc, data)
		if err != nil {
			return []string{}
		}

		// if there isno key snake part in your area its a fine move
		if data.KSD.minKeySnakePart() == nil {
			ret = append(ret, direc)
			continue
		}
		p2 := data.KSD.minKeySnakePart().pnt
		distFromHead := head.Dist(p2)
		if p == nil {
			continue
		}
		distFromPnt := p.Dist(p2)

		// prefer to move in the opposite direction
		//fmt.Printf("%v %v %v %v %v %v %v\n", distFromPnt.X, distFromHead.X, distFromPnt.Y, distFromHead.Y, p2, direc, maxDist)
		if distFromPnt.X < maxDist && distFromPnt.Y < maxDist {
			continue
		}

		// prefer to move in the opposite direction
		if distFromPnt.X > distFromHead.X {
			maxDist = distFromPnt.X
			ret = append(ret, direc)
		} else if distFromPnt.Y > distFromHead.Y {
			maxDist = distFromPnt.Y
			ret = append(ret, direc)
		} else if distFromHead.X > distFromHead.Y {
			if direc == UP || direc == DOWN {
				ret = append(ret, direc)
				maxDist = distFromPnt.X
			}
		} else if distFromHead.X < distFromHead.Y {
			if direc == RIGHT || direc == LEFT {
				ret = append(ret, direc)
				maxDist = distFromPnt.Y
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
