package handler

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUseHandler struct {
	mock.Mock
}

func (m *MockUseHandler) Handle(ctx context.Context, use *domain.Use) error {
	args := m.Called(ctx, use)
	return args.Error(0)
}
