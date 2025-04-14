package service

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/mock"
)

type SharedServiceMock struct {
	mock.Mock
}

func (m *SharedServiceMock) CheckUserHome(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) CheckVersion(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) DownloadVersion(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) Checksum(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) RemoveVersion(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) UntarFiles(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) AddToPath(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) RemoveFromPath(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) CheckInstalledVersion(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}

func (m *SharedServiceMock) CheckAvailableUpdates(ctx context.Context, action *domain.Action) error {
	args := m.Called(ctx, action)
	return args.Error(0)
}
