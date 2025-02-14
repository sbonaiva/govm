package handler_test

import (
	"context"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type installHandlerSuite struct {
	suite.Suite
	ctx     context.Context
	gateway *gateway.HttpGatewayMock
	handler handler.InstallHandler
}

func TestInstallHandler(t *testing.T) {
	suite.Run(t, new(installHandlerSuite))
}

func (r *installHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.gateway = new(gateway.HttpGatewayMock)
	r.handler = handler.NewInstall(r.gateway)
}

func (r *installHandlerSuite) TearDownTest() {
	r.gateway.AssertExpectations(r.T())
}

func (r *installHandlerSuite) TestVersionNotAvailable() {
	install := &domain.Install{Version: "1.16"}

	r.gateway.On("VersionExists", r.ctx, install.Version).Return(false, nil)

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrVersionNotAvailable(install.Version), err)
}

func (r *installHandlerSuite) TestDownloadError() {
	install := &domain.Install{Version: "1.16"}

	r.gateway.On("VersionExists", r.ctx, install.Version).Return(true, nil)
	r.gateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(domain.ErrUnexpected)

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrUnexpected, err)
}

func (r *installHandlerSuite) TestChecksumError() {
	install := &domain.Install{Version: "1.16"}

	r.gateway.On("VersionExists", r.ctx, install.Version).Return(true, nil)
	r.gateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil)
	r.gateway.On("GetChecksum", r.ctx, install.Version).Return("", nil)

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrUnexpected, err)
}
