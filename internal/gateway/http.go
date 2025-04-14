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

type HttpGateway interface {
	GetVersions(ctx context.Context) (domain.VersionsResponse, error)
	GetChecksum(ctx context.Context, version string) (string, error)
	VersionExists(ctx context.Context, version string) (bool, error)
	DownloadVersion(ctx context.Context, action *domain.Action, file *os.File) error
}

type HttpConfig struct {
	GoVersionURL  string
	GoDownloadURL string
}

type httpClient struct {
	config *HttpConfig
	client *http.Client
}

func NewHttpGateway(config *HttpConfig) HttpGateway {
	return &httpClient{
		client: http.DefaultClient,
		config: config,
	}
}

func (r *httpClient) GetVersions(ctx context.Context) (domain.VersionsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.config.GoVersionURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Error while creating request", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return domain.VersionsResponse{}, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Error while making request", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return domain.VersionsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected status code", slog.String("GoDevClient", "GetVersions"), slog.String("status", resp.Status))
		return domain.VersionsResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versions []domain.VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		slog.ErrorContext(ctx, "Error decoding body", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return domain.VersionsResponse{}, err
	}

	compatibleVersions := make([]domain.VersionResponse, 0, len(versions))

	for _, v := range versions {
		if v.IsCompatible() && v.Stable {
			compatibleVersions = append(compatibleVersions, v)
		}
	}

	return domain.VersionsResponse{
		Versions: compatibleVersions,
	}, nil
}

func (r *httpClient) GetChecksum(ctx context.Context, version string) (string, error) {
	res, err := r.GetVersions(ctx)
	if err != nil {
		return "", err
	}

	for _, v := range res.Versions {
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
	res, err := r.GetVersions(ctx)
	if err != nil {
		return false, err
	}

	for _, v := range res.Versions {
		if v.Version == version && v.IsCompatible() && v.Stable {
			return true, nil
		}
	}

	return false, err
}

func (r *httpClient) DownloadVersion(ctx context.Context, action *domain.Action, file *os.File) error {
	resp, err := r.client.Get(fmt.Sprintf(r.config.GoDownloadURL, action.Filename()))
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
