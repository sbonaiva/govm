package handler_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
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
	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(false, nil).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewVersionNotAvailableError(install.Version), err)
}

func (r *installHandlerSuite) TestUserHomeError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return("", errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallCheckUserHome), err)
}

func (r *installHandlerSuite) TestVersionExistsError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(false, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallCheckVersion), err)
}

func (r *installHandlerSuite) TestDownloadRemoveDirError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallDownloadRemoveDir), err)
}

func (r *installHandlerSuite) TestDownloadCreateDirError() {
	var osNilFile *os.File
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(osNilFile, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallDownloadCreateFile), err)
}

func (r *installHandlerSuite) TestDownloadVersionError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(os.CreateTemp("", "")).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallDownloadVersion), err)
}

func (r *installHandlerSuite) TestChecksumDownloadError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(os.CreateTemp("", "")).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return("", errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallChecksumDownload), err)
}

func (r *installHandlerSuite) TestChecksumOpenFileError() {
	var osNilFile *os.File
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(os.CreateTemp("", "")).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return("", nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(osNilFile, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallChecksumOpenFile), err)
}

func (r *installHandlerSuite) TestChecksumMismatchError() {
	install := &domain.Install{Version: "1.16"}

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(os.CreateTemp("", "")).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return("", nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(os.CreateTemp("", "")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallChecksumMismatch), err)
}

func (r *installHandlerSuite) TestRemovePreviousVersionError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallRemovePreviousVersion), err)
}

func (r *installHandlerSuite) TestUntarCreateDirError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallUntarCreateDir), err)
}

func (r *installHandlerSuite) TestUntarExtractError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallUntarExtract), err)
}

func (r *installHandlerSuite) TestAddToPathStatError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin:/usr/local/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fileInfoMmock, errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathStat), err)
}

func (r *installHandlerSuite) TestAddToPathReadFileError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin:/usr/local/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fileInfoMmock, nil).Once()
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte{}, errors.New("error")).Once()
	//r.osGateway.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathRead), err)
}

func (r *installHandlerSuite) TestAddToPathWriteFileError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin:/usr/local/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fileInfoMmock, nil).Once()
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("content"), nil).Once()
	r.osGateway.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(errors.New("error")).Once()

	err := r.handler.Handle(r.ctx, install)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathWrite), err)
}

func (r *installHandlerSuite) TestAddToPathWriteNotShellsFoundError() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin:/usr/local/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fileInfoMmock, errors.New("error")).Times(8)

	err := r.handler.Handle(r.ctx, install)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathNoShellsFound), err)
}

func (r *installHandlerSuite) TestSuccessWithEmptyShellEnvVar() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin:/usr/local/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fileInfoMmock, nil).Times(8)
	r.osGateway.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil).Times(8)
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("export PATH=$PATH:/home/fake/go/bin"), nil).Times(8)

	err := r.handler.Handle(r.ctx, install)

	r.NoError(err)
}

func (r *installHandlerSuite) TestSuccessWithFilledShellEnvVar() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin:/usr/local/bin", nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("/bin/bash", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(fileInfoMmock, nil).Once()
	r.osGateway.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("export PATH=$PATH:/home/fake/go/bin"), nil).Once()

	err := r.handler.Handle(r.ctx, install)

	r.NoError(err)
}

func (r *installHandlerSuite) TestSuccessWithGoAlreadyInPath() {
	install := &domain.Install{Version: "1.16"}

	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)
	fileInfoMmock := new(gateway.FileInfoMock)
	fileInfoMmock.On("IsDir").Return(true)

	r.osGateway.On("GetUserHomeDir").Return(r.fakeHomeDir, nil).Once()
	r.httpGateway.On("VersionExists", r.ctx, install.Version).Return(true, nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("CreateFile", mock.AnythingOfType("string")).Return(downloadFile, nil).Once()
	r.osGateway.On("RemoveFile", mock.AnythingOfType("string")).Return(nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	r.httpGateway.On("GetChecksum", r.ctx, install.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", mock.AnythingOfType("string")).Return(checksumFile, nil).Once()
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), mock.AnythingOfType("fs.FileMode")).Return(nil).Once()
	r.osGateway.On("RemoveDir", mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("Untar", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
	r.osGateway.On("GetEnv", "PATH").Return("/home/fake/.govm/go/bin", nil).Once()

	err := r.handler.Handle(r.ctx, install)

	r.NoError(err)
}
