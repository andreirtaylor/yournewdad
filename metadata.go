package kaa

import (
	"container/heap"
	"fmt"
)

func keepMDFMT() {
	fmt.Printf("")
}

// This file contains all of the functions which build
// the metadata in each direction

func pushOntoPQ(
	p *Point,
	seen map[string]bool,
	pq *PriorityQueue,
	priority int) {

	if p != nil {
		if !seen[p.String()] {
			seen[p.String()] = true

			heap.Push(pq, &Item{
				value:    *p,
				priority: priority + 1,
			})
		}
	}
}

func GenerateMetaData(data *MoveRequest) error {
	data.init()
	data.Direcs = make(MoveMetaData)
	data.Direcs[UP] = &MetaDataDirec{}
	data.Direcs[DOWN] = &MetaDataDirec{}
	data.Direcs[LEFT] = &MetaDataDirec{}
	data.Direcs[RIGHT] = &MetaDataDirec{}

	tightSpace := true
	for direc, direcMD := range data.Direcs {
		head, err := getMyHead(data)
		if err != nil {
			return err
		}
		newHead, err := GetPointInDirection(head, direc, data)
		if err != nil {
			return err
		}
		stats := fullStats(newHead, data)
		//fmt.Printf("%#v\n%#v\n", sd, direc)
		if stats != nil {
			ksd := stats.KeySnakeData.minKeySnakePart()
			if ksd != nil {
				direcMD.MovesVsSpace = stats.Moves - stats.Food - ksd.lengthLeft
			}
			direcMD.KeySnakeData = stats.KeySnakeData
			direcMD.TotalMoves = stats.Moves
			direcMD.TotalFood = stats.Food
			direcMD.sortedFood = stats.sortedFood
			direcMD.FoodHash = stats.FoodHash
			direcMD.myTail = stats.SeeTail
			if direcMD.MovesVsSpace > 20 {
				tightSpace = false
			}
		}
	}
	data.tightSpace = tightSpace
	return nil
}
