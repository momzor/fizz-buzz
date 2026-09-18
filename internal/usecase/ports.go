// Package usecase implements the application services (business use cases)
// of the fizzbuzz API. It orchestrates the domain layer and depends only on
// the ports (interfaces) it declares here, never on concrete adapters.
package usecase

import (
	"context"

	"fizzbuzz/internal/domain"
)

// StatsRepository is the port used to persist and query request statistics.
// It is implemented by driven adapters (e.g. MongoDB).
type StatsRepository interface {
	// IncrementHit atomically increments the hit counter for the given
	// request, creating the entry if it does not exist yet.
	IncrementHit(ctx context.Context, req domain.FizzBuzzRequest) error

	// MostFrequent returns the request with the highest hit count.
	// It returns (nil, nil) when no statistics have been recorded yet.
	MostFrequent(ctx context.Context) (*domain.StatEntry, error)
}
