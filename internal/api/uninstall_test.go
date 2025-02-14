package api_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestUninstallCmd_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockUninstallHandler)
	mockHandler.On("Handle", ctx, &domain.Uninstall{}).Return(nil)

	input := []byte("y")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	cmd := api.NewUninstallCmd(ctx, mockHandler)

	// Act
	output, _ := util.CaptureOutput(func() error {
		cmd.Run(cmd, []string{})
		return nil
	})

	// Assert
	assert.Contains(t, output, "Go uninstalled successfully!")
	mockHandler.AssertExpectations(t)
}

func TestUninstallCmd_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockUninstallHandler)
	mockHandler.On("Handle", ctx, &domain.Uninstall{}).Return(errors.New("failed to uninstall Go"))

	input := []byte("y")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	cmd := api.NewUninstallCmd(ctx, mockHandler)

	// Act
	output, _ := util.CaptureOutput(func() error {
		cmd.Run(cmd, []string{})
		return nil
	})

	// Assert
	assert.Contains(t, output, "Error: failed to uninstall Go")
	mockHandler.AssertExpectations(t)
}

func TestUninstallCmd_Abort(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockHandler := new(handler.MockUninstallHandler)

	input := []byte("n")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}
	w.Close()

	// Restore stdin right after the test.
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	cmd := api.NewUninstallCmd(ctx, mockHandler)

	// Act
	output, _ := util.CaptureOutput(func() error {
		cmd.Run(cmd, []string{})
		return nil
	})

	// Assert
	assert.Contains(t, output, "Uninstall aborted by user")
	mockHandler.AssertExpectations(t)
}
