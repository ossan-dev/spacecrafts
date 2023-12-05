package handlers_test

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"spacecraft/cmd/server/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/spacecraft.json
var data []byte

func TestGetspacecraft(t *testing.T) {
	err := json.Unmarshal(data, &handlers.Spacecraft)
	require.Nil(t, err)
	srv := http.NewServeMux()
	srv.HandleFunc("/spacecraft", handlers.GetSpacecraft)
	testSuite := []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "get paginated spacecraft",
			url:                "/spacecraft?pageSize=100&pageNumber=0",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "pageSize NAN",
			url:                "/spacecraft?pageSize=abc&pageNumber=0",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "pageNumber NAN",
			url:                "/spacecraft?pageSize=100&pageNumber=abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "pageSize not provided",
			url:                "/spacecraft?pageNumber=0",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "get last spacecraft page",
			url:                "/spacecraft?pageSize=100&pageNumber=14",
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
