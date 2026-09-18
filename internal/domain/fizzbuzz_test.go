package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFizzBuzzRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     FizzBuzzRequest
		wantErr string
	}{
		{
			name: "valid request",
			req:  FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"},
		},
		{
			name:    "int1 zero",
			req:     FizzBuzzRequest{Int1: 0, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"},
			wantErr: "int1 must be a strictly positive integer",
		},
		{
			name:    "int1 negative",
			req:     FizzBuzzRequest{Int1: -3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"},
			wantErr: "int1 must be a strictly positive integer",
		},
		{
			name:    "int2 zero",
			req:     FizzBuzzRequest{Int1: 3, Int2: 0, Limit: 100, Str1: "fizz", Str2: "buzz"},
			wantErr: "int2 must be a strictly positive integer",
		},
		{
			name:    "limit zero",
			req:     FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 0, Str1: "fizz", Str2: "buzz"},
			wantErr: "limit must be a strictly positive integer",
		},
		{
			name:    "limit too large",
			req:     FizzBuzzRequest{Int1: 3, Int2: 5, Limit: MaxLimit + 1, Str1: "fizz", Str2: "buzz"},
			wantErr: "limit must not exceed",
		},
		{
			name:    "str1 empty",
			req:     FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "  ", Str2: "buzz"},
			wantErr: "str1 must not be empty",
		},
		{
			name:    "str2 empty",
			req:     FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: ""},
			wantErr: "str2 must not be empty",
		},
		{
			name:    "all invalid combines errors",
			req:     FizzBuzzRequest{Int1: 0, Int2: 0, Limit: 0, Str1: "", Str2: ""},
			wantErr: "int1 must be a strictly positive integer; int2 must be a strictly positive integer; limit must be a strictly positive integer; str1 must not be empty; str2 must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.True(t, strings.Contains(err.Error(), tt.wantErr), "got error: %s", err.Error())
		})
	}
}

func TestFizzBuzzRequest_Generate(t *testing.T) {
	req := FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"}
	require.NoError(t, req.Validate())
	got := req.Generate()
	want := []string{
		"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz",
		"11", "fizz", "13", "14", "fizzbuzz",
	}
	assert.Equal(t, want, got)
}

func TestFizzBuzzRequest_Generate_CustomStrings(t *testing.T) {
	req := FizzBuzzRequest{Int1: 2, Int2: 7, Limit: 14, Str1: "foo", Str2: "bar"}
	require.NoError(t, req.Validate())
	got := req.Generate()
	want := []string{
		"1", "foo", "3", "foo", "5", "foo", "bar",
		"foo", "9", "foo", "11", "foo", "13", "foobar",
	}
	assert.Equal(t, want, got)
}

func TestFizzBuzzRequest_Generate_RequiresValidation(t *testing.T) {
	req := FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"}

	assert.PanicsWithValue(t, "fizzbuzz: Generate called before Validate", func() {
		req.Generate()
	})
}

func TestFizzBuzzRequest_Key(t *testing.T) {
	req := FizzBuzzRequest{Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz"}
	assert.Equal(t, "3:5:100:fizz:buzz", req.Key())
}
