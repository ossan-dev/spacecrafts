package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ITSuite struct {
	suite.Suite
	Endpoint  string
	Ctx       context.Context
	Container testcontainers.Container
}

func (s *ITSuite) SetupSuite() {
	var err error
	s.Ctx = context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "../",
		},
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	s.Container, err = testcontainers.GenericContainer(
		s.Ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	require.Nil(s.T(), err)
	endpoint, err := s.Container.Endpoint(s.Ctx, "")
	require.Nil(s.T(), err)
	s.Endpoint = endpoint
}

func (s *ITSuite) TearDownSuite() {
	if s.Container != nil {
		if err := s.Container.Terminate(s.Ctx); err != nil {
			require.Nil(s.T(), err)
		}
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ITSuite))
}
