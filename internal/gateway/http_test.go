package gateway_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/stretchr/testify/assert"
)

func TestGetVersionsSuccess(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []domain.VersionResponse{
			{
				Version: "1.17",
				Stable:  true,
				Files: []domain.FileResponse{
					{Kind: "archive", OS: "linux", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "linux", Arch: "arm64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "arm64", SHA256: "dummychecksum"},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	result, err := gatewayInstance.GetVersions(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.Versions, 1)
	assert.Equal(t, "1.17", result.Versions[0].Version)
	assert.Len(t, result.Versions[0].Files, 4)
}

func TestGetVersionsErrorCreatingRequest(t *testing.T) {
	// Arrange
	config := &gateway.HttpConfig{
		GoVersionURL: "://invalid-url",
	}
	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	_, err := gatewayInstance.GetVersions(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "parse \"://invalid-url\": missing protocol scheme", err.Error())
}

func TestGetVersionsErrorUnexpectedStatusCode(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	_, err := gatewayInstance.GetVersions(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 404", err.Error())
}

func TestGetVersionsErrorDecode(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := "invalid"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	_, err := gatewayInstance.GetVersions(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "json: cannot unmarshal string into Go value of type []domain.VersionResponse", err.Error())
}

func TestGetChecksumVersionFound(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []domain.VersionResponse{
			{
				Version: "1.17",
				Stable:  true,
				Files: []domain.FileResponse{
					{Kind: "archive", OS: "linux", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "linux", Arch: "arm64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "arm64", SHA256: "dummychecksum"},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	result, err := gatewayInstance.GetChecksum(context.Background(), "1.17")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "dummychecksum", result)
}

func TestGetChecksumVersionNotFound(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []domain.VersionResponse{
			{
				Version: "1.16",
				Stable:  true,
				Files: []domain.FileResponse{
					{Kind: "archive", OS: "linux", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "linux", Arch: "arm64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "arm64", SHA256: "dummychecksum"},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	result, err := gatewayInstance.GetChecksum(context.Background(), "1.17")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "version 1.17 not found", err.Error())
	assert.Equal(t, "", result)
}

func TestGetChecksumUnexpectedStatusCode(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	result, err := gatewayInstance.GetChecksum(context.Background(), "1.17")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 404", err.Error())
	assert.Equal(t, "", result)
}

func TestVersionExistsVersionFound(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []domain.VersionResponse{
			{
				Version: "1.17",
				Stable:  true,
				Files: []domain.FileResponse{
					{Kind: "archive", OS: "linux", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "linux", Arch: "arm64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "arm64", SHA256: "dummychecksum"},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	exists, err := gatewayInstance.VersionExists(context.Background(), "1.17")

	// Assert
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestVersionExistsVersionNotFound(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []domain.VersionResponse{
			{
				Version: "1.16",
				Stable:  true,
				Files: []domain.FileResponse{
					{Kind: "archive", OS: "linux", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "linux", Arch: "arm64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "amd64", SHA256: "dummychecksum"},
					{Kind: "archive", OS: "darwin", Arch: "arm64", SHA256: "dummychecksum"},
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(versions)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	exists, err := gatewayInstance.VersionExists(context.Background(), "1.17")

	// Assert
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestVersionExistsUnexpectedStatusCode(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL,
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	// Act
	exists, err := gatewayInstance.VersionExists(context.Background(), "1.17")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 500", err.Error())
	assert.False(t, exists)
}

func TestDownloadVersion_Success(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "file content")
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL + "/%s",
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	file, err := os.CreateTemp("", "downloaded_file")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Act
	err = gatewayInstance.DownloadVersion(context.Background(), &domain.Action{Version: "1.17"}, file)

	// Assert
	assert.NoError(t, err)
	fileContent, _ := os.ReadFile(file.Name())
	assert.Equal(t, "file content", string(fileContent))
}

func TestDownloadVersion_ErrorDownloading(t *testing.T) {
	// Arrange
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := &gateway.HttpConfig{
		GoVersionURL:  server.URL,
		GoDownloadURL: server.URL + "/%s",
	}

	gatewayInstance := gateway.NewHttpGateway(config)

	file, err := os.CreateTemp("", "downloaded_file")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	// Act
	err = gatewayInstance.DownloadVersion(context.Background(), &domain.Action{Version: "1.17"}, file)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "unexpected status code: 404", err.Error())
}
