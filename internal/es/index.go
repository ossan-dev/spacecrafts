package es

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

func IndexSpacecraftAsDocuments(ctx context.Context, esClient *elasticsearch.Client) error {
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
	spacecrafts, err := domain.GetSpacecraftsFromCtx(ctx)
	if err != nil {
		return err
	}
	for spacecraftID, spacecraft := range spacecrafts {
		res, err := esClient.Index("spacecrafts", esutil.NewJSONReader(spacecraft), esClient.Index.WithDocumentID(strconv.Itoa(spacecraftID)))
		if err != nil {
			return err
		}
		res.Body.Close()
	}
	return nil
}

func DeleteIndex(esClient *elasticsearch.Client, indexName string) error {
	res, err := esClient.Indices.Delete([]string{indexName})
	if err != nil {
		return fmt.Errorf("failed to delete index %q with err: %v", indexName, err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("err while reading response body: %v", err)
		}
		return fmt.Errorf("failed to delete index %q with message: %v", indexName, string(data))
	}
	return nil
}

func IndexSpacecraftAsDocumentsAsync(ctx context.Context, esClient *elasticsearch.Client) error {
	spacecrafts, err := domain.GetSpacecraftsFromCtx(ctx)
	if err != nil {
		return err
	}
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      "spacecrafts",
		Client:     esClient,
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
