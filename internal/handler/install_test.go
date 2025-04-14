package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/service"
	"github.com/stretchr/testify/suite"
)

type installHandlerSuite struct {
	suite.Suite
	ctx       context.Context
	action    *domain.Action
	sharedSvc *service.SharedServiceMock
	handler   handler.InstallHandler
}

func TestInstallHandler(t *testing.T) {
	suite.Run(t, new(installHandlerSuite))
}

func (r *installHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.action = &domain.Action{
		Version: "1.20.5",
		HomeDir: "/home/fake",
	}
	r.sharedSvc = new(service.SharedServiceMock)
	r.handler = handler.NewInstall(r.sharedSvc)
}

func (r *installHandlerSuite) TearDownTest() {
	r.sharedSvc.AssertExpectations(r.T())
}

func (r *installHandlerSuite) TestSuccess() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("UntarFiles", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("AddToPath", r.ctx, r.action).Return(nil)

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.NoError(err)
}

func (r *installHandlerSuite) TestCheckUserHomeError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *installHandlerSuite) TestCheckVersionError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *installHandlerSuite) TestDownloadVersionError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *installHandlerSuite) TestChecksumError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *installHandlerSuite) TestRemoveVersionError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *installHandlerSuite) TestUntarFilesError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("UntarFiles", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *installHandlerSuite) TestAddToPathError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("UntarFiles", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("AddToPath", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}
