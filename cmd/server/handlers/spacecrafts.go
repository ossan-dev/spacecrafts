package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"esdemov8/domain"
)

var Spacecrafts []*domain.Spacecraft

func extractIntQueryParam(r *http.Request, paramName string, defaultValue int) (*int, error) {
	urlValues := r.URL.Query()
	var result int
	rawParam := urlValues.Get(paramName)
	if rawParam == "" {
		return &defaultValue, nil
	}
	result, err := strconv.Atoi(rawParam)
	if err != nil {
		return nil, fmt.Errorf("err while fetching value from query string: %v", err)
	}
	return &result, nil
}

func GetSpacecrafts(w http.ResponseWriter, r *http.Request) {
	pageSize, err := extractIntQueryParam(r, "pageSize", 100)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	pageNumber, err := extractIntQueryParam(r, "pageNumber", 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	low := *pageNumber * *pageSize
	high := low + *pageSize
	var result []*domain.Spacecraft
	if high >= len(Spacecrafts) {
		result = Spacecrafts[low:]
		fmt.Println("high >= len(Spacecrafts):", len(result), "spacecrafts")
	} else {
		result = Spacecrafts[low:high]
		fmt.Println("else:", len(result), "spacecrafts")
	}
	data, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	time.Sleep(2 * time.Second)
	w.Write([]byte(data))
}
