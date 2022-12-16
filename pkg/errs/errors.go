package errs

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("measurement with given ID not found")
	ErrNotImplemented = errors.New("not implemented")
)
