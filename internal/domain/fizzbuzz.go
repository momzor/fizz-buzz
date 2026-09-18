// Package domain contains the core business entities and rules
package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fizzbuzz/internal/apperror"
)

// MaxLimit bounds the size of a generated sequence to protect the service
const MaxLimit = 100_000

// FizzBuzzRequest is the value object describing a fizzbuzz computation.
type FizzBuzzRequest struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string

	validated bool
}

// Validate ensures the request parameters are usable to generate a sequence.
func (r *FizzBuzzRequest) Validate() error {
	r.validated = false
	var errs []string

	if r.Int1 <= 0 {
		errs = append(errs, "int1 must be a strictly positive integer")
	}
	if r.Int2 <= 0 {
		errs = append(errs, "int2 must be a strictly positive integer")
	}
	if r.Limit <= 0 {
		errs = append(errs, "limit must be a strictly positive integer")
	}
	if r.Limit > MaxLimit {
		errs = append(errs, fmt.Sprintf("limit must not exceed %d", MaxLimit))
	}
	if strings.TrimSpace(r.Str1) == "" {
		errs = append(errs, "str1 must not be empty")
	}
	if strings.TrimSpace(r.Str2) == "" {
		errs = append(errs, "str2 must not be empty")
	}

	if len(errs) > 0 {
		return apperror.Invalid("validate fizzbuzz request", errors.New(strings.Join(errs, "; ")))
	}

	r.validated = true
	return nil
}

// Key returns a canonical, order-independent string representation of the
// request. It is used as the deduplication/statistics key.
func (r FizzBuzzRequest) Key() string {
	return strconv.Itoa(r.Int1) + ":" + strconv.Itoa(r.Int2) + ":" +
		strconv.Itoa(r.Limit) + ":" + r.Str1 + ":" + r.Str2
}

// Generate computes the fizzbuzz sequence described by the request.
func (r FizzBuzzRequest) Generate() []string {
	if !r.validated {
		panic("fizzbuzz: Generate called before Validate")
	}

	result := make([]string, r.Limit)
	for i := 1; i <= r.Limit; i++ {
		switch {
		case i%r.Int1 == 0 && i%r.Int2 == 0:
			result[i-1] = r.Str1 + r.Str2
		case i%r.Int1 == 0:
			result[i-1] = r.Str1
		case i%r.Int2 == 0:
			result[i-1] = r.Str2
		default:
			result[i-1] = strconv.Itoa(i)
		}
	}
	return result
}
