// Package apperror defines typed application errors shared across layers.
package apperror

import "fmt"

// Kind identifies the public category of an application error.
type Kind uint8

const (
	InvalidArgument Kind = iota + 1
	NotFound
	Internal
	Unavailable
)

// Error enriched with operation and underlying cause.
type Error struct {
	Kind Kind
	Op   string
	Err  error
}

func (e *Error) Error() string {
	if e.Op == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

// Unwrap allows errors.Is and errors.As to inspect the original cause.
func (e *Error) Unwrap() error {
	return e.Err
}

// New creates a typed application error.
func New(kind Kind, op string, err error) *Error {
	return &Error{Kind: kind, Op: op, Err: err}
}

// Invalid creates an invalid-argument application error.
func Invalid(op string, err error) *Error {
	return New(InvalidArgument, op, err)
}

// InternalError creates an internal application error.
func InternalError(op string, err error) *Error {
	return New(Internal, op, err)
}

// UnavailableError creates a dependency-unavailable application error.
func UnavailableError(op string, err error) *Error {
	return New(Unavailable, op, err)
}
