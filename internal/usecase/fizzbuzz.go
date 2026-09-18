package usecase

import (
	"context"
	"log/slog"

	"fizzbuzz/internal/domain"
)

// FizzBuzz generates fizzbuzz sequences and records usage statistics.
type FizzBuzz struct {
	stats  StatsRepository
	logger *slog.Logger
}

// NewFizzBuzz wires the fizzbuzz use case to its dependencies.
func NewFizzBuzz(stats StatsRepository, logger *slog.Logger) *FizzBuzz {
	return &FizzBuzz{stats: stats, logger: logger}
}

// Execute validates the request, generates the sequence and records the
// request for statistics purposes. Statistics recording failures are logged
// but never fail the use case: they must not degrade the main feature.
func (uc *FizzBuzz) Execute(ctx context.Context, req domain.FizzBuzzRequest) ([]string, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	result := req.Generate()

	if err := uc.stats.IncrementHit(ctx, req); err != nil {
		uc.logger.ErrorContext(ctx, "failed to record fizzbuzz statistics", "error", err)
	}

	return result, nil
}
