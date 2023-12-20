package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"spacecraft/internal/domain"

	"golang.org/x/sync/errgroup"
)

func LoadspacecraftAsync(ctx context.Context, url string) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		fmt.Println("LoadspacecraftAsync() took", elapsedTime)
	}()
	res, err := http.Get(fmt.Sprintf("%v/spacecraft?pageNumber=0&pageSize=100", url))
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
	var spacecrafts []*domain.Spacecraft
	if len(spacecraftWrapper.Data) == 0 {
		return ctx, fmt.Errorf("no spacecraft in the page")
	}
	spacecrafts = append(spacecrafts, spacecraftWrapper.Data...)
	group, gctx := errgroup.WithContext(ctx)
	ch := make(chan []*domain.Spacecraft, spacecraftWrapper.TotalElements)
	for i := 1; i < spacecraftWrapper.TotalPages; i++ {
		routineUrl := fmt.Sprintf("%v/spacecraft?pageNumber=%d&pageSize=100", url, i)
		group.Go(func() error {
			if err := fetchSpacecraft(gctx, routineUrl, ch); err != nil {
				return err
			}
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return ctx, err
	}
	close(ch)
	for msg := range ch {
		spacecrafts = append(spacecrafts, msg...)
	}
	ctx = domain.WithSpacecrafts(ctx, spacecrafts)
	return ctx, nil
}

func fetchSpacecraft(ctx context.Context, url string, ch chan []*domain.Spacecraft) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		fmt.Println(url)
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode > 299 {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("err while parsing not positive answer: %w", err)
			}
			return fmt.Errorf("err fetching data: %v", string(body))
		}
		var spacecraftWrapper domain.SpacecraftWrapper
		if err := json.NewDecoder(res.Body).Decode(&spacecraftWrapper); err != nil {
			return err
		}
		ch <- spacecraftWrapper.Data
		return nil
	}
}
