package domain_test

import (
	"fmt"
	"os"
	"path"
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

	filename := fmt.Sprintf("go1.19.13.%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH)

	assert.Equal(t, filename, action.Filename())
	assert.Equal(t, path.Join(os.TempDir(), filename), action.DownloadFile())

	assert.Equal(t, "/home/user/.govm/go/bin", action.HomeGoBinDir())
	assert.Equal(t, "/home/user/.govm/go", action.HomeGoDir())
	assert.Equal(t, "/home/user/.govm", action.HomeGovmDir())

	assert.Equal(t, "# The next lines are added by govm\nexport GOROOT=/home/user/.govm/go\nexport GOPATH=$HOME/go\nexport PATH=$PATH:/home/user/.govm/go/bin\n# End of govm path", action.Export())
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
