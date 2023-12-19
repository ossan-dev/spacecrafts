package es

import (
	"context"
	"fmt"

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

// major: should return an elasticsearch.Client?
func Connect(ctx context.Context, url string) (context.Context, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return ctx, fmt.Errorf("failed to connect to elasticsearch: %v", err)
	}
	return context.WithValue(ctx, domain.ClientKey, es), nil
}
