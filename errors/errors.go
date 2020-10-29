package errors

import (
	"errors"
)

var (
	ErrBadRequest  = errors.New("invalid request")
	ErrBadResponse = errors.New("bad response error")
)
