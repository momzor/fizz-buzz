package usecase

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"fizzbuzz/internal/domain"
	repomocks "fizzbuzz/internal/usecase/mocks"
)

func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestFizzBuzzUseCase_Execute_Success(t *testing.T) {
	repo := repomocks.NewStatsRepository(t)
	uc := NewFizzBuzz(repo, noopLogger())

	req := domain.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 5, Str1: "fizz", Str2: "buzz"}
	repo.EXPECT().IncrementHit(context.Background(), requestMatcher(req)).Return(nil).Once()
	result, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "fizz", "4", "buzz"}, result)
}

func TestFizzBuzzUseCase_Execute_InvalidRequest(t *testing.T) {
	repo := repomocks.NewStatsRepository(t)
	uc := NewFizzBuzz(repo, noopLogger())

	req := domain.FizzBuzzRequest{Int1: 0, Int2: 5, Limit: 5, Str1: "fizz", Str2: "buzz"}
	result, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "IncrementHit", mock.Anything, mock.Anything)
}

func TestFizzBuzzUseCase_Execute_StatsFailureDoesNotFailRequest(t *testing.T) {
	repo := repomocks.NewStatsRepository(t)
	uc := NewFizzBuzz(repo, noopLogger())

	req := domain.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 3, Str1: "fizz", Str2: "buzz"}
	repo.EXPECT().IncrementHit(context.Background(), requestMatcher(req)).Return(errors.New("boom")).Once()
	result, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "fizz"}, result)
}

func requestMatcher(expected domain.FizzBuzzRequest) interface{} {
	return mock.MatchedBy(func(actual domain.FizzBuzzRequest) bool {
		return actual.Int1 == expected.Int1 &&
			actual.Int2 == expected.Int2 &&
			actual.Limit == expected.Limit &&
			actual.Str1 == expected.Str1 &&
			actual.Str2 == expected.Str2
	})
}

func TestStatsUseCase_MostFrequent(t *testing.T) {
	entry := &domain.StatEntry{
		Request: domain.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"},
		Hits:    42,
	}
	repo := repomocks.NewStatsRepository(t)
	uc := NewStats(repo)
	repo.EXPECT().MostFrequent(context.Background()).Return(entry, nil).Once()

	got, err := uc.MostFrequent(context.Background())

	require.NoError(t, err)
	assert.Equal(t, entry, got)
}

func TestStatsUseCase_MostFrequent_Error(t *testing.T) {
	repo := repomocks.NewStatsRepository(t)
	uc := NewStats(repo)
	repo.EXPECT().MostFrequent(context.Background()).Return(nil, errors.New("db down")).Once()

	got, err := uc.MostFrequent(context.Background())

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestStatsUseCase_MostFrequent_NoData(t *testing.T) {
	repo := repomocks.NewStatsRepository(t)
	uc := NewStats(repo)
	repo.EXPECT().MostFrequent(context.Background()).Return(nil, nil).Once()

	got, err := uc.MostFrequent(context.Background())

	require.NoError(t, err)
	assert.Nil(t, got)
}
