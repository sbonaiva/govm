package gateway

import (
	"context"
	"net/http"
	"os"

	"github.com/sbonaiva/govm/internal/domain"
)

type HttpGateway interface {
	GetVersions(ctx context.Context) ([]domain.GoVersionResponse, error)
	GetChecksum(ctx context.Context, version string) (string, error)
	VersionExists(ctx context.Context, version string) (bool, error)
	DownloadVersion(ctx context.Context, install domain.Install, file *os.File) error
}

func NewHttpGateway() HttpGateway {
	return &httpClient{
		client: &http.Client{},
	}
}

type OsGateway interface {
	GetUserHomeDir() (string, error)
}

func NewOsGateway() OsGateway {
	return &osClient{}
}
