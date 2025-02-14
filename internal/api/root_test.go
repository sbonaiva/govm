package api_test

import (
	"context"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	// Arrange
	ctx := context.Background()

	cmd := api.NewRootCmd(ctx, new(gateway.HttpGatewayMock), new(gateway.OsGatewayMock))

	// Act
	output, err := test.CaptureOutput(func() error {
		return cmd.Execute()
	})

	// Assert
	assert.Equal(t, "::: Go Version Manager :::\n\nUsage:\n  govm [command]\n\nAvailable Commands:\n  completion  Generate the autocompletion script for the specified shell\n  help        Help about any command\n  install     Install a Go version\n  list        List all Go versions\n  uninstall   Uninstall a Go version\n  use         Use a Go version\n\nFlags:\n  -h, --help      help for govm\n  -v, --version   version for govm\n\nUse \"govm [command] --help\" for more information about a command.\n", output)
	assert.NoError(t, err)
}
