package gateway

import (
	"context"
	"os"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/mock"
)

type HttpGatewayMock struct {
	mock.Mock
}

func (m *HttpGatewayMock) GetVersions(ctx context.Context) ([]domain.GoVersionResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.GoVersionResponse), args.Error(1)
}

func (m *HttpGatewayMock) GetChecksum(ctx context.Context, version string) (string, error) {
	args := m.Called(ctx, version)
	return args.String(0), args.Error(1)
}

func (m *HttpGatewayMock) VersionExists(ctx context.Context, version string) (bool, error) {
	args := m.Called(ctx, version)
	return args.Bool(0), args.Error(1)
}

func (m *HttpGatewayMock) DownloadVersion(ctx context.Context, install domain.Install, file *os.File) error {
	args := m.Called(ctx, install, file)
	return args.Error(0)
}

type OsGatewayMock struct {
	mock.Mock
}

func (m *OsGatewayMock) GetUserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *OsGatewayMock) CreateDir(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *OsGatewayMock) RemoveDir(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *OsGatewayMock) CreateFile(path string) (*os.File, error) {
	args := m.Called(path)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *OsGatewayMock) OpenFile(path string) (*os.File, error) {
	args := m.Called(path)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *OsGatewayMock) RemoveFile(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *OsGatewayMock) GetEnv(key string) string {
	args := m.Called(key)
	return args.String(0)
}
