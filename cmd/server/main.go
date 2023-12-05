package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"spacecraft/cmd/server/handlers"
	"spacecraft/domain"
)

//go:embed store/spacecraft.json
var data []byte

func main() {
	var res []*domain.Spacecraft
	var err error
	if err := json.Unmarshal(data, &res); err != nil {
		fmt.Fprintf(os.Stderr, "err while unmarshaling file: %v", err)
		return
	}
	handlers.Spacecraft = res
	if err != nil {
		// nit: why don't use log/slog package? explicit writer here in main.go is not needed (main can't be tested)
		// [x]: switch to log
		fmt.Fprintf(os.Stderr, "err while fetching spacecraft: %v", err)
		return
	}
	r := http.NewServeMux()
	r.HandleFunc("/spacecraft", handlers.GetSpacecraft)
	// nit:
	// [x]: check return values of log.Fatal (it might write `nil` on the shell)
	// use log.Fatal(http.ListenAndServe(":7000", r)) to save few lines and be more idiomatic
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Fprintf(os.Stderr, "failed to launch web server: %v", err)
		return
	}
}
