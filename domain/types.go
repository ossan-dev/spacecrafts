package domain

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
	SpacecraftClass *SpacecraftInfo `json:"spacecraftClass"`
	Owner           *SpacecraftInfo `json:"owner"`
	Operator        *SpacecraftInfo `json:"operator"`
}

type SpacecraftWrapper struct {
	PageNumber       int           `json:"pageNumber"`
	PageSize         int           `json:"pageSize"`
	NumberOfElements int           `json:"numberOfElements"`
	TotalPages       int           `json:"totalPages"`
	TotalElements    int           `json:"totalElements"`
	Data             []*Spacecraft `json:"data"`
}
