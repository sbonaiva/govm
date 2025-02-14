package handler

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type ListHandlerMock struct {
	mock.Mock
}

func (m *ListHandlerMock) Handle(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
