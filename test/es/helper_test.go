package es

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"spacecraft/internal/domain"
	"spacecraft/internal/es"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type EsSuite struct {
	suite.Suite
	Endpoint  string
	Ctx       context.Context
	Container testcontainers.Container
	EsClient  *elasticsearch.Client
}

//go:embed testdata/spacecraft.json
var data []byte

func (s *EsSuite) SetupSuite() {
	var err error
	var spacecrafts []*domain.Spacecraft
	err = json.Unmarshal(data, &spacecrafts)
	require.Nil(s.T(), err)
	s.Ctx = domain.WithSpacecrafts(context.Background(), spacecrafts)
	req := testcontainers.ContainerRequest{
		Image:        "docker.elastic.co/elasticsearch/elasticsearch:7.14.0",
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"bootstrap.memory_lock": "true",
			"ES_JAVA_OPTS":          "-Xms1g -Xmx1g",
			"discovery.type":        "single-node",
			"node.name":             "lonely-gopher",
			"cluster.name":          "es4gophers",
		},
		WaitingFor: wait.ForListeningPort("9200/tcp"),
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
	s.EsClient, err = es.Connect(fmt.Sprintf("http://%v", s.Endpoint))
	require.NotNil(s.T(), s.EsClient)
	require.Nil(s.T(), err)
}

func (s *EsSuite) TearDownSuite() {
	if s.Container != nil {
		if err := s.Container.Terminate(s.Ctx); err != nil {
			require.Nil(s.T(), err)
		}
	}
}

func TestEsSuite(t *testing.T) {
	suite.Run(t, new(EsSuite))
}
