package main

import (
	"fmt"
	"net/http"
	"os"

	"esdemov8/cmd/server/handlers"
	"esdemov8/cmd/server/store"
)

func main() {
	var err error
	handlers.Spacecrafts, err = store.LoadSpacecraftsFromFile("store/spacecrafts.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err while fetching spacecrafts: %v", err)
		return
	}
	r := http.NewServeMux()
	r.HandleFunc("/spacecrafts", handlers.GetSpacecrafts)
	http.ListenAndServe(":8080", r)
}