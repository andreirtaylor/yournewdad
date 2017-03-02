package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/move", handleMove)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	log.Printf("Running server on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
