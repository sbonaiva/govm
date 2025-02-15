package handler

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UninstallHandlerMock struct {
	mock.Mock
}

func (m *UninstallHandlerMock) Handle(ctx context.Context, uninstall *domain.Uninstall) error {
	args := m.Called(ctx, uninstall)
	return args.Error(0)
}
