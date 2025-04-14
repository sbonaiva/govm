package domain_test

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	action := domain.Action{
		Version:        "go1.19.13",
		HomeDir:        "/home/user",
		UpdateStrategy: domain.MinorStrategy,
	}

	assert.Equal(t, fmt.Sprintf("go1.19.13.%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH), action.Filename())
	assert.Equal(t, fmt.Sprintf("%s/go1.19.13.%s-%s.tar.gz", os.TempDir(), runtime.GOOS, runtime.GOARCH), action.DownloadFile())

	assert.Equal(t, "/home/user/.govm/go/bin", action.HomeGoBinDir())
	assert.Equal(t, "/home/user/.govm/go", action.HomeGoDir())
	assert.Equal(t, "/home/user/.govm", action.HomeGovmDir())

	assert.Equal(t, "# The next line is added by govm\nexport PATH=$PATH:/home/user/.govm/go/bin\nexport GOPATH=$HOME/go\n", action.Export())
	assert.Equal(t, domain.MinorStrategy, action.UpdateStrategy)
	assert.NoError(t, action.CheckUpdateStrategy())
}

func TestAction_UpdateStrategyError(t *testing.T) {
	action := domain.Action{
		Version:        "go1.19.13",
		HomeDir:        "/home/user",
		UpdateStrategy: "invalid_strategy",
	}

	assert.Error(t, action.CheckUpdateStrategy())
}
