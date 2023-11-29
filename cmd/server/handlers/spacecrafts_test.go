package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"spacecrafts/cmd/server/handlers"
	"spacecrafts/cmd/server/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSpacecrafts(t *testing.T) {
	var err error
	handlers.Spacecrafts, err = store.LoadSpacecraftsFromFile("../store/spacecrafts.json")
	require.NoError(t, err)
	srv := http.NewServeMux()
	srv.HandleFunc("/spacecrafts", handlers.GetSpacecrafts)
	testSuite := []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "get paginated spacecrafts",
			url:                "/spacecrafts?pageSize=100&pageNumber=0",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "pageSize NAN",
			url:                "/spacecrafts?pageSize=abc&pageNumber=0",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "pageNumber NAN",
			url:                "/spacecrafts?pageSize=100&pageNumber=abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "pageSize not provided",
			url:                "/spacecrafts?pageNumber=0",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "get last spacecrafts page",
			url:                "/spacecrafts?pageSize=100&pageNumber=14",
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tt := range testSuite {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, tt.url, nil)
			require.NoError(t, err)
			// Act
			srv.ServeHTTP(w, r)
			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
