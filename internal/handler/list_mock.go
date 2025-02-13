package handler

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockListHandler struct {
	mock.Mock
}

func (m *MockListHandler) Handle(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
