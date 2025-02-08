package domain

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var currentVersion string

type GoFileResponse struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Kind     string `json:"kind"`
	SHA256   string `json:"sha256"`
}

type GoVersionResponse struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []GoFileResponse
}

func init() {
	outputBytes, err := exec.Command("go", "version").Output()
	if err != nil {
		currentVersion = ""
		return
	}

	outputParts := strings.Split(string(outputBytes), " ")

	if len(outputParts) < 3 {
		currentVersion = ""
		return
	}

	currentVersion = outputParts[2]
}

func (v GoVersionResponse) IsCompatible() bool {
	for _, f := range v.Files {
		if f.Kind == "archive" && f.OS == runtime.GOOS && f.Arch == runtime.GOARCH {
			return true
		}
	}
	return false
}

func (v GoVersionResponse) String() string {

	if v.Version == currentVersion {
		return fmt.Sprintf("* %s", v.Version)
	}

	return v.Version
}
