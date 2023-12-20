package es

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func Connect(url string) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	return es, nil
}
