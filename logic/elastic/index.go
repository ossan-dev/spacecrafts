package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"spacecraft/domain"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func IndexSpacecraftAsDocuments(ctx context.Context) error {
	spacecrafts := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)
	for spacecraftID, spacecraft := range spacecrafts {
		res, err := client.Index("spacecrafts", esutil.NewJSONReader(spacecraft), client.Index.WithDocumentID(strconv.Itoa(spacecraftID)))
		if err == nil {
			defer res.Body.Close()
			fmt.Println(res)
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func DeleteIndex(ctx context.Context, indexName string) error {
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("failed to get info for index %q with err: %v", indexName, err)
	}
	if res.StatusCode != http.StatusNotFound {
		res, err := client.Indices.Delete([]string{indexName})
		if err != nil {
			return fmt.Errorf("failed to delete index %q with err: %v", indexName, err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			data, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("err while reading response body: %v", err)
			}
			return fmt.Errorf("failed to delete index %q with message: %v", indexName, string(data))
		}
	}
	return nil
}

func IndexSpacecraftAsDocumentsAsync(ctx context.Context) error {
	spacecrafts := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      "spacecrafts",
		Client:     client,
		NumWorkers: 5,
	})
	if err != nil {
		return fmt.Errorf("err while creating bulk indexer: %v", err)
	}
	for spacecraftID, spacecraft := range spacecrafts {
		data, err := json.Marshal(spacecraft)
		if err != nil {
			return fmt.Errorf("err while marshaling object: %v", spacecraft)
		}
		err = bulkIndexer.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: strconv.Itoa(spacecraftID),
			Body:       strings.NewReader(string(data)),
		})
		if err != nil {
			return fmt.Errorf("failed to add %v to the bulk indexer: %v", spacecraft, err)
		}
	}
	if err = bulkIndexer.Close(ctx); err != nil {
		return fmt.Errorf("failed to close the bulk indexer: %v", err)
	}
	stats := bulkIndexer.Stats()
	fmt.Fprintf(os.Stdout, "Spacecrafts indexed on Elasticsearch: %d\n", stats.NumIndexed)
	return nil
}
