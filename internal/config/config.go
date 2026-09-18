// Package config centralizes environment-based configuration loading so the
// rest of the application never reads environment variables directly.
package config

import (
	"os"
	"time"
)

// Config holds all the runtime configuration of the service.
type Config struct {
	// HTTP
	HTTPAddr        string
	HTTPPort        string
	PublicBaseURL   string
	ShutdownTimeout time.Duration

	// MongoDB
	MongoURI      string
	MongoDatabase string

	// OpenTelemetry
	ServiceName    string
	ServiceVersion string
	OTLPEndpoint   string
	TracingEnabled bool

	// Misc
	Environment string
}

// Load reads configuration from environment variables, applying sane
// defaults for local development.
func Load() Config {
	return Config{
		HTTPAddr:        getEnv("HTTP_ADDR", "0.0.0.0:8080"),
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		PublicBaseURL:   getEnv("PUBLIC_BASE_URL", "http://localhost:8080"),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:   getEnv("MONGO_DATABASE", "fizzbuzz"),
		ServiceName:     getEnv("OTEL_SERVICE_NAME", "fizzbuzz-api"),
		ServiceVersion:  getEnv("SERVICE_VERSION", "dev"),
		OTLPEndpoint:    getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
		TracingEnabled:  getEnv("TRACING_ENABLED", "true") == "true",
		Environment:     getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
