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

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	respond(res, GameStartResponse{
		Taunt:   toStringPointer("battlesnake-go!"),
		Color:   "#00FF00",
		Name:    fmt.Sprintf("%v (%vx%v)", data.GameId, data.Width, data.Height),
		HeadUrl: toStringPointer(fmt.Sprintf("%v://%v/static/head.png", scheme, req.Host)),
	})
}

func handleMove(res http.ResponseWriter, req *http.Request) {
	data, err := NewMoveRequest(req)
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

	if err != nil {
		respond(res, MoveResponse{
			Move:  "up",
			Taunt: toStringPointer("can't parse this!"),
		})
		return
	}

	respond(res, MoveResponse{
		Move:  move,
		Taunt: &data.You,
	})
}

func getMove(data *MoveRequest, req *http.Request) (string, error) {
	ctx := appengine.NewContext(req)

	metadata, err := GenerateMetaData(data)
	if err != nil {
		log.Errorf(ctx, "generating MetaData: %v", err)
		return "", err
	}

	move := bestMove(metadata)
	//moves := bestMoves(metadata)

	//log.Infof(ctx, "%v", moves)
	//for direc, direcData := range metadata {
	//	log.Infof(ctx, "Meta data %v\ndirec %#v", direcData, direc)
	//}
	return move, err
}
