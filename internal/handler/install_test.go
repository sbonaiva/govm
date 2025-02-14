package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type installHandlerSuite struct {
	suite.Suite
	ctx         context.Context
	fakeHomeDir string
	httpGateway *gateway.HttpGatewayMock
	osGateway   *gateway.OsGatewayMock
	handler     handler.InstallHandler
}

func TestInstallHandler(t *testing.T) {
	suite.Run(t, new(installHandlerSuite))
}

func (r *installHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.fakeHomeDir = "/home/fake"
	r.httpGateway = new(gateway.HttpGatewayMock)
	r.osGateway = new(gateway.OsGatewayMock)
	r.handler = handler.NewInstall(r.httpGateway, r.osGateway)
}

func (r *installHandlerSuite) TearDownTest() {
	r.httpGateway.AssertExpectations(r.T())
}

func (r *installHandlerSuite) TestVersionNotAvailable() {
	install := &domain.Install{Version: "1.16"}
	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil)
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(false, nil)

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrVersionNotAvailable(install.Version), err)
}

func (r *installHandlerSuite) TestUserHomeError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return("", errors.New("error"))

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrUnexpected, err)
}

func (r *installHandlerSuite) TestDownloadError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil)
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil)
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(domain.ErrUnexpected)

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrUnexpected, err)
}

func (r *installHandlerSuite) TestChecksumError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil)
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil)
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil)
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return("", nil)

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.ErrUnexpected, err)
}
