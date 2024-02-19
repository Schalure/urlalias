package aliasmaker

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLWasDeleted = errors.New("url was deleted")
)