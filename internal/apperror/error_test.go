package apperror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_WrapsCause(t *testing.T) {
	cause := errors.New("database failed")
	err := InternalError("query statistics", cause)

	assert.Equal(t, "query statistics: database failed", err.Error())
	assert.True(t, errors.Is(err, cause))

	var typed *Error
	require.True(t, errors.As(err, &typed))
	assert.Equal(t, Internal, typed.Kind)
}

func TestError_WithoutOperation(t *testing.T) {
	cause := errors.New("invalid request")
	err := New(InvalidArgument, "", cause)

	assert.Equal(t, "invalid request", err.Error())
}
