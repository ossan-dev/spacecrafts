package es

import (
	"fmt"
	"net/http"
	"time"

	"spacecraft/internal/es"

	"github.com/stretchr/testify/require"
)

func (s *EsSuite) TestConnect() {
	s.Run("right url", func() {
		// act
		esClient, err := es.Connect(fmt.Sprintf("http://%v", s.Endpoint))

		// assert
		require.Nil(s.T(), err)
		require.NotNil(s.T(), esClient)
		res, err := esClient.Ping()
		require.Nil(s.T(), err)
		require.NotNil(s.T(), res)
		require.Equal(s.T(), http.StatusOK, res.StatusCode)
	})

	s.Run("wrong url", func() {
		// act
		esClient, err := es.Connect(fmt.Sprintf("http://%v", fmt.Sprintf("%v%d", s.Endpoint, time.Now().UnixNano())))

		// assert
		require.Nil(s.T(), err)
		require.NotNil(s.T(), esClient)
		res, err := esClient.Ping()
		require.Nil(s.T(), res)
		require.NotNil(s.T(), err)
	})
}
