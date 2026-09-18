package domain

// StatEntry is the number of times a given FizzBuzzRequest has been requested.
type StatEntry struct {
	Request FizzBuzzRequest
	Hits    int64
}
