package webclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"spacecraft/domain"
)

// Client is a quick refactoring of the LoadSpacecraft sequential function.
type Client struct {
	base   string // eg http://localhost:7000/
	client *http.Client
}

const getSpacecraftOp = "/spacecraft"

// Load fetch all pages from the server using getSpacecraftOp.
// todo: this might be done with ES Limit and Offset to just select a portion of the elements in the index
func (c Client) Load(ctx context.Context) ([]*domain.Spacecraft, error) {
	var (
		out  []*domain.Spacecraft
		page int
	)

	for {
		url := fmt.Sprintf("%s%s?pageNumber=%d&pageSize=100", c.base, getSpacecraftOp, page)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

		var spacecraftWrapper domain.SpacecraftWrapper
		err := c.do(req, &spacecraftWrapper)
		if err != nil {
			return nil, fmt.Errorf("can't perform Load(): %w", err)
		}

		out = append(out, spacecraftWrapper.Data...)
		page++
		if page == spacecraftWrapper.TotalPages {
			break
		}
	}

	return out, nil // return all docs in the index not very efficient...but ok for the sake of this demo!
}

// do is a generic method that takes a request and, if status code is 2xx, unmarshal the response into out.
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
