package util

import (
	"fmt"
	"runtime"
)

var GoVersionManager string

func init() {
	GoVersionManager = fmt.Sprintf("%s %s/%s", "0.0.2", runtime.GOOS, runtime.GOARCH)
}
