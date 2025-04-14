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

type uninstallHandlerSuite struct {
	suite.Suite
	ctx       context.Context
	action    *domain.Action
	sharedSvc *service.SharedServiceMock
	handler   handler.UninstallHandler
}

func TestUninstallHandler(t *testing.T) {
	suite.Run(t, new(uninstallHandlerSuite))
}

func (r *uninstallHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.action = &domain.Action{
		Version: "1.20.5",
		HomeDir: "/home/fake",
	}
	r.sharedSvc = new(service.SharedServiceMock)
	r.handler = handler.NewUninstall(r.sharedSvc)
}

func (r *uninstallHandlerSuite) TearDownTest() {
	r.sharedSvc.AssertExpectations(r.T())
}

func (r *uninstallHandlerSuite) TestSuccess() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveFromPath", r.ctx, r.action).Return(nil)

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.NoError(err)
}

func (r *uninstallHandlerSuite) TestCheckUserHomeError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *uninstallHandlerSuite) TestCheckVersionError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *uninstallHandlerSuite) TestRemoveVersionError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}

func (r *uninstallHandlerSuite) TestRemoveFromPathError() {
	// Arrange
	r.sharedSvc.On("CheckUserHome", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("CheckVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveVersion", r.ctx, r.action).Return(nil)
	r.sharedSvc.On("RemoveFromPath", r.ctx, r.action).Return(errors.New("error"))

	// Act
	err := r.handler.Handle(r.ctx, r.action)

	// Assert
	r.Error(err)
	r.Equal("error", err.Error())
}
