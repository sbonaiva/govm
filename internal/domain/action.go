package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type UpdateStrategy string

const (
	exportBegin  = "# The next lines are added by govm"
	exportGoRoot = "export GOROOT=%s"
	exportGoPath = "export GOPATH=$HOME/go"
	exportPath   = "export PATH=$PATH:%s"
	exportEnd    = "# End of govm path"

	MajorStrategy UpdateStrategy = "major"
	MinorStrategy UpdateStrategy = "minor"
	PatchStrategy UpdateStrategy = "patch"
)

type Action struct {
	Version          string
	HomeDir          string
	InstalledVersion string
	UpdateStrategy   UpdateStrategy
}

func (r Action) Filename() string {
	return fmt.Sprintf("%s.%s-%s.tar.gz", r.Version, runtime.GOOS, runtime.GOARCH)
}

func (r Action) DownloadFile() string {
	return filepath.Join(os.TempDir(), r.Filename())
}

func (r Action) HomeGovmDir() string {
	return filepath.Join(r.HomeDir, ".govm")
}

func (r Action) HomeGoDir() string {
	return filepath.Join(r.HomeGovmDir(), "go")
}

func (r Action) HomeGoBinDir() string {
	return filepath.Join(r.HomeGoDir(), "bin")
}

func (r Action) Export() string {
	return strings.Join([]string{
		exportBegin,
		fmt.Sprintf(exportGoRoot, r.HomeGoDir()),
		exportGoPath,
		fmt.Sprintf(exportPath, r.HomeGoBinDir()),
		exportEnd,
	}, "\n")
}

func (r *Action) CheckUpdateStrategy() error {
	switch r.UpdateStrategy {
	case MajorStrategy, MinorStrategy, PatchStrategy:
		return nil
	default:
		return NewInvalidUpdateStrategyError(r.UpdateStrategy)
	}
}
