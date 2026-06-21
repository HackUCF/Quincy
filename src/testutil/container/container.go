package container

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

// RunContainer starts a container, skipping t with a warning on any failure —
// including a broken or unavailable Docker daemon (which testcontainers panics on).
// Registers container termination as a t.Cleanup.
func RunContainer(ctx context.Context, t *testing.T, req testcontainers.GenericContainerRequest) testcontainers.Container {
	t.Helper()

	c, err := tryStartContainer(ctx, req)
	if err != nil {
		t.Logf("WARNING: container unavailable, skipping %s: %v", t.Name(), err)
		t.Skipf("container failed to start: %v", err)
	}

	t.Cleanup(func() { c.Terminate(ctx) })
	return c
}

// tryStartContainer wraps GenericContainer in a recover so panics from a broken
// Docker daemon are returned as errors rather than crashing the test binary.
func tryStartContainer(ctx context.Context, req testcontainers.GenericContainerRequest) (c testcontainers.Container, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return testcontainers.GenericContainer(ctx, req)
}
