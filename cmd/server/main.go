package main

import (
	"fmt"
	"net/http"
	"os"

	"spacecraft/cmd/server/handlers"
	"spacecraft/cmd/server/store"
)

func main() {
	var err error
	handlers.Spacecraft, err = store.LoadspacecraftFromFile("store/spacecraft.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err while fetching spacecraft: %v", err)
		return
	}
	r := http.NewServeMux()
	r.HandleFunc("/spacecraft", handlers.Getspacecraft)
	http.ListenAndServe(":8080", r)
}
