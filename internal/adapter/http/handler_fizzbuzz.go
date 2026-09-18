package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"fizzbuzz/internal/apperror"
	"fizzbuzz/internal/domain"
)

// FizzBuzzExecutor is the contract the injected use case must satisfy.
type FizzBuzzExecutor interface {
	Execute(ctx context.Context, req domain.FizzBuzzRequest) ([]string, error)
}

// FizzBuzzHandler serves the fizzbuzz generation endpoint.
type FizzBuzzHandler struct {
	usecase FizzBuzzExecutor
}

// NewFizzBuzzHandler builds a FizzBuzzHandler.
func NewFizzBuzzHandler(usecase FizzBuzzExecutor) *FizzBuzzHandler {
	return &FizzBuzzHandler{usecase: usecase}
}

// ServeHTTP handles GET /api/v1/fizzbuzz.
//
// @Summary      Generate a customized fizzbuzz sequence
// @Description  Returns numbers from 1 to limit, replacing multiples of int1 by str1, multiples of int2 by str2, and multiples of both by str1+str2.
// @Tags         fizzbuzz
// @Produce      json
// @Param        int1   query     int     true  "First divisor"        example(3)
// @Param        int2   query     int     true  "Second divisor"       example(5)
// @Param        limit  query     int     true  "Upper bound (inclusive)" example(100)
// @Param        str1   query     string  true  "Replacement for multiples of int1" example(fizz)
// @Param        str2   query     string  true  "Replacement for multiples of int2" example(buzz)
// @Success      200 {object} fizzBuzzResponse
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /api/v1/fizzbuzz [get]
func (h *FizzBuzzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := parseFizzBuzzRequest(r)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	result, err := h.usecase.Execute(r.Context(), req)
	if err != nil {
		writeApplicationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, fizzBuzzResponse{Result: result})
}

func parseFizzBuzzRequest(r *http.Request) (domain.FizzBuzzRequest, error) {
	q := r.URL.Query()

	int1, err := parseIntParam(q, "int1")
	if err != nil {
		return domain.FizzBuzzRequest{}, err
	}
	int2, err := parseIntParam(q, "int2")
	if err != nil {
		return domain.FizzBuzzRequest{}, err
	}
	limit, err := parseIntParam(q, "limit")
	if err != nil {
		return domain.FizzBuzzRequest{}, err
	}

	str1, err := parseStringParam(q, "str1")
	if err != nil {
		return domain.FizzBuzzRequest{}, err
	}
	str2, err := parseStringParam(q, "str2")
	if err != nil {
		return domain.FizzBuzzRequest{}, err
	}

	return domain.FizzBuzzRequest{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}, nil
}

func parseIntParam(q map[string][]string, name string) (int, error) {
	values, ok := q[name]
	if !ok || len(values) == 0 || values[0] == "" {
		return 0, apperror.Invalid("parse fizzbuzz request", &paramError{name: name, reason: "is required"})
	}
	v, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, apperror.Invalid("parse fizzbuzz request", &paramError{name: name, reason: "must be an integer"})
	}
	return v, nil
}

func parseStringParam(q map[string][]string, name string) (string, error) {
	values, ok := q[name]
	if !ok || len(values) == 0 || values[0] == "" {
		return "", apperror.Invalid("parse fizzbuzz request", &paramError{name: name, reason: "is required"})
	}
	return values[0], nil
}

type paramError struct {
	name   string
	reason string
}

func (e *paramError) Error() string {
	return e.name + " " + e.reason
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeApplicationError(w http.ResponseWriter, err error) {
	status, message := applicationErrorResponse(err)
	writeError(w, status, message)
}

func applicationErrorResponse(err error) (int, string) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		return http.StatusInternalServerError, "internal server error"
	}

	switch appErr.Kind {
	case apperror.InvalidArgument:
		return http.StatusBadRequest, appErr.Err.Error()
	case apperror.NotFound:
		return http.StatusNotFound, "resource not found"
	case apperror.Unavailable:
		return http.StatusServiceUnavailable, "service temporarily unavailable"
	case apperror.Internal:
		fallthrough
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
