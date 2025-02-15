package domain

import (
	"fmt"
)

const (
	errMessageUnexpected          = "an unexpected error occurred, please verify govm.log for more information"
	errMessageVersionNotAvailable = "go version \"%s\" is not available"

	ErrCodeInstallCheckUserHome          = 1
	ErrCodeInstallCheckVersion           = 2
	ErrCodeInstallDownloadRemoveDir      = 3
	ErrCodeInstallDownloadCreateFile     = 4
	ErrCodeInstallDownloadVersion        = 5
	ErrCodeInstallChecksumDownload       = 6
	ErrCodeInstallChecksumOpenFile       = 7
	ErrCodeInstallChecksumCopy           = 8
	ErrCodeInstallChecksumMismatch       = 9
	ErrCodeInstallRemovePreviousVersion  = 10
	ErrCodeInstallUntarCreateDir         = 11
	ErrCodeInstallUntarExtract           = 12
	ErrCodeInstallAddToPathNoShellsFound = 13
	ErrCodeInstallAddToPathStat          = 14
	ErrCodeInstallAddToPathRead          = 15
	ErrCodeInstallAddToPathWrite         = 16

	ErrCodeUninstallCheckUserHome               = 1
	ErrCodeUninstallCheckVersionStat            = 2
	ErrCodeUninstallCheckVersionNotDir          = 3
	ErrCodeUninstallRemoveDir                   = 4
	ErrCodeUninstallRemoveFromPathNoShellsFound = 5
	ErrCodeUninstallRemoveFromPathStat          = 6
	ErrCodeUninstallRemoveFromPathRead          = 7
	ErrCodeUninstallRemoveFromPathWrite         = 8
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
