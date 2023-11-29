package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"spacecraft/domain"
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

func Getspacecraft(w http.ResponseWriter, r *http.Request) {
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
	var result domain.SpacecraftWrapper
	result.PageNumber = *pageNumber
	result.PageSize = *pageSize
	result.NumberOfElements = *pageSize
	result.TotalPages = (len(Spacecraft) / *pageSize) + 1
	result.TotalElements = len(Spacecraft)
	if high >= len(Spacecraft) {
		result.NumberOfElements = len(Spacecraft[low:])
		result.Data = Spacecraft[low:]
	} else {
		result.Data = Spacecraft[low:high]
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
