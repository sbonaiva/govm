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

type updateHandlerSuite struct {
	suite.Suite
	ctx       context.Context
	action    *domain.Action
	sharedSvc *service.SharedServiceMock
	handler   handler.UpdateHandler
}

func TestUpdateHandler(t *testing.T) {
	suite.Run(t, new(updateHandlerSuite))
}

func (r *updateHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.action = &domain.Action{
		Version:        "1.20.5",
		HomeDir:        "/home/fake",
		UpdateStrategy: domain.PatchStrategy,
	}
	r.sharedSvc = new(service.SharedServiceMock)
	r.handler = handler.NewUpdate(r.sharedSvc)
}

func (r *updateHandlerSuite) TearDownTest() {
	r.sharedSvc.AssertExpectations(r.T())
}

func (r *updateHandlerSuite) TestSuccess() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("UntarFiles", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("AddToPath", r.ctx, r.action).Return(nil)

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.NoError(err)
	r.Equal(r.action.Version, version)
}

func (r *updateHandlerSuite) TestCheckUpdateStrategyError() {
	// Arrange
	action := &domain.Action{
		Version:        "1.20.5",
		HomeDir:        "/home/fake",
		UpdateStrategy: domain.UpdateStrategy("invalid"),
	}
	// Act
	version, err := r.handler.Handle(r.ctx, action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal(domain.NewInvalidUpdateStrategyError(action.UpdateStrategy), err)
}

func (r *updateHandlerSuite) TestCheckInstalledVersionError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestCheckAvailableUpdatesError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestCheckUserHomeError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestCheckVersionError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestDownloadVersionError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestChecksumError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestRemoveVersionError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestUntarFilesError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("UntarFiles", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}

func (r *updateHandlerSuite) TestAddToPathError() {
	// Arrange
	r.sharedSvc.On("CheckInstalledVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckAvailableUpdates", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("DownloadVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("Checksum", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("UntarFiles", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("AddToPath", r.ctx, r.action).Return(errors.New("error"))

	// Act
	version, err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("", version)
	r.Equal("error", err.Error())
}
