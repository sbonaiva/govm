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

func TestNewNoUpdatesAvailableError(t *testing.T) {
	// Arrange
	version := "1.16.0"
	strategy := PatchStrategy

	// Act
	err := NewNoUpdatesAvailableError(strategy, version)

	// Assert
	baseErr, ok := err.(*baseError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf(errMessageNoUpdatesAvailable, string(strategy), version), baseErr.Message)
	assert.Equal(t, 1, baseErr.Code)
	assert.Equal(t, "Error: no patch updates available for version \"1.16.0\" Code: 1", err.Error())
}

func TestNewNoGoInstallationsFoundError(t *testing.T) {
	// Act
	err := NewNoGoInstallationsFoundError()

	// Assert
	baseErr, ok := err.(*baseError)
	assert.True(t, ok)
	assert.Equal(t, errMessageNoGoInstallationsFound, baseErr.Message)
	assert.Equal(t, 1, baseErr.Code)
	assert.Equal(t, "Error: no go installations found Code: 1", err.Error())
}

func TestNewInvalidUpdateStrategyError(t *testing.T) {
	// Arrange
	strategy := MajorStrategy

	// Act
	err := NewInvalidUpdateStrategyError(strategy)

	// Assert
	baseErr, ok := err.(*baseError)
	assert.True(t, ok)
	assert.Equal(t, fmt.Sprintf(errMessageInvalidUpdateStrategy, string(strategy)), baseErr.Message)
	assert.Equal(t, 1, baseErr.Code)
	assert.Equal(t, "Error: \"major\" is not a valid update strategy Code: 1", err.Error())
}
