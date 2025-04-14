package handler_test

import (
	"context"
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
