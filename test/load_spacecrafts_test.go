package test

import (
	"context"
	"fmt"

	"spacecraft/internal/clients"
	"spacecraft/internal/domain"

	"github.com/stretchr/testify/require"
)

func (s *ITSuite) TestLoadspacecraft() {
	// act
	ctx, err := clients.Loadspacecraft(context.Background(), fmt.Sprintf("http://%v", s.Endpoint))

	// assert
	require.Nil(s.T(), err)
	spacecraft, ok := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	s.True(ok)
	s.NotEmpty(spacecraft)
}
