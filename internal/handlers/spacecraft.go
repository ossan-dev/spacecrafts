package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"spacecraft/internal/domain"
)

var Spacecraft []*domain.Spacecraft

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

func GetSpacecraft(w http.ResponseWriter, r *http.Request) {
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
	result := domain.GetPaginatedSpacecrafts(Spacecraft, *pageSize, *pageNumber)
	data, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	time.Sleep(time.Second)
	w.Write(data)
}
