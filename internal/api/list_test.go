package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/handler"
)

func TestNewListCmd_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockListHandler)
	mockHandler.On("Handle", ctx).Return(nil)

	cmd := api.NewListCmd(ctx, mockHandler)

	// Act
	cmd.Run(cmd, []string{})

	// Assert
	mockHandler.AssertExpectations(t)
}

func TestNewListCmd_ErrorHandling(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockListHandler)
	mockHandler.On("Handle", ctx).Return(errors.New("error"))

	cmd := api.NewListCmd(ctx, mockHandler)

	// Act
	cmd.Run(cmd, []string{})

	// Assert
	mockHandler.AssertExpectations(t)
}
