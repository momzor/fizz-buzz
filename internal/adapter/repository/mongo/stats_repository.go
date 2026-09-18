// Package mongorepo implements the usecase.StatsRepository port on top of
// MongoDB. It is a driven adapter: it depends on the usecase package, never
// the other way around.
package mongorepo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fizzbuzz/internal/domain"
)

const collectionName = "fizzbuzz_stats"

// statDocument is the MongoDB representation of a domain.StatEntry.
// The five request fields form the document's natural (unique) key.
type statDocument struct {
	Int1  int    `bson:"int1"`
	Int2  int    `bson:"int2"`
	Limit int    `bson:"limit"`
	Str1  string `bson:"str1"`
	Str2  string `bson:"str2"`
	Hits  int64  `bson:"hits"`
}

// StatsRepository persists fizzbuzz request statistics in MongoDB.
type StatsRepository struct {
	collection *mongo.Collection
}

// NewStatsRepository builds a StatsRepository bound to the given database.
func NewStatsRepository(db *mongo.Database) *StatsRepository {
	return &StatsRepository{collection: db.Collection(collectionName)}
}

// EnsureIndexes creates the unique index required to make hit increments
// atomic and idempotent-safe. It should be called once at startup.
func (r *StatsRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "int1", Value: 1},
			{Key: "int2", Value: 1},
			{Key: "limit", Value: 1},
			{Key: "str1", Value: 1},
			{Key: "str2", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("uniq_request"),
	})
	if err != nil {
		return fmt.Errorf("creating unique index: %w", err)
	}
	return nil
}

// IncrementHit atomically increments the hit counter for the given request.
func (r *StatsRepository) IncrementHit(ctx context.Context, req domain.FizzBuzzRequest) error {
	filter := bson.M{
		"int1":  req.Int1,
		"int2":  req.Int2,
		"limit": req.Limit,
		"str1":  req.Str1,
		"str2":  req.Str2,
	}
	update := bson.M{"$inc": bson.M{"hits": 1}}
	opts := options.Update().SetUpsert(true)

	if _, err := r.collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("incrementing hit counter: %w", err)
	}
	return nil
}

// MostFrequent returns the request with the highest hit count, or nil if no
// statistics have been recorded yet.
func (r *StatsRepository) MostFrequent(ctx context.Context) (*domain.StatEntry, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "hits", Value: -1}})

	var doc statDocument
	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying most frequent request: %w", err)
	}

	return &domain.StatEntry{
		Request: domain.FizzBuzzRequest{
			Int1:  doc.Int1,
			Int2:  doc.Int2,
			Limit: doc.Limit,
			Str1:  doc.Str1,
			Str2:  doc.Str2,
		},
		Hits: doc.Hits,
	}, nil
}
