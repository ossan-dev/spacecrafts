package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"spacecraft/internal/domain"
)

func Loadspacecraft(ctx context.Context, url string) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		fmt.Println("Loadspacecraft() took", elapsedTime)
	}()
	var spacecrafts []*domain.Spacecraft
	pageNumber := 0
	for {
		res, err := http.Get(fmt.Sprintf("%v/spacecraft?pageNumber=%d&pageSize=100", url, pageNumber))
		if err != nil {
			return ctx, err
		}
		defer res.Body.Close()
		var spacecraftWrapper domain.SpacecraftWrapper
		if err = json.NewDecoder(res.Body).Decode(&spacecraftWrapper); err != nil {
			return ctx, err
		}
		if spacecraftWrapper.TotalPages == 0 {
			return ctx, fmt.Errorf("empty resultset")
		}
		if len(spacecraftWrapper.Data) == 0 {
			return ctx, fmt.Errorf("no spacecraft in the page")
		}
		spacecrafts = append(spacecrafts, spacecraftWrapper.Data...)
		pageNumber++
		if pageNumber == spacecraftWrapper.TotalPages {
			break
		}
	}
	ctx = domain.WithSpacecrafts(ctx, spacecrafts)
	return ctx, nil
}
