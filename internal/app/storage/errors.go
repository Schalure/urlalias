package storage

import "errors"

var (
	ErrConflictURL = errors.New("this URL already exists")
)
