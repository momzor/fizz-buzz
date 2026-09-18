package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewRouter builds the complete HTTP router of the service: health checks,
// versioned API routes and Swagger UI documentation.
func NewRouter(fizzBuzzHandler *FizzBuzzHandler, statsHandler *StatsHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodOptions},
	}))

	r.Get("/healthz", healthHandler)

	r.Route("/api/v1", func(api chi.Router) {
		api.Handle("/fizzbuzz", otelhttp.NewHandler(fizzBuzzHandler, "fizzbuzz"))
		api.Handle("/statistics", otelhttp.NewHandler(statsHandler, "statistics"))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	return r
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
