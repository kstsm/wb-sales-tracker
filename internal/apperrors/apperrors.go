package apperrors

import "errors"

var (
	ErrItemNotFound = errors.New("item not found")
	ErrEmptyDate    = errors.New("empty date string")
)
