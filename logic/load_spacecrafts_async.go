package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"spacecraft/domain"
)

var errChan chan error = make(chan error)

func FetchspacecraftAsync(url string, wg *sync.WaitGroup, ch chan *domain.Spacecraft) {
	fmt.Println(url)
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		errChan <- err
		return
	}
	defer res.Body.Close()
	var stapiRes domain.StapiResponse
	if err := json.NewDecoder(res.Body).Decode(&stapiRes); err != nil {
		errChan <- err
		return
	}
	for _, v := range stapiRes.Spacecraft {
		select {
		case err := <-errChan:
			close(errChan)
			fmt.Printf("err while fetching: %v", err)
			return
		default:
			ch <- v
		}
	}
}

func LoadspacecraftAsync(ctx context.Context) (context.Context, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		fmt.Println("LoadspacecraftAsync() took", elapsedTime)
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
	var wg sync.WaitGroup
	ch := make(chan *domain.Spacecraft, stapiRes.Page.TotalElements)
	for i := 1; i < stapiRes.Page.TotalPages; i++ {
		wg.Add(1)
		go FetchspacecraftAsync(fmt.Sprintf("https://stapi.co/api/v1/rest/spacecraft/search?pageNumber=%d&pageSize=100", i), &wg, ch)
	}
	wg.Wait()
	close(ch)
	select {
	case err = <-errChan:
		return ctx, err
	default:
		fmt.Println("no errors received!")
	}
	for msg := range ch {
		spacecraft = append(spacecraft, msg)
	}
	ctx = context.WithValue(ctx, domain.ModelsKey, spacecraft)
	return ctx, nil
}
