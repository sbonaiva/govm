package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnexpectedError(t *testing.T) {
	// Arrange
	code := ErrCodeListVersions

	// Act
	err := NewUnexpectedError(code)

	// Assert
	baseErr, ok := err.(*baseError)
	assert.True(t, ok)
	assert.Equal(t, errMessageUnexpected, baseErr.Message)
	assert.Equal(t, code, baseErr.Code)
	assert.Equal(t, "Error: an unexpected error occurred, please verify govm.log for more information Code: 1", err.Error())
}

func TestNewVersionNotAvailableError(t *testing.T) {
	// Arrange
	version := "1.16.0"

	// Act
	err := NewVersionNotAvailableError(version)

	// Assert
	baseErr, ok := err.(*baseError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf(errMessageVersionNotAvailable, version), baseErr.Message)
	assert.Equal(t, 1, baseErr.Code)
	assert.Equal(t, "Error: go version \"1.16.0\" is not available Code: 1", err.Error())
}
