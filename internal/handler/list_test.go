package handler_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestHandle_Success(t *testing.T) {
	ctx := context.Background()
	mockHttpGateway := new(gateway.MockHttpGateway)
	versions := []domain.GoVersionResponse{
		{Version: "1.16"}, {Version: "1.17"}, {Version: "1.18"},
		{Version: "1.19"}, {Version: "1.20"}, {Version: "1.21"},
	}
	mockHttpGateway.On("GetVersions", ctx).Return(versions, nil)
	listHandler := handler.NewList(mockHttpGateway)

	output, err := util.CaptureOutput(func() error {
		err := listHandler.Handle(ctx)
		return err
	})

	assert.NoError(t, err)
	assert.Contains(t, output, fmt.Sprintf("Available Go versions for %s/%s", runtime.GOOS, runtime.GOARCH))
}

func TestHandle_GetVersionsError(t *testing.T) {
	ctx := context.Background()
	mockHttpGateway := new(gateway.MockHttpGateway)
	mockHttpGateway.On("GetVersions", ctx).Return([]domain.GoVersionResponse{}, errors.New("gateway error"))

	listHandler := handler.NewList(mockHttpGateway)

	err := listHandler.Handle(ctx)

	assert.Error(t, err)
	assert.Equal(t, "gateway error", err.Error())
}
