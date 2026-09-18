package httpadapter

import (
	"context"
	"net/http"

	"fizzbuzz/internal/domain"
)

// StatsQuery is the subset of the stats use case needed by the HTTP handler.
type StatsQuery interface {
	MostFrequent(ctx context.Context) (*domain.StatEntry, error)
}

// StatsHandler serves the statistics endpoint.
type StatsHandler struct {
	usecase StatsQuery
}

// NewStatsHandler builds a StatsHandler.
func NewStatsHandler(usecase StatsQuery) *StatsHandler {
	return &StatsHandler{usecase: usecase}
}

// ServeHTTP handles GET /api/v1/statistics.
//
// @Summary      Get the most frequently requested fizzbuzz parameters
// @Description  Returns the request parameters that were used the most, along with their hit count.
// @Tags         statistics
// @Produce      json
// @Success      200 {object} statisticsResponse
// @Failure      500 {object} errorResponse
// @Router       /api/v1/statistics [get]
func (h *StatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	entry, err := h.usecase.MostFrequent(r.Context())
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toStatisticsResponse(entry))
}
