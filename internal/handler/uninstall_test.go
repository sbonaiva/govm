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

type uninstallHandlerSuite struct {
	suite.Suite
	ctx       context.Context
	osGateway *gateway.OsGatewayMock
	handler   handler.UninstallHandler
}

func TestUninstallHandler(t *testing.T) {
	suite.Run(t, new(uninstallHandlerSuite))
}

func (r *uninstallHandlerSuite) SetupTest() {
	r.ctx = context.Background()
	r.osGateway = new(gateway.OsGatewayMock)
	r.handler = handler.NewUninstall(r.osGateway)
}

func (r *uninstallHandlerSuite) TearDownTest() {
	r.osGateway.AssertExpectations(r.T())
}

func (r *uninstallHandlerSuite) TestCheckUserHomeError() {
	uninstall := &domain.Uninstall{}

	r.osGateway.On("GetUserHomeDir").Return("", errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallCheckUserHome), err)
}

func (r *uninstallHandlerSuite) TestCheckVersionStatError() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallCheckVersionStat), err)
}

func (r *uninstallHandlerSuite) TestCheckVersionNotDirError() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(false)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallCheckVersionNotDir), err)
}

func (r *uninstallHandlerSuite) TestRemoveCurrentVersionError() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveDir), err)
}

func (r *uninstallHandlerSuite) TestRemoveFromPathNoShellCommandsFound() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/fake/home/.govm/go/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(fakeFileInfo, errors.New("error")).Times(8)

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathNoShellsFound), err)
}

func (r *uninstallHandlerSuite) TestRemoveFromPathStatError() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/fake/home/.govm/go/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(fakeFileInfo, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathStat), err)
}

func (r *uninstallHandlerSuite) TestRemoveFromShellRunCommandsReadError() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/fake/home/.govm/go/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte{}, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathRead), err)
}

func (r *uninstallHandlerSuite) TestRemoveFromShellRunCommandsWriteError() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/fake/home/.govm/go/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte("content"), nil).Once()
	r.osGateway.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathWrite), err)
}

func (r *uninstallHandlerSuite) TestSuccessAlreadyRemovedFromPath() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin", nil).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.NoError(err)
}

func (r *uninstallHandlerSuite) TestSuccessRemovingFromPathWithEmptyShellEnvVar() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/fake/home/.govm/go/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(fakeFileInfo, nil).Times(8)
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte("path content"), nil).Times(8)
	r.osGateway.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(8)

	err := r.handler.Handle(r.ctx, uninstall)

	r.NoError(err)
}

func (r *uninstallHandlerSuite) TestSuccessRemovingFromPathWithFilledShellEnvVar() {
	uninstall := &domain.Uninstall{}
	fakeFileInfo := new(gateway.FileInfoMock)
	fakeFileInfo.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return("/fake/home", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/fake/home/.govm/go/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(fakeFileInfo, nil).Once()
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte("path content"), nil).Once()
	r.osGateway.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := r.handler.Handle(r.ctx, uninstall)

	r.NoError(err)
}
