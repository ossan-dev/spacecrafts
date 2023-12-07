package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"spacecraft/internal/domain"
)

const getSpacecraftOp = "/spacecraft"

type Client struct {
	base   string
	client *http.Client
}

func NewClient(url string, client *http.Client) *Client {
	return &Client{
		base:   url,
		client: client,
	}
}

func (c Client) Load(ctx context.Context) ([]*domain.Spacecraft, error) {
	var (
		out  []*domain.Spacecraft
		page int
	)

	for {
		url := fmt.Sprintf("%s%s?pageNumber=%d&pageSize=100", c.base, getSpacecraftOp, page)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		var spacecraftWrapper domain.SpacecraftWrapper
		err = c.do(req, &spacecraftWrapper)
		if err != nil {
			return nil, fmt.Errorf("can't perform Load(): %w", err)
		}

		out = append(out, spacecraftWrapper.Data...)
		page++
		if page == spacecraftWrapper.TotalPages {
			break
		}
	}

	return out, nil
}

func (c Client) do(req *http.Request, out interface{}) error {
	if out == nil {
		return fmt.Errorf("out interface can't be nil")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("when performing request: %w ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errors.New("status code is not a 2xx")
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}
	return nil
}
