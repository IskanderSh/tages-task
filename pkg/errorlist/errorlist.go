package errorlist

import "errors"

var (
	ErrNotFound           = errors.New("not found error")
	ErrNoFileWithSuchName = errors.New("no file with such name")
	ErrInvalidValues      = errors.New("error taking invalid values")
)
