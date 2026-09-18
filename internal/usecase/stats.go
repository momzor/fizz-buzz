package usecase

import (
	"context"

	"fizzbuzz/internal/apperror"
	"fizzbuzz/internal/domain"
)

// Stats exposes read access to fizzbuzz usage statistics.
type Stats struct {
	stats StatsRepository
}

// NewStats wires the statistics use case to its dependencies.
func NewStats(stats StatsRepository) *Stats {
	return &Stats{stats: stats}
}

// MostFrequent returns the most requested fizzbuzz parameters along with
// their hit count, or nil if no request has been recorded yet.
func (uc *Stats) MostFrequent(ctx context.Context) (*domain.StatEntry, error) {
	entry, err := uc.stats.MostFrequent(ctx)
	if err != nil {
		return nil, apperror.InternalError("get most frequent fizzbuzz request", err)
	}
	return entry, nil
}
