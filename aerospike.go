package aerospike

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

const (
	aerospikeServicePort  = "3000/tcp"
	defaultAerospikeImage = "aerospike:ce-6.4.0.0"
)

type AerospikeContainer struct {
	testcontainers.Container
}

// RunContainer creates an instance of the Aerospike container type.
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*AerospikeContainer, error) {
	containerRequest := testcontainers.ContainerRequest{
		Image:        defaultAerospikeImage,
		ExposedPorts: []string{"3000/tcp"},
		WaitingFor:   newAerospikeWaitStrategy(),
	}
	genericContainerRequest := testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	}
	for _, opt := range opts {
		_ = opt.Customize(&genericContainerRequest)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to start Aerospike: %w", err)
	}

	return &AerospikeContainer{Container: container}, nil
}

// Port returns the port on which the Aerospike container is listening.
func (c AerospikeContainer) ServicePort(ctx context.Context) (int, error) {
	port, err := c.Container.MappedPort(ctx, aerospikeServicePort)
	if err != nil {
		return 0, err
	}
	return port.Int(), nil
}

// WithImage sets the image for the Aerospike container.
func WithImage(image string) testcontainers.CustomizeRequestOption {
	return testcontainers.WithImage(image)
}

// WithNamespace sets the default namespace that is created when Aerospike
// starts. By default, this is set to "test".
func WithNamespace(namespace string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		if req.Env == nil {
			req.Env = make(map[string]string)
		}
		req.Env["NAMESPACE"] = namespace

		return nil
	}
}
