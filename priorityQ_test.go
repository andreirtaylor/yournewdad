package main

import (
	"container/heap"
	"testing"
)

func TestPQ(t *testing.T) {
	// Some items and their priorities.
	points := []Point{
		Point{1, 1},
		Point{2, 1},
		Point{3, 2},
		Point{4, 2},
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, 4)

	// Insert a new item and then modify its priority.
	for i, p := range points {
		pq[i] = &Item{
			value:    p,
			priority: p.Y,
			index:    i,
		}
	}

	heap.Init(&pq)

	heap.Push(&pq, &Item{
		value:    Point{4, 3},
		priority: 3,
	})

	// Take the items out; they arrive in decreasing priority order.
	for _, val := range []int{1, 1, 2, 2, 3} {
		item := heap.Pop(&pq).(*Item)
		if val != item.priority {
			t.Errorf("Expected priority to be %d got %d", val, item.priority)
		}
	}

}
