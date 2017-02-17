package kaa

import (
	"net/http"
)

func init() {
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/move", handleMove)
}
