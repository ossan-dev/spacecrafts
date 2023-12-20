package es

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func GetByID(esClient *elasticsearch.Client, index, documentID string) (*domain.Spacecraft, error) {
	res, err := esClient.Get(index, documentID)
	if err != nil {
		return nil, fmt.Errorf("err while looking up the document with ID: %v with err: %w", documentID, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("err while reading response body: %w", err)
	}
	if res.StatusCode == http.StatusOK {
		var lookupRes domain.LookupResponse
		if err = json.Unmarshal(data, &lookupRes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the elastic search result: %w", err)
		}
		return lookupRes.Source, nil
	}
	return nil, fmt.Errorf("failed to lookup doc: %w", fmt.Errorf(string(data)))
}
