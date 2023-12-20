package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func SearchByStatusAndUidPrefix(esClient *elasticsearch.Client, index, uidPrefix, status string) (res []*domain.Spacecraft, count int, err error) {
	var searchBuffer bytes.Buffer
	search := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"match_phrase_prefix": map[string]string{
						"uid": uidPrefix,
					},
				},
				"filter": map[string]interface{}{
					"term": map[string]string{
						"status": strings.ToLower(status),
					},
				},
			},
		},
	}
	if err = json.NewEncoder(&searchBuffer).Encode(search); err != nil {
		return nil, 0, fmt.Errorf("err while encoding the search req: %v", err)
	}
	response, err := esClient.Search(
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(&searchBuffer),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithSize(30), // this should come from pageSize
		esClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("err while invoking elasticsearch: %v", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("err while reading response body: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected elasticsearch err: %v", err)
	}
	var searchRes domain.SearchResponse
	if err = json.Unmarshal(data, &searchRes); err != nil {
		return nil, 0, fmt.Errorf("err while unmarshaling data: %v", err)
	}
	count = searchRes.Hits.Total.Value
	if searchRes.Hits.Total.Value > 0 { // nit: early return if searchRes.Hits.Total.Value == 0
		for _, v := range searchRes.Hits.Hits {
			res = append(res, v.Source)
		}
	}
	return
}
