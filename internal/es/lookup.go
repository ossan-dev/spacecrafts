package es

import (
	"encoding/json"
	"fmt"
	"io"

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func GetByID(esClient *elasticsearch.Client, index, documentID string) (*domain.Spacecraft, error) {
	res, err := esClient.Get(index, documentID)
	if err != nil {
		return nil, fmt.Errorf("err while looking up the document with ID: %v with err: %w", documentID, err)
	}
	defer res.Body.Close()
	var lookupRes domain.LookupResponse
	jsonDecoder := json.NewDecoder(res.Body)
	for {
		if err := jsonDecoder.Decode(&lookupRes); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("err while reading the stream: %w", err)
		}
	}
	if lookupRes.Found {
		return lookupRes.Source, nil
	}
	return nil, fmt.Errorf("document with id: %q not found", documentID)
}
