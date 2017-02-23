package kaa

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

func respond(res http.ResponseWriter, obj interface{}) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(obj)
}

func handleStart(res http.ResponseWriter, req *http.Request) {
	data, err := NewGameStartRequest(req)

	//saveGame(data, req)

	if err != nil {
		respond(res, GameStartResponse{
			Taunt:   toStringPointer("battlesnake-go!"),
			Color:   "#00FF00",
			Name:    fmt.Sprintf("%v (%vx%v)", data.GameId, data.Width, data.Height),
			HeadUrl: toStringPointer(fmt.Sprintf("%v://%v/static/head.png")),
		})
	}

	color := "gold"
	if appengine.IsDevAppServer() {
		color = "gold"
	}

	respond(res, GameStartResponse{
		Taunt:   toStringPointer("battlesnake-go!"),
		Color:   color,
		Name:    fmt.Sprintf("%v (%vx%v)", data.GameId, data.Width, data.Height),
		HeadUrl: toStringPointer("https://media.giphy.com/media/I2v9aehFlQBQ4/giphy.gif"),
	})
}

func handleMove(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	str := getMoveRequestString(req)

	// log each and every json blob that comes in
	log.Infof(ctx, str)

	data, err := NewMoveRequest(str)
	if err != nil {
		respond(res, MoveResponse{
			Move:  "up",
			Taunt: toStringPointer("can't parse this!"),
		})
		return
	}

	// its tooooooooo sloooooooooooooooww :(
	//SaveMove(data, req)

	move, err := getMove(data, req)
	log.Infof(ctx, "%v", data.tightSpace)

	if err != nil {
		respond(res, MoveResponse{
			Move:  "up",
			Taunt: toStringPointer("can't parse this!"),
		})
		log.Errorf(ctx, "Could not find a move for this data")
		return
	}

	respond(res, MoveResponse{
		Move:  move,
		Taunt: &data.You,
	})
}

func getMove(data *MoveRequest, req *http.Request) (string, error) {
	ctx := appengine.NewContext(req)

	move, err := bestMove(data)

	if err != nil {
		log.Errorf(ctx, "generating MetaData: %v", err)
		return "", err
	}
	log.Infof(ctx, "%v\n", move)
	return move, nil
}
