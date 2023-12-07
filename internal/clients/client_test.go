package clients_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spacecraft/internal/clients"
	"spacecraft/internal/domain"
)

func TestClient_Load(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get Request: ", r.URL)
		pageN, e := strconv.Atoi(r.URL.Query().Get("pageNumber"))
		require.NoError(t, e)

		mockResp := domain.SpacecraftWrapper{
			PageNumber:       pageN,
			PageSize:         100,
			NumberOfElements: 1,
			TotalPages:       3,
			TotalElements:    3,
			Data: []*domain.Spacecraft{{
				Uid:    "1234",
				Name:   "ABC",
				Status: "OK",
			}},
		}
		b, err := json.Marshal(&mockResp)
		require.NoError(t, err)
		_, _ = w.Write(b)
	}))
	defer api.Close()

	c := clients.NewClient(api.URL, api.Client())

	sp, err := c.Load(context.Background())
	assert.NoError(t, err)
	assert.Len(t, sp, 3)
}
