package logic

import (
	"context"

	"esdemov8/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func ConnectWithElasticSearch(ctx context.Context) context.Context {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		panic(err)
	}
	return context.WithValue(ctx, domain.ClientKey, es)
}
