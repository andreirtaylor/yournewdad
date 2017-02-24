package kaa

import (
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

func respond(res http.ResponseWriter, obj interface{}) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(obj)
}

func handleStart(res http.ResponseWriter, req *http.Request) {
	color := "gold"
	if appengine.IsDevAppServer() {
		color = "gold"
	}
	respond(res, GameStartResponse{
		Taunt:   toStringPointer("Dad 2.0 Ready"),
		Color:   color,
		Name:    "Your New Dad",
		HeadUrl: toStringPointer("http://i.imgur.com/MLo4AQI.png"),
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
			Taunt: "can't parse this!",
		})
		return
	}

	// its tooooooooo sloooooooooooooooww :(
	//SaveMove(data, req)

	move, err := getMove(data, req)

	if err != nil {
		respond(res, MoveResponse{
			Move:  "up",
			Taunt: "Couldn't parse",
		})
		log.Errorf(ctx, "Could not find a move for this data")
		return
	}
	taunt := getTaunt(data.Turn)
	respond(res, MoveResponse{
		Move:  move,
		Taunt: taunt,
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
