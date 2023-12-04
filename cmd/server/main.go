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
	handlers.Spacecraft, err = store.LoadspacecraftFromFile("spacecraft.json")
	if err != nil {
		// nit: why don't use log/slog package? explicit writer here in main.go is not needed (main can't be tested)
		fmt.Fprintf(os.Stderr, "err while fetching spacecraft: %v", err)
		return
	}
	r := http.NewServeMux()
	r.HandleFunc("/spacecraft", handlers.Getspacecraft)
	// nit:
	// use log.Fatal(http.ListenAndServe(":7000", r)) to save few lines and be more idiomatic
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Fprintf(os.Stderr, "failed to launch web server: %v", err)
		return
	}
}
