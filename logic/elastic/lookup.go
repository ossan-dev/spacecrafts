package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

// FIXME: move below "internal"
// nit: QuerySpacecraftByDocumentID it' s tecnically a GetByID, as no ES Search is performed.
// see: https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-get.html
func QuerySpacecraftByDocumentID(ctx context.Context, index, documentID string) (*domain.Spacecraft, error) {
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)
	res, err := client.Get(index, documentID)
	if err != nil {
		// major: should be %w to wrap an error in fmt.Errorf instead of %v? (same in other places)
		return nil, fmt.Errorf("err while looking up the document with ID: %v with err: %v", documentID, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			// nit: here we lose the not 200 error context
			return nil, fmt.Errorf("err while reading response body: %v", err)
		}
		return nil, fmt.Errorf("err with the Elasticsearch request processing: %v", string(data))
	}
	var lookupRes domain.LookupResponse
	// a note on correct usage of json.NewDecoder:
	// https://mottaquikarim.github.io/dev/posts/you-might-not-be-using-json.decoder-correctly-in-golang/
	if err = json.NewDecoder(res.Body).Decode(&lookupRes); err != nil {
		return nil, fmt.Errorf("failed to decode the elastic search result: %v", err)
	}
	return lookupRes.Source, nil
}
