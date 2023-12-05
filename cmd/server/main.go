package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"

	"spacecraft/internal/handlers"
)

//go:embed spacecraft.json
var data []byte

func main() {
	if err := json.Unmarshal(data, &handlers.Spacecraft); err != nil {
		log.Fatalf("err while unmarshaling file: %v", err)
	}
	r := http.NewServeMux()
	r.HandleFunc("/spacecraft", handlers.GetSpacecraft)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to launch web server: %v", err)
	}
}
