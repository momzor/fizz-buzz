package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	httpmocks "fizzbuzz/internal/adapter/http/mocks"
	"fizzbuzz/internal/apperror"
	"fizzbuzz/internal/domain"
)

func TestFizzBuzzHandler_ServeHTTP_Success(t *testing.T) {
	executor := httpmocks.NewFizzBuzzExecutor(t)
	reqParams := domain.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 3, Str1: "fizz", Str2: "buzz"}
	executor.EXPECT().Execute(mock.Anything, reqParams).Return([]string{"1", "2", "fizz"}, nil).Once()
	handler := NewFizzBuzzHandler(executor)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?int1=3&int2=5&limit=3&str1=fizz&str2=buzz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body fizzBuzzResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, []string{"1", "2", "fizz"}, body.Result)
}

func TestFizzBuzzHandler_ServeHTTP_MissingParam(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"missing int1", "int2=5&limit=3&str1=fizz&str2=buzz"},
		{"missing int2", "int1=3&limit=3&str1=fizz&str2=buzz"},
		{"missing limit", "int1=3&int2=5&str1=fizz&str2=buzz"},
		{"missing str1", "int1=3&int2=5&limit=3&str2=buzz"},
		{"missing str2", "int1=3&int2=5&limit=3&str1=fizz"},
		{"non integer int1", "int1=abc&int2=5&limit=3&str1=fizz&str2=buzz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := httpmocks.NewFizzBuzzExecutor(t)
			handler := NewFizzBuzzHandler(executor)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?"+tt.query, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, http.StatusBadRequest, rec.Code)

			var body errorResponse
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
			assert.NotEmpty(t, body.Error)
		})
	}
}

func TestFizzBuzzHandler_ServeHTTP_UseCaseValidationError(t *testing.T) {
	executor := httpmocks.NewFizzBuzzExecutor(t)
	executor.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, apperror.Invalid("validate fizzbuzz request", errors.New("invalid request"))).Once()
	handler := NewFizzBuzzHandler(executor)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/fizzbuzz?int1=3&int2=5&limit=3&str1=fizz&str2=buzz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestApplicationErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		statusCode int
		message    string
	}{
		{
			name:       "invalid argument",
			err:        apperror.Invalid("validate request", errors.New("bad input")),
			statusCode: http.StatusBadRequest,
			message:    "bad input",
		},
		{
			name:       "not found",
			err:        apperror.New(apperror.NotFound, "find request", errors.New("missing")),
			statusCode: http.StatusNotFound,
			message:    "resource not found",
		},
		{
			name:       "unavailable",
			err:        apperror.UnavailableError("query statistics", errors.New("database unavailable")),
			statusCode: http.StatusServiceUnavailable,
			message:    "service temporarily unavailable",
		},
		{
			name:       "internal",
			err:        apperror.InternalError("query statistics", errors.New("database failed")),
			statusCode: http.StatusInternalServerError,
			message:    "internal server error",
		},
		{
			name:       "unknown error",
			err:        errors.New("unexpected failure"),
			statusCode: http.StatusInternalServerError,
			message:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, message := applicationErrorResponse(tt.err)
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, tt.message, message)
		})
	}
}

func TestStatsHandler_ServeHTTP_Success(t *testing.T) {
	entry := &domain.StatEntry{
		Request: domain.FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"},
		Hits:    7,
	}
	query := httpmocks.NewStatsQuery(t)
	query.EXPECT().MostFrequent(mock.Anything).Return(entry, nil).Once()
	handler := NewStatsHandler(query)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body statisticsResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, statisticsResponse{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz", Hits: 7}, body)
}

func TestStatsHandler_ServeHTTP_NoData(t *testing.T) {
	query := httpmocks.NewStatsQuery(t)
	query.EXPECT().MostFrequent(mock.Anything).Return(nil, nil).Once()
	handler := NewStatsHandler(query)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body statisticsResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.Equal(t, statisticsResponse{}, body)
}

func TestStatsHandler_ServeHTTP_Error(t *testing.T) {
	query := httpmocks.NewStatsQuery(t)
	query.EXPECT().MostFrequent(mock.Anything).Return(nil, errors.New("db down")).Once()
	handler := NewStatsHandler(query)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/statistics", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestNewRouter_HealthAndRoutes(t *testing.T) {
	executor := httpmocks.NewFizzBuzzExecutor(t)
	executor.EXPECT().Execute(mock.Anything, mock.Anything).Return([]string{"1"}, nil).Once()
	query := httpmocks.NewStatsQuery(t)
	query.EXPECT().MostFrequent(mock.Anything).Return(nil, nil).Once()
	router := NewRouter(NewFizzBuzzHandler(executor), NewStatsHandler(query))

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp2, err := http.Get(server.URL + "/api/v1/fizzbuzz?int1=3&int2=5&limit=1&str1=fizz&str2=buzz")
	require.NoError(t, err)
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	resp3, err := http.Get(server.URL + "/api/v1/statistics")
	require.NoError(t, err)
	defer resp3.Body.Close()
	assert.Equal(t, http.StatusOK, resp3.StatusCode)
}
