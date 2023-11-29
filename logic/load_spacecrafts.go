package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"spacecraft/domain"
)

func Fetchspacecraft(url string) ([]*domain.Spacecraft, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var stapiRes domain.StapiResponse
	if err = json.NewDecoder(res.Body).Decode(&stapiRes); err != nil {
		return nil, err
	}
	if stapiRes.Spacecraft == nil || len(stapiRes.Spacecraft) == 0 {
		return nil, fmt.Errorf("no spacecraft in the page")
	}
	var spacecraft []*domain.Spacecraft
	spacecraft = append(spacecraft, stapiRes.Spacecraft...)
	return spacecraft, nil
}

func Loadspacecraft(ctx context.Context) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		fmt.Println("Loadspacecraft() took", elapsedTime)
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
	var spacecraft []*domain.Spacecraft
	if stapiRes.Spacecraft == nil || len(stapiRes.Spacecraft) == 0 {
		return ctx, fmt.Errorf("no spacecraft in the page")
	}
	spacecraft = append(spacecraft, stapiRes.Spacecraft...)
	for i := 1; i < stapiRes.Page.TotalPages; i++ {
		newSpacrafts, err := Fetchspacecraft(fmt.Sprintf("https://stapi.co/api/v1/rest/spacecraft/search?pageNumber=%d&pageSize=100", i))
		if err != nil {
			return ctx, err
		}
		spacecraft = append(spacecraft, newSpacrafts...)
	}
	ctx = context.WithValue(ctx, domain.ModelsKey, spacecraft)
	return ctx, nil
}
