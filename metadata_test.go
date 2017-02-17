package kaa

import (
	"fmt"
	"testing"
)

func TestMetaData(t *testing.T) {
	req := &MoveRequest{
		GameId: "d0bd244e-91da-4e63-86e6-ea575376c3be",
		Height: 20,
		Width:  20,
		Turn:   4,
		Food: []Point{
			Point{X: 5, Y: 13},
		},
		Snakes: []Snake{
			Snake{
				Coords: []Point{
					Point{X: 14, Y: 12},
					Point{X: 13, Y: 12},
					Point{X: 13, Y: 13},
				},
				HealthPoints: 96,
				Id:           "639fb7cd-2590-4418-abcc-3da577559fc6",
				Name:         "d0bd244e-91da-4e63-86e6-ea575376c3be (20x20)",
				Taunt:        "639fb7cd-2590-4418-abcc-3da577559fc6",
			},
		},
		You: "639fb7cd-2590-4418-abcc-3da577559fc6",
	}

	data, _ := GenerateMetaData(req)
	fmt.Printf("%#v", data)
}
