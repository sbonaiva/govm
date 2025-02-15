package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"github.com/sbonaiva/govm/internal/domain"
)

const (
	goVersionsURL = "https://go.dev/dl/?mode=json&include=all"
	goDownloadURL = "https://go.dev/dl/%s"
)

type HttpGateway interface {
	GetVersions(ctx context.Context) ([]domain.GoVersionResponse, error)
	GetChecksum(ctx context.Context, version string) (string, error)
	VersionExists(ctx context.Context, version string) (bool, error)
	DownloadVersion(ctx context.Context, install domain.Install, file *os.File) error
}

type httpClient struct {
	client *http.Client
}

func NewHttpGateway() HttpGateway {
	return &httpClient{
		client: &http.Client{},
	}
}

func (r *httpClient) GetVersions(ctx context.Context) ([]domain.GoVersionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, goVersionsURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Error while creating request", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return []domain.GoVersionResponse{}, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Error while making request", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return []domain.GoVersionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected status code", slog.String("GoDevClient", "GetVersions"), slog.String("status", resp.Status))
		return []domain.GoVersionResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versions []domain.GoVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		slog.ErrorContext(ctx, "Error decoding body", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return []domain.GoVersionResponse{}, err
	}

	compatibleVersions := make([]domain.GoVersionResponse, 0, len(versions))

	for _, v := range versions {
		if v.IsCompatible() && v.Stable {
			compatibleVersions = append(compatibleVersions, v)
		}
	}

	return compatibleVersions, nil
}

func (r *httpClient) GetChecksum(ctx context.Context, version string) (string, error) {
	versions, err := r.GetVersions(ctx)
	if err != nil {
		return "", err
	}

	for _, v := range versions {
		if v.Version == version {
			for _, f := range v.Files {
				if f.Kind == "archive" && f.OS == runtime.GOOS && f.Arch == runtime.GOARCH {
					return f.SHA256, nil
				}
			}
		}
	}

	return "", fmt.Errorf("version %s not found", version)
}

func (r *httpClient) VersionExists(ctx context.Context, version string) (bool, error) {
	versions, err := r.GetVersions(ctx)
	if err != nil {
		return false, err
	}

	for _, v := range versions {
		if v.Version == version && v.IsCompatible() && v.Stable {
			return true, nil
		}
	}

	return false, err
}

func (r *httpClient) DownloadVersion(ctx context.Context, install domain.Install, file *os.File) error {
	resp, err := r.client.Get(fmt.Sprintf(goDownloadURL, install.Filename()))
	if err != nil {
		slog.ErrorContext(ctx, "Error while downloading file", slog.String("GoDevClient", "DownloadVersion"), slog.String("error", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected status code", slog.String("GoDevClient", "DownloadVersion"), slog.String("status", resp.Status))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		slog.ErrorContext(ctx, "Error while copying file", slog.String("GoDevClient", "DownloadVersion"), slog.String("error", err.Error()))
		return err
	}

	return nil
}
