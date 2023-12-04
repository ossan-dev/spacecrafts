package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"spacecraft/cmd/server/handlers"
	"spacecraft/cmd/server/store"
	"spacecraft/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetspacecraft(t *testing.T) {
	var err error
	// [x]: fix it
	// bug: the  "../store/spacecraft.json" path is wrong.
	// also best to create a testdata folder under handlers to store the  spacecraft.json `fixture file`
	// see for example: https://dave.cheney.net/2016/05/10/test-fixtures-in-go
	handlers.Spacecraft, err = store.LoadSpacecraftFromFile("../../../spacecraft.json")
	require.NoError(t, err)

	srv := http.NewServeMux()
	srv.HandleFunc("/spacecraft", handlers.GetSpacecraft)
	testSuite := []struct {
		name               string
		url                string
		expectedStatusCode int
		expectedResponse   domain.SpacecraftWrapper // major: it's important to test the expectedResponse
	}{
		{
			name:               "get paginated spacecraft",
			url:                "/spacecraft?pageSize=100&pageNumber=0",
			expectedStatusCode: http.StatusOK,
			expectedResponse: domain.SpacecraftWrapper{
				PageNumber:       0,
				PageSize:         100,
				NumberOfElements: 100,
				TotalPages:       15,
				TotalElements:    1443,
			},
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
		// major: this test was missing
		{
			name:               "get a specific page",
			url:                "/spacecraft?pageSize=2&pageNumber=3",
			expectedStatusCode: http.StatusOK,
			expectedResponse: domain.SpacecraftWrapper{
				PageNumber:       3,
				PageSize:         2,
				NumberOfElements: 2,
				TotalPages:       722,
				TotalElements:    1443,
			},
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

			var sp domain.SpacecraftWrapper
			if reflect.DeepEqual(tt.expectedResponse, sp) {
				return
			} // assert expectedResponse, if not zero

			err = json.Unmarshal(w.Body.Bytes(), &sp)
			require.NoError(t, err)
			// set Data to nil  from `sp` response: this simplifies the test...
			// ...ideally in prod code at least len(Data) should be tested
			sp.Data = nil
			assert.EqualValues(t, tt.expectedResponse, sp)
		})
	}
}
