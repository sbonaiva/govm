package domain_test

import (
	"runtime"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestVersionsResponse(t *testing.T) {

	versions := domain.VersionsResponse{
		Versions: []domain.VersionResponse{
			{
				Version: "1.20.5",
				Stable:  true,
				Files: []domain.FileResponse{
					{
						Filename: "go1.20.5.tar.gz",
						OS:       runtime.GOOS,
						Arch:     runtime.GOARCH,
						Kind:     "archive",
					},
				},
			},
			{
				Version: "1.20.6",
				Stable:  true,
				Files: []domain.FileResponse{
					{
						Filename: "go1.20.6.tar.gz",
						OS:       "solaris",
						Arch:     runtime.GOARCH,
						Kind:     "archive",
					},
				},
			},
		},
	}

	assert.True(t, versions.Versions[0].IsCompatible())
	assert.Equal(t, "* 1.20.5", versions.Versions[0].String("1.20.5"))
	assert.False(t, versions.Versions[1].IsCompatible())
	assert.Equal(t, "1.20.6", versions.Versions[1].String("1.20.5"))
	assert.Contains(t, versions.StringSlice(), "1.20.5")
	assert.Contains(t, versions.StringSlice(), "1.20.6")
}
