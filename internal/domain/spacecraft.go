package domain

import (
	"encoding/json"
	"fmt"
	"os"
)

type SpacecraftInfo struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

type Spacecraft struct {
	Uid             string          `json:"uid"`
	Name            string          `json:"name"`
	Registry        string          `json:"registry"`
	Status          string          `json:"status"`
	DateStatus      string          `json:"dateStatus"`
	SpacecraftClass *SpacecraftInfo `json:"spacecraftClass,omitempty"`
	Owner           *SpacecraftInfo `json:"owner,omitempty"`
	Operator        *SpacecraftInfo `json:"operator,omitempty"`
}

func GetPaginatedSpacecrafts(spacecrafts []*Spacecraft, pageSize, pageNumber int) (res SpacecraftWrapper) {
	low := pageNumber * pageSize
	high := low + pageSize
	res.PageNumber = pageNumber
	res.PageSize = pageSize
	res.NumberOfElements = pageSize
	res.TotalPages = (len(spacecrafts) / pageSize) + 1
	res.TotalElements = len(spacecrafts)
	if high >= len(spacecrafts) {
		res.NumberOfElements = len(spacecrafts[low:])
		res.Data = spacecrafts[low:]
	} else {
		res.Data = spacecrafts[low:high]
	}
	return
}

func (s Spacecraft) String() {
	data, _ := json.MarshalIndent(s, "", "\t")
	fmt.Fprintln(os.Stdout, string(data))
}

type SpacecraftWrapper struct {
	PageNumber       int           `json:"pageNumber"`
	PageSize         int           `json:"pageSize"`
	NumberOfElements int           `json:"numberOfElements"`
	TotalPages       int           `json:"totalPages"`
	TotalElements    int           `json:"totalElements"`
	Data             []*Spacecraft `json:"data"`
}
