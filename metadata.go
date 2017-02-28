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
