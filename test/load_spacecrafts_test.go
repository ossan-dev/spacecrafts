package test

import (
	"context"
	"fmt"
	"testing"

	"spacecraft/domain"
	"spacecraft/logic/webclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLoadspacecraft(t *testing.T) {
	// arrange
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "../",
		},
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForExec([]string{"/webserver"}),
	}
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	require.Nil(t, err)
	endpoint, err := container.Endpoint(ctx, "")
	require.Nil(t, err)

	// act
	ctx, err = webclient.Loadspacecraft(ctx, fmt.Sprintf("http://%v", endpoint))

	// assert
	require.Nil(t, err)
	spacecraft, ok := ctx.Value(domain.ModelsKey).([]*domain.Spacecraft)
	assert.True(t, ok)
	assert.NotEmpty(t, spacecraft)

	// teardown
	if err := container.Terminate(ctx); err != nil {
		t.Fatal(err)
	}
}
