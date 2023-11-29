package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"spacecraft/cmd/server/handlers"
	"spacecraft/cmd/server/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetspacecraft(t *testing.T) {
	var err error
	handlers.Spacecraft, err = store.LoadspacecraftFromFile("../store/spacecraft.json")
	require.NoError(t, err)
	srv := http.NewServeMux()
	srv.HandleFunc("/spacecraft", handlers.Getspacecraft)
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
