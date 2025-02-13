package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewInstallCmd_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockInstallHandler)
	mockHandler.On("Handle", ctx, &domain.Install{Version: "1.15.0"}).Return(nil)
	cmd := api.NewInstallCmd(ctx, mockHandler)
	cmd.SetArgs([]string{"1.15.0"})

	// Act
	err := cmd.Execute()

	// Assert
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestNewInstallCmd_ErrorHandling(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockInstallHandler)
	mockHandler.On("Handle", ctx, &domain.Install{Version: "1.15.0"}).Return(errors.New("install error"))
	cmd := api.NewInstallCmd(ctx, mockHandler)
	cmd.SetArgs([]string{"1.15.0"})

	// Act
	err := cmd.Execute()

	// Assert
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestNewInstallCmd_InvalidArgument(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockInstallHandler)
	cmd := api.NewInstallCmd(ctx, mockHandler)

	// Act
	err := cmd.Execute()

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "accepts 1 arg(s), received 0")
}
