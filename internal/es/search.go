package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"spacecraft/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
)

func SearchByStatusAndUidPrefix(esClient *elasticsearch.Client, index, uidPrefix, status string) ([]*domain.Spacecraft, int, error) {
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
	if err := json.NewEncoder(&searchBuffer).Encode(search); err != nil {
		return nil, 0, fmt.Errorf("err while encoding the search req: %w", err)
	}
	response, err := esClient.Search(
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(&searchBuffer),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithSize(30),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("err while invoking elasticsearch: %w", err)
	}
	defer response.Body.Close()
	var searchRes domain.SearchResponse
	jsonDec := json.NewDecoder(response.Body)
	for {
		if err := jsonDec.Decode(&searchRes); err != nil {
			if err == io.EOF {
				break
			}
			return nil, 0, fmt.Errorf("err while reading streaming response: %w", err)
		}
	}
	var res []*domain.Spacecraft
	for _, v := range searchRes.Hits.Hits {
		res = append(res, v.Source)
	}
	return res, searchRes.Hits.Total.Value, nil
}
