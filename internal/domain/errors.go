package domain

import (
	"fmt"
)

const (
	errMessageUnexpected             = "an unexpected error occurred, please verify govm.log for more information"
	errMessageVersionNotAvailable    = "go version \"%s\" is not available"
	errMessageNoUpdatesAvailable     = "no %s updates available for version \"%s\""
	errMessageNoGoInstallationsFound = "no go installations found"

	ErrCodeListVersions = 1

	ErrCodeCheckUserHome               = 1
	ErrCodeCheckVersion                = 2
	ErrCodeDownloadRemoveDir           = 3
	ErrCodeDownloadCreateFile          = 4
	ErrCodeDownloadVersion             = 5
	ErrCodeChecksumDownload            = 6
	ErrCodeChecksumOpenFile            = 7
	ErrCodeChecksumCopy                = 8
	ErrCodeChecksumMismatch            = 9
	ErrCodeRemoveVersion               = 10
	ErrCodeUntarCreateDir              = 11
	ErrCodeUntarExtract                = 12
	ErrCodeAddToPathNoShellsFound      = 13
	ErrCodeAddToPathStat               = 14
	ErrCodeAddToPathRead               = 15
	ErrCodeAddToPathWrite              = 16
	ErrCodeRemoveCurrentVersion        = 17
	ErrCodeRemoveFromPathNoShellsFound = 18
	ErrCodeRemoveFromPathStat          = 19
	ErrCodeRemoveFromPathRead          = 20
	ErrCodeRemoveFromPathWrite         = 21
)

type baseError struct {
	Message string
	Code    int
}

func (e *baseError) Error() string {
	return fmt.Sprintf("Error: %s Code: %d", e.Message, e.Code)
}

func NewUnexpectedError(code int) error {
	return &baseError{
		Message: errMessageUnexpected,
		Code:    code,
	}
}

func NewVersionNotAvailableError(version string) error {
	return &baseError{
		Message: fmt.Sprintf(errMessageVersionNotAvailable, version),
		Code:    1,
	}
}

func NewNoUpdatesAvailableError(strategy UpdateStrategy, version string) error {
	return &baseError{
		Message: fmt.Sprintf(errMessageNoUpdatesAvailable, string(strategy), version),
		Code:    1,
	}
}

func NewNoGoInstallationsFoundError() error {
	return &baseError{
		Message: errMessageNoGoInstallationsFound,
		Code:    1,
	}
}
