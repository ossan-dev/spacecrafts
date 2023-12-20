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
	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			// nit: here we lose the not 200 error context
			return nil, fmt.Errorf("err while reading response body: %w", err)
		}
		return nil, fmt.Errorf("err with the Elasticsearch request processing: %v", string(data))
	}
	var lookupRes domain.LookupResponse
	// a note on correct usage of json.NewDecoder:
	// https://mottaquikarim.github.io/dev/posts/you-might-not-be-using-json.decoder-correctly-in-golang/
	if err = json.NewDecoder(res.Body).Decode(&lookupRes); err != nil {
		return nil, fmt.Errorf("failed to decode the elastic search result: %w", err)
	}
	return lookupRes.Source, nil
}
