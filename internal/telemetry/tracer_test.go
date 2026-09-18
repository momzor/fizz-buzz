package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetup_Disabled(t *testing.T) {
	shutdown, err := Setup(context.Background(), "test-service", "1.0.0", "localhost:4318", false)
	require.NoError(t, err)
	require.NoError(t, shutdown(context.Background()))
}

func TestSetup_Enabled(t *testing.T) {
	shutdown, err := Setup(context.Background(), "test-service", "1.0.0", "localhost:4318", true)
	require.NoError(t, err)
	require.NoError(t, shutdown(context.Background()))
}
