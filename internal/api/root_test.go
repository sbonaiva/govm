package api_test

import (
	"context"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	// Arrange
	ctx := context.Background()

	cmd := api.NewRootCmd(ctx, new(gateway.MockHttpGateway))

	// Act
	err := cmd.Execute()

	// Assert
	assert.NoError(t, err)
}
