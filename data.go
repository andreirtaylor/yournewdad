package kaa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

type StaticData struct {
	Food   int
	Snakes int
	Moves  int
}

type MetaData struct {
	// denotes the number of moves until you reach the closest piece of food
	MyLength  int
	Hazards   map[string]bool
	FoodMap   map[string]bool
	SnakeHash map[string]SnakeData
	Direcs    MoveMetaData
}

type MetaDataDirec struct {
	// denotes the number of moves until you reach the closest piece of food
	ClosestFood int
	// totals up your length and the ammount of food in a direction
	// if you would fill up the space make it unlikely to go that direction
	MovesVsSpace int
	// definied by the itoa above
	MovesAway []*StaticData
}

func (m *MetaDataDirec) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("\n{")
	buffer.WriteString(fmt.Sprintf("ClosestFood:%v\n", m.ClosestFood))
	buffer.WriteString(fmt.Sprintf("movesVsSpace:%v\n", m.MovesVsSpace))
	for _, x := range m.MovesAway {
		buffer.WriteString(fmt.Sprintf("%#v\n", x))
	}
	buffer.WriteString("}\n")

	return buffer.String()
}

// returns the meta data for the maximum distance you can travel i.e. the whole board
func (m *MetaDataDirec) moveMax() (*StaticData, error) {
	if len(m.MovesAway) == 0 {
		return nil, errors.New("Array is empty")
	}
	return m.MovesAway[len(m.MovesAway)-1], nil
}

type MoveMetaData map[string]*MetaDataDirec

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

	MetaData
}

func (m *MoveRequest) init() {
	m.SetMyLength(m)
	m.GenHazards(m)
	m.GenFoodMap(m)
	m.GenSnakeHash(m)
}

func getMoveRequestString(req *http.Request) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	return buf.String()
}

func NewMoveRequest(str string) (*MoveRequest, error) {
	res := new(MoveRequest)
	err := json.Unmarshal([]byte(str), res)
	if err != nil {
		return nil, err
	}
	err = GenerateMetaData(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *MetaData) SetMyLength(data *MoveRequest) {
	for _, snake := range data.Snakes {
		if snake.Id == data.You && len(data.You) > 0 {
			m.MyLength = len(snake.Coords)
		}
	}
}

type SnakeData struct {
	id         int
	lengthLeft int
}

func (m *MetaData) GenSnakeHash(data *MoveRequest) {
	m.SnakeHash = make(map[string]SnakeData)
	for _, snake := range data.Snakes {
		for i, coord := range snake.Coords {
			m.SnakeHash[coord.String()] = *&SnakeData{
				id:         i,
				lengthLeft: len(snake.Coords) - i - 1,
			}
		}
	}
}

func (m *MetaData) GenHazards(data *MoveRequest) {
	m.Hazards = make(map[string]bool)
	for _, snake := range data.Snakes {
		if len(snake.Coords) >= m.MyLength && data.You != snake.Id {
			head := snake.Head()
			d := head.Down(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}
			d = head.Up(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}
			d = head.Right(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}
			d = head.Left(data)
			if d != nil {
				m.Hazards[d.String()] = true
			}

		}
		for _, coord := range snake.Coords {
			m.Hazards[coord.String()] = true
		}
	}
}

func (m *MetaData) GenFoodMap(data *MoveRequest) {
	m.FoodMap = make(map[string]bool)
	for _, food := range data.Food {
		m.FoodMap[food.String()] = true
	}
}

type MoveResponse struct {
	Move  string  `json:"move"`
	Taunt *string `json:"taunt,omitempty"`
}

type Snake struct {
	Coords       []Point `json:"coords"`
	HealthPoints int     `json:"health_points"`
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Taunt        string  `json:"taunt"`
}

func getJson(data *MoveRequest) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewGameStartRequest(req *http.Request) (*GameStartRequest, error) {
	decoded := GameStartRequest{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	return &decoded, err
}

func (snake Snake) Head() Point     { return snake.Coords[0] }
func (snake *Snake) String() string { return fmt.Sprintf("%#v", snake) }

// Decode [number, number] JSON array into a Point

// Allows decoding a string or number identifier in JSON
// by removing any surrounding quotes and storing in a string
type Identifier string

func (t *Identifier) UnmarshalJSON(data []byte) error {
	*t = Identifier(strings.Trim(string(data), `"`))
	return nil
}
