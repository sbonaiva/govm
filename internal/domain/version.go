package domain

import (
	"fmt"
	"runtime"
)

type FileResponse struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Kind     string `json:"kind"`
	SHA256   string `json:"sha256"`
}

type VersionResponse struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []FileResponse
}

func (v VersionResponse) IsCompatible() bool {
	for _, f := range v.Files {
		if f.Kind == "archive" && f.OS == runtime.GOOS && f.Arch == runtime.GOARCH {
			return true
		}
	}
	return false
}

func (v VersionResponse) String(currentVersion string) string {

	if v.Version == currentVersion {
		return fmt.Sprintf("* %s", v.Version)
	}

	return v.Version
}

type VersionsResponse struct {
	Versions []VersionResponse
}

func (v VersionsResponse) StringSlice() []string {
	versions := make([]string, len(v.Versions))
	for i, version := range v.Versions {
		versions[i] = version.Version
	}
	return versions
}
