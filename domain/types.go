package domain

import (
	"encoding/json"
	"fmt"
	"io"
)

/****************************************************/
/************* Internal for Context *****************/
/****************************************************/
type contextKey struct {
	Key int
}

var (
	ClientKey contextKey = contextKey{Key: 1}
	ModelsKey contextKey = contextKey{Key: 2}
)

/****************************************************/
/************* Domain Types *************************/
/****************************************************/
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

func (s Spacecraft) Print(fs io.Writer) {
	data, _ := json.MarshalIndent(s, "", "\t")
	fmt.Fprintln(fs, string(data))
}

type SpacecraftWrapper struct {
	PageNumber       int           `json:"pageNumber"`
	PageSize         int           `json:"pageSize"`
	NumberOfElements int           `json:"numberOfElements"`
	TotalPages       int           `json:"totalPages"`
	TotalElements    int           `json:"totalElements"`
	Data             []*Spacecraft `json:"data"`
}

/****************************************************/
/************* Elasticsearch Types ******************/
/****************************************************/
type LookupResponse struct {
	Index   string      `json:"_index"`
	ID      string      `json:"id"`
	Version int         `json:"_version"`
	Source  *Spacecraft `json:"_source"`
}
