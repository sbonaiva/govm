package handler_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/service"
	"github.com/sbonaiva/govm/internal/test"
	"github.com/stretchr/testify/suite"
)

type listHandlerSuite struct {
	suite.Suite
	ctx       context.Context
	sharedSvc *service.SharedServiceMock
	handler   handler.ListHandler
}

func TestListHandler(t *testing.T) {
	suite.Run(t, new(listHandlerSuite))
}

func (r *listHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.sharedSvc = new(service.SharedServiceMock)
	r.handler = handler.NewList(r.sharedSvc)
}

func (r *listHandlerSuite) TearDownTest() {
	r.sharedSvc.AssertExpectations(r.T())
}

func (r *listHandlerSuite) TestSuccess() {
	r.sharedSvc.On("GetAvailableGoVersions", r.ctx).Return(domain.VersionsResponse{
		Versions: []domain.VersionResponse{
			{Version: "1.16"},
			{Version: "1.17"},
			{Version: "1.18"},
			{Version: "1.19"},
			{Version: "1.20"},
			{Version: "1.21"},
		},
	}, nil)

	r.sharedSvc.On("GetInstalledGoVersion", r.ctx).Return("", nil)

	output, err := test.CaptureOutput(func() error {
		err := r.handler.Handle(r.ctx)
		return err
	})

	expected := strings.Join(
		[]string{
			strings.Repeat("=", 100) + "\n",
			"Available Go versions for %s/%s \n",
			strings.Repeat("=", 100) + "\n",
			"1.16           1.17           1.18           1.19           1.20           1.21           \n",
			strings.Repeat("=", 100) + "\n",
			"* currently in use\n",
			strings.Repeat("=", 100) + "\n",
		},
		"",
	)

	r.NoError(err)
	r.Equal(fmt.Sprintf(expected, runtime.GOOS, runtime.GOARCH), output)
}

func (r *listHandlerSuite) TestError() {
	r.sharedSvc.On("GetAvailableGoVersions", r.ctx).Return(domain.VersionsResponse{}, errors.New("error"))

	output, err := test.CaptureOutput(func() error {
		err := r.handler.Handle(r.ctx)
		return err
	})

	r.Error(err)
	r.Equal("error", err.Error())
	r.Empty(output)
}
