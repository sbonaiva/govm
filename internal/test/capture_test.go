package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOutput(t *testing.T) {
	// Arrange
	expected := "Hello, World!\n"

	// Act
	output, err := CaptureOutput(func() error {
		fmt.Println("Hello, World!")
		return nil
	})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expected, output)
}

func TestCaptureOutputError(t *testing.T) {
	// Act
	output, err := CaptureOutput(func() error {
		return fmt.Errorf("error")
	})

	// Assert
	assert.Error(t, err)
	assert.Empty(t, output)
}
