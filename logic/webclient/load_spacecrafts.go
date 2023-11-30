package webclient

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
	var spacecraftWrapper domain.SpacecraftWrapper
	if err = json.NewDecoder(res.Body).Decode(&spacecraftWrapper); err != nil {
		return nil, err
	}
	if len(spacecraftWrapper.Data) == 0 {
		return nil, fmt.Errorf("no spacecraft in the page")
	}
	var spacecraft []*domain.Spacecraft
	spacecraft = append(spacecraft, spacecraftWrapper.Data...)
	return spacecraft, nil
}

func Loadspacecraft(ctx context.Context) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		fmt.Println("Loadspacecraft() took", elapsedTime)
	}()
	res, err := http.Get("http://localhost:8080/spacecraft?pageNumber=0&pageSize=100")
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
	var spacecraft []*domain.Spacecraft
	if len(spacecraftWrapper.Data) == 0 {
		return ctx, fmt.Errorf("no spacecraft in the page")
	}
	spacecraft = append(spacecraft, spacecraftWrapper.Data...)
	for i := 1; i < spacecraftWrapper.TotalPages; i++ {
		newSpacrafts, err := Fetchspacecraft(fmt.Sprintf("http://localhost:8080/spacecraft?pageNumber=%d&pageSize=100", i))
		if err != nil {
			return ctx, err
		}
		spacecraft = append(spacecraft, newSpacrafts...)
	}
	ctx = context.WithValue(ctx, domain.ModelsKey, spacecraft)
	return ctx, nil
}
