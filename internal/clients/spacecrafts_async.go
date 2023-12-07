package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"spacecraft/internal/domain"
)

var errChan = make(chan error)

// refactor it similarly to the sequantial version in client.go
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
	var spacecraft []*domain.Spacecraft
	if len(spacecraftWrapper.Data) == 0 {
		return ctx, fmt.Errorf("no spacecraft in the page")
	}
	spacecraft = append(spacecraft, spacecraftWrapper.Data...)
	var wg sync.WaitGroup
	ch := make(chan []*domain.Spacecraft, spacecraftWrapper.TotalElements)
	for i := 1; i < spacecraftWrapper.TotalPages; i++ {
		wg.Add(1)
		go func(url string) {
			fmt.Println(url)
			defer wg.Done()
			res, err := http.Get(url)
			if err != nil {
				errChan <- err
				return
			}
			defer res.Body.Close()
			var spacecraftWrapper domain.SpacecraftWrapper
			if err := json.NewDecoder(res.Body).Decode(&spacecraftWrapper); err != nil {
				errChan <- err
				return
			}
			ch <- spacecraftWrapper.Data
		}(fmt.Sprintf("%v/spacecraft?pageNumber=%d&pageSize=100", url, i))
	}
	// nit: possibly, to make it really parallel this should be written as
	// otherwise the async blocks until all the Get request are completed
	// go func() {
	//	wg.Wait()
	//	close(ch)
	//}
	wg.Wait()
	close(ch)
	select { // not sure if this select here is needed...err can be checked differently or using a errGroup
	// https://pkg.go.dev/golang.org/x/sync/errgroup
	case err = <-errChan:
		return ctx, err
	default:
		fmt.Println("no errors received!")
	}
	for msg := range ch {
		spacecraft = append(spacecraft, msg...)
	}
	ctx = context.WithValue(ctx, domain.ModelsKey, spacecraft)
	return ctx, nil
}
