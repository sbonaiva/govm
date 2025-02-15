package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	comment = "# The next line is added by govm"
)

type Install struct {
	Version string
	HomeDir string
}

func (r Install) Filename() string {
	return fmt.Sprintf("%s.%s-%s.tar.gz", r.Version, runtime.GOOS, runtime.GOARCH)
}

func (r Install) DownloadFile() string {
	return filepath.Join(os.TempDir(), r.Filename())
}

func (r Install) HomeGovmDir() string {
	return filepath.Join(r.HomeDir, ".govm")
}

func (r Install) HomeGoDir() string {
	return filepath.Join(r.HomeGovmDir(), "go")
}

func (r Install) HomeGoBinDir() string {
	return filepath.Join(r.HomeGoDir(), "bin")
}

func (r Install) Export() string {
	return fmt.Sprintf(
		"%s\nexport PATH=$PATH:%s\nexport GOPATH=$HOME/go\n",
		comment,
		r.HomeGoBinDir(),
	)
}
