package elastic

import (
	"context"
	"fmt"

	"spacecraft/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func ConnectWithElasticSearch(ctx context.Context) (context.Context, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		return ctx, fmt.Errorf("failed to connect to elasticsearch: %v", err)
	}
	return context.WithValue(ctx, domain.ClientKey, es), nil
}
