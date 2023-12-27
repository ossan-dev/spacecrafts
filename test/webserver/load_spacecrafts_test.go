package webserver

import (
	"context"
	"fmt"

	"spacecraft/internal/clients"
	"spacecraft/internal/domain"

	"github.com/stretchr/testify/require"
)

func (s *WebServerSuite) TestLoadspacecraft() {
	// act
	ctx, err := clients.Loadspacecraft(context.Background(), fmt.Sprintf("http://%v", s.Endpoint))

	// assert
	require.Nil(s.T(), err)
	spacecrafs, err := domain.GetSpacecraftsFromCtx(ctx)
	require.Nil(s.T(), err)
	s.NotEmpty(spacecrafs)
}
