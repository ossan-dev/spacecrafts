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

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func IndexSpacecraftAsDocuments(ctx context.Context) error {
	// Major: do not store non-scoped request "objects" in the context.
	//  I want to explicitly state that context.Value() should NEVER be used for values that are not created and destroyed
	// during the lifetime of the request.
	// You shouldn’t store a logger there if it isn’t created specifically to be scoped to this request,
	// and likewise you shouldn’t store a generic database connection in a context value.
	// see: https://www.calhoun.io/pitfalls-of-context-values-and-how-to-avoid-or-mitigate-them/
	// or
	// https://pkg.go.dev/context
	// Package context defines the Context type, which carries deadlines, cancellation signals,
	// and other request-scoped values across API boundaries and between processes.
	spacecrafts := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)

	for spacecraftID, spacecraft := range spacecrafts {
		res, err := client.Index("spacecrafts", esutil.NewJSONReader(spacecraft), client.Index.WithDocumentID(strconv.Itoa(spacecraftID)))
		if err == nil { // nit: err patter should be in the form `if err != nil`
			defer res.Body.Close() // bug: this defer is called ina for loop, it won't defer as expected
			fmt.Println(res)
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

func DeleteIndex(ctx context.Context, indexName string) error {
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)

	res, err := client.Indices.Exists([]string{indexName}) // nit: no need to perform an extra query to verify existence
	if err != nil {
		return fmt.Errorf("failed to get info for index %q with err: %v", indexName, err)
	}
	// nit: early return on res.StatusCode == http.StatusNotFound: less nesting and more idiomatic
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

// FIXME: understand how to wait for all of the data to be indexed
func IndexSpacecraftAsDocumentsAsync(ctx context.Context) error {
	spacecrafts := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      "spacecrafts",
		Client:     client,
		NumWorkers: 5,
		Refresh:    "wait_for",
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
