package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/sbonaiva/govm/internal/domain"
)

const (
	goVersionsURL = "https://go.dev/dl/?mode=json&include=all"
	goDownloadURL = "https://go.dev/dl/%s"
)

type GoDevClient interface {
	GetVersions(ctx context.Context) ([]domain.GoVersionResponse, error)
	VersionExists(ctx context.Context, version string) (bool, error)
	DownloadVersion(ctx context.Context, install domain.Install, file *os.File) error
}

type goDevClient struct {
	httpClient *http.Client
}

func NewGoDevClient() GoDevClient {
	return &goDevClient{
		httpClient: &http.Client{},
	}
}

func (r *goDevClient) GetVersions(ctx context.Context) ([]domain.GoVersionResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, goVersionsURL, nil)
	if err != nil {
		slog.Error("Error while creating request", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return []domain.GoVersionResponse{}, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		slog.Error("Error while making request", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
		return []domain.GoVersionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Unexpected status code", slog.String("GoDevClient", "GetVersions"), slog.String("status", resp.Status))
		return []domain.GoVersionResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versions []domain.GoVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		slog.Error("Error decoding body", slog.String("GoDevClient", "GetVersions"), slog.String("error", err.Error()))
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

func (r *goDevClient) VersionExists(ctx context.Context, version string) (bool, error) {
	versions, err := r.GetVersions(ctx)
	if err != nil {
		return false, err
	}

	for _, v := range versions {
		if v.Version == version {
			return true, nil
		}
	}

	return false, nil
}

func (r *goDevClient) DownloadVersion(ctx context.Context, install domain.Install, file *os.File) error {
	resp, err := r.httpClient.Get(fmt.Sprintf(goDownloadURL, install.Filename()))
	if err != nil {
		slog.Error("Error while downloading file", slog.String("GoDevClient", "DownloadVersion"), slog.String("error", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Unexpected status code", slog.String("GoDevClient", "DownloadVersion"), slog.String("status", resp.Status))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		slog.Error("Error while copying file", slog.String("GoDevClient", "DownloadVersion"), slog.String("error", err.Error()))
		return err
	}

	return nil
}
