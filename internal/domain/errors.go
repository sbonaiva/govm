package domain

import (
	"errors"
)

var (
	ErrUnexpected = errors.New("An unexpected error occurred, please verify govm.log for more information")
)
