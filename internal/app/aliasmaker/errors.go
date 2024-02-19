package aliasmaker

import "errors"

var (
	ErrInternal = errors.New("something went wrong")

	ErrURLNotFound = errors.New("url not found")
	ErrURLWasDeleted = errors.New("url was deleted")

	ErrConflictURL = errors.New("this URL already exists")
)