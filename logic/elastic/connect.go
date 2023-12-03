package elastic

import (
	"context"
	"fmt"

	"spacecraft/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

// ConnectWithElasticSearch nit naming convention: can just be called Connect or Dial
// major: should return an elasticsearch.Client?
func ConnectWithElasticSearch(ctx context.Context) (context.Context, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, //nit: allow to be passed in as a config param
	})
	if err != nil {
		return ctx, fmt.Errorf("failed to connect to elasticsearch: %v", err)
	}
	return context.WithValue(ctx, domain.ClientKey, es), nil
}
