package handler

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/mock"
)

type InstallHandlerMock struct {
	mock.Mock
}

func (m *InstallHandlerMock) Handle(ctx context.Context, install *domain.Install) error {
	args := m.Called(ctx, install)
	return args.Error(0)
}
