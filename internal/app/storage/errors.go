package storage

import "errors"

// storage errors
var (
	ErrConflictURL = errors.New("this URL already exists")
)
