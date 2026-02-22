package db

import "errors"

var (
	// Concurrency
	ErrConcurrentModification = errors.New("concurrent modification detected")
)
