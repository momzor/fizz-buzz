// Package httpadapter is the  HTTP adapter of the application
// package never knows about it.
package httpadapter

import "fizzbuzz/internal/domain"

// fizzBuzzResponse is the JSON body returned by GET /api/v1/fizzbuzz.
type fizzBuzzResponse struct {
	Result []string `json:"result"`
}

// statisticsResponse is the JSON body returned by GET /api/v1/statistics.
type statisticsResponse struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
	Hits  int64  `json:"hits"`
}

func toStatisticsResponse(entry *domain.StatEntry) statisticsResponse {
	if entry == nil {
		return statisticsResponse{}
	}
	return statisticsResponse{
		Int1:  entry.Request.Int1,
		Int2:  entry.Request.Int2,
		Limit: entry.Request.Limit,
		Str1:  entry.Request.Str1,
		Str2:  entry.Request.Str2,
		Hits:  entry.Hits,
	}
}

// errorResponse is the JSON body returned for any 4xx/5xx response.
type errorResponse struct {
	Error string `json:"error"`
}
