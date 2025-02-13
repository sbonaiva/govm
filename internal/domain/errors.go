package domain

import (
	"errors"
	"fmt"
)

const (
	errUnexpected          = "an unexpected error occurred, please verify govm.log for more information"
	errVersionNotAvailable = "go version \"%s\" is not available"
)

var (
	ErrUnexpected          = errors.New(errUnexpected)
	ErrVersionNotAvailable = func(version string) error {
		return fmt.Errorf(errVersionNotAvailable, version)
	}
)
