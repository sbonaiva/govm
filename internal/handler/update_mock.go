package handler

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/mock"
)

type UpdateHandlerMock struct {
	mock.Mock
}

func (m *UpdateHandlerMock) Handle(ctx context.Context, update *domain.Action) error {
	args := m.Called(ctx, update)
	return args.Error(0)
}
