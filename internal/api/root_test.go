package api_test

import (
	"context"
	"fmt"
	"runtime"
	"strings"
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
	actual, err := test.CaptureOutput(func() error {
		return cmd.Execute()
	})

	expected := strings.Join([]string{
		"::: Go Version Manager :::\n\n",
		"Usage:\n",
		"  govm [command]\n\n",
		"Available Commands:\n",
		"  completion  Generate the autocompletion script for the specified shell\n",
		"  help        Help about any command\n",
		"  install     Install a Go version\n",
		"  list        List all Go versions\n",
		"  uninstall   Uninstall a Go version\n\n",
		"Flags:\n",
		"  -h, --help      help for govm\n",
		"  -v, --version   version for govm\n\n",
		"Use \"govm [command] --help\" for more information about a command.\n",
	}, "")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, fmt.Sprintf("0.0.2 %s/%s", runtime.GOOS, runtime.GOARCH), cmd.Version)
}
