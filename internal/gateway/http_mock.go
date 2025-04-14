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

func (m *HttpGatewayMock) GetVersions(ctx context.Context) (domain.VersionsResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.VersionsResponse), args.Error(1)
}

func (m *HttpGatewayMock) GetVersionsList(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *HttpGatewayMock) GetChecksum(ctx context.Context, version string) (string, error) {
	args := m.Called(ctx, version)
	return args.String(0), args.Error(1)
}

func (m *HttpGatewayMock) VersionExists(ctx context.Context, version string) (bool, error) {
	args := m.Called(ctx, version)
	return args.Bool(0), args.Error(1)
}

func (m *HttpGatewayMock) DownloadVersion(ctx context.Context, action *domain.Action, file *os.File) error {
	args := m.Called(ctx, action, file)
	return args.Error(0)
}
