package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/test"
	"github.com/stretchr/testify/suite"
)

type listHandlerSuite struct {
	suite.Suite
	ctx     context.Context
	gateway *gateway.HttpGatewayMock
	handler handler.ListHandler
}

func TestListHandler(t *testing.T) {
	suite.Run(t, new(listHandlerSuite))
}

func (r *listHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.gateway = new(gateway.HttpGatewayMock)
	r.handler = handler.NewList(r.gateway)
}

func (r *listHandlerSuite) TearDownTest() {
	r.gateway.AssertExpectations(r.T())
}

func (r *listHandlerSuite) TestSuccess() {
	versions := []domain.GoVersionResponse{
		{Version: "1.16"}, {Version: "1.17"}, {Version: "1.18"},
		{Version: "1.19"}, {Version: "1.20"}, {Version: "1.21"},
	}
	r.gateway.On("GetVersions", r.ctx).Return(versions, nil)

	output, err := test.CaptureOutput(func() error {
		err := r.handler.Handle(r.ctx)
		return err
	})

	r.NoError(err)
	r.Equal("====================================================================================================\nAvailable Go versions for linux/amd64 \n====================================================================================================\n1.16           1.17           1.18           1.19           1.20           1.21           \n====================================================================================================\n* currently in use\n====================================================================================================\n", output)
}

func (r *listHandlerSuite) TestError() {
	r.gateway.On("GetVersions", r.ctx).Return([]domain.GoVersionResponse{}, errors.New("gateway error"))

	output, err := test.CaptureOutput(func() error {
		err := r.handler.Handle(r.ctx)
		return err
	})

	r.Error(err)
	r.EqualError(err, "gateway error")
	r.Empty(output)
}
