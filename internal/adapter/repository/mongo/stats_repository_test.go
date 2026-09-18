package mongorepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	"fizzbuzz/internal/domain"
)

func TestStatsRepository_IncrementHit(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := NewStatsRepository(mt.DB)
		err := repo.IncrementHit(context.Background(), domain.FizzBuzzRequest{
			Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz",
		})

		require.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "boom"}})

		repo := NewStatsRepository(mt.DB)
		err := repo.IncrementHit(context.Background(), domain.FizzBuzzRequest{
			Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz",
		})

		require.Error(t, err)
	})
}

func TestStatsRepository_MostFrequent(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("found", func(mt *mtest.T) {
		expected := mtest.CreateCursorResponse(1, "fizzbuzz.fizzbuzz_stats", mtest.FirstBatch, bson.D{
			{Key: "int1", Value: 3},
			{Key: "int2", Value: 5},
			{Key: "limit", Value: 100},
			{Key: "str1", Value: "fizz"},
			{Key: "str2", Value: "buzz"},
			{Key: "hits", Value: int64(42)},
		})
		killCursors := mtest.CreateCursorResponse(0, "fizzbuzz.fizzbuzz_stats", mtest.NextBatch)
		mt.AddMockResponses(expected, killCursors)

		repo := NewStatsRepository(mt.DB)
		entry, err := repo.MostFrequent(context.Background())

		require.NoError(t, err)
		require.NotNil(t, entry)
		assert.Equal(t, int64(42), entry.Hits)
		assert.Equal(t, domain.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}, entry.Request)
	})

	mt.Run("no documents", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "fizzbuzz.fizzbuzz_stats", mtest.FirstBatch))

		repo := NewStatsRepository(mt.DB)
		entry, err := repo.MostFrequent(context.Background())

		require.NoError(t, err)
		assert.Nil(t, entry)
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "boom"}})

		repo := NewStatsRepository(mt.DB)
		entry, err := repo.MostFrequent(context.Background())

		require.Error(t, err)
		assert.Nil(t, entry)
	})
}

func TestStatsRepository_EnsureIndexes(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := NewStatsRepository(mt.DB)
		err := repo.EnsureIndexes(context.Background())

		require.NoError(t, err)
	})

	mt.Run("error", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "boom"}})

		repo := NewStatsRepository(mt.DB)
		err := repo.EnsureIndexes(context.Background())

		require.Error(t, err)
	})
}
