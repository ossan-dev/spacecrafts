package domain

type LookupResponse struct {
	Index   string      `json:"_index"`
	ID      string      `json:"id"`
	Version int         `json:"_version"`
	Found   bool        `json:"found"`
	Source  *Spacecraft `json:"_source"`
}

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source *Spacecraft `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
