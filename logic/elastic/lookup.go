package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"spacecraft/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func QuerySpacecraftByDocumentID(ctx context.Context, index, documentID string) (*domain.Spacecraft, error) {
	client := ctx.Value(domain.ClientKey).(*elasticsearch.Client)
	res, err := client.Get(index, documentID)
	if err != nil {
		return nil, fmt.Errorf("err while looking up the document with ID: %v with err: %v", documentID, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("err while reading response body: %v", err)
		}
		return nil, fmt.Errorf("err with the Elasticsearch request processing: %v", string(data))
	}
	var lookupRes domain.LookupResponse
	if err = json.NewDecoder(res.Body).Decode(&lookupRes); err != nil {
		return nil, fmt.Errorf("failed to decode the elastic search result: %v", err)
	}
	return lookupRes.Source, nil
}
