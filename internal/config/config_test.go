package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := Load()

	assert.Equal(t, "8080", cfg.HTTPPort)
	assert.Equal(t, "0.0.0.0:8080", cfg.HTTPAddr)
	assert.Equal(t, "http://localhost:8080", cfg.PublicBaseURL)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
	assert.Equal(t, "fizzbuzz", cfg.MongoDatabase)
	assert.Equal(t, "fizzbuzz-api", cfg.ServiceName)
	assert.True(t, cfg.TracingEnabled)
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("SHUTDOWN_TIMEOUT", "5s")
	t.Setenv("MONGO_URI", "mongodb://mongo:27017")
	t.Setenv("MONGO_DATABASE", "custom_db")
	t.Setenv("TRACING_ENABLED", "false")
	t.Setenv("ENVIRONMENT", "production")

	cfg := Load()

	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "0.0.0.0:8080", cfg.HTTPAddr)
	assert.Equal(t, "http://localhost:8080", cfg.PublicBaseURL)
	assert.Equal(t, 5*time.Second, cfg.ShutdownTimeout)
	assert.Equal(t, "mongodb://mongo:27017", cfg.MongoURI)
	assert.Equal(t, "custom_db", cfg.MongoDatabase)
	assert.False(t, cfg.TracingEnabled)
	assert.Equal(t, "production", cfg.Environment)
}

func TestLoad_InvalidDurationFallsBackToDefault(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "not-a-duration")

	cfg := Load()
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
}
