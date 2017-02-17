package kaa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const ( // iota is reset to 0
	UP    = "up"    // c0 == 0
	DOWN  = "down"  // c1 == 1
	LEFT  = "left"  // c2 == 2
	RIGHT = "right" // c2 == 2
)

type GameStartRequest struct {
	GameId string `json:"game_id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type GameStartResponse struct {
	Color   string  `json:"color"`
	HeadUrl *string `json:"head_url,omitempty"`
	Name    string  `json:"name"`
	Taunt   *string `json:"taunt,omitempty"`
}

type MetaData struct {
	Food   int
	Snakes int
	Moves  int
	score  float64
}

type MoveRequest struct {
	// static
	GameId string `json:"game_id"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Turn   int    `json:"turn"`

	// dynamic
	Food   []Point `json:"food"`
	Snakes []Snake `json:"snakes"`
	You    string  `json:"you"`

	// added by me
	// lists all the points that are hazards this turn
	Hazards map[string]bool
}

func (m *MoveRequest) GenHazards() {
	m.Hazards = make(map[string]bool)
	for _, snake := range m.Snakes {
		for _, coord := range snake.Coords {
			m.Hazards[coord.String()] = true
		}
	}
}

type MoveResponse struct {
	Move  string  `json:"move"`
	Taunt *string `json:"taunt,omitempty"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Coords       []Point `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
}

func NewMoveRequest(req *http.Request) (*MoveRequest, error) {
	decoded := MoveRequest{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	decoded.GenHazards()
	return &decoded, err
}

func NewGameStartRequest(req *http.Request) (*GameStartRequest, error) {
	decoded := GameStartRequest{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	return &decoded, err
}

func (snake Snake) Head() Point { return snake.Coords[0] }

// Decode [number, number] JSON array into a Point
func (point *Point) UnmarshalJSON(data []byte) error {
	var coords []int
	json.Unmarshal(data, &coords)
	if len(coords) != 2 {
		return errors.New("Bad set of coordinates: " + string(data))
	}
	*point = Point{X: coords[0], Y: coords[1]}
	return nil
}

func (point *Point) String() string {
	return fmt.Sprintf("{%d,%d}", point.X, point.Y)
}

// directional functions return a new point or nil if the point is out of the
// board
func (point *Point) Up(data *MoveRequest) *Point {
	if point.Y == 0 {
		return nil
	}
	ret := &Point{point.X, point.Y - 1}
	if data.Hazards[ret.String()] {
		return nil
	}
	return ret
}

func (point *Point) Down(data *MoveRequest) *Point {
	if point.Y == data.Height-1 {
		return nil
	}
	ret := &Point{point.X, point.Y + 1}
	if data.Hazards[ret.String()] {
		return nil
	}
	return ret
}

func (point *Point) Left(data *MoveRequest) *Point {
	if point.X == 0 {
		return nil
	}
	ret := &Point{point.X - 1, point.Y}
	if data.Hazards[ret.String()] {
		return nil
	}
	return ret
}

func (point *Point) Right(data *MoveRequest) *Point {
	if point.X == data.Width-1 {
		return nil
	}
	ret := &Point{point.X + 1, point.Y}
	if data.Hazards[ret.String()] {
		return nil
	}
	return ret
}

// Allows decoding a string or number identifier in JSON
// by removing any surrounding quotes and storing in a string
type Identifier string

func (t *Identifier) UnmarshalJSON(data []byte) error {
	*t = Identifier(strings.Trim(string(data), `"`))
	return nil
}
