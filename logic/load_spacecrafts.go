package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"spacecrafts/domain"
)

func FetchSpacecrafts(url string) ([]*domain.Spacecraft, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var stapiRes domain.StapiResponse
	if err = json.NewDecoder(res.Body).Decode(&stapiRes); err != nil {
		return nil, err
	}
	if stapiRes.Spacecrafts == nil || len(stapiRes.Spacecrafts) == 0 {
		return nil, fmt.Errorf("no spacecrafts in the page")
	}
	var spacecrafts []*domain.Spacecraft
	spacecrafts = append(spacecrafts, stapiRes.Spacecrafts...)
	return spacecrafts, nil
}

func LoadSpacecrafts(ctx context.Context) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		fmt.Println("LoadSpacecrafts() took", elapsedTime)
	}()
	res, err := http.Get("https://stapi.co/api/v1/rest/spacecraft/search?pageNumber=0&pageSize=100")
	if err != nil {
		return ctx, err
	}
	defer res.Body.Close()
	var stapiRes domain.StapiResponse
	if err = json.NewDecoder(res.Body).Decode(&stapiRes); err != nil {
		return ctx, err
	}
	if stapiRes.Page == nil || stapiRes.Page.TotalPages == 0 {
		return ctx, fmt.Errorf("empty resultset")
	}
	var spacecrafts []*domain.Spacecraft
	if stapiRes.Spacecrafts == nil || len(stapiRes.Spacecrafts) == 0 {
		return ctx, fmt.Errorf("no spacecrafts in the page")
	}
	spacecrafts = append(spacecrafts, stapiRes.Spacecrafts...)
	for i := 1; i < stapiRes.Page.TotalPages; i++ {
		newSpacrafts, err := FetchSpacecrafts(fmt.Sprintf("https://stapi.co/api/v1/rest/spacecraft/search?pageNumber=%d&pageSize=100", i))
		if err != nil {
			return ctx, err
		}
		spacecrafts = append(spacecrafts, newSpacrafts...)
	}
	ctx = context.WithValue(ctx, domain.ModelsKey, spacecrafts)
	return ctx, nil
}
