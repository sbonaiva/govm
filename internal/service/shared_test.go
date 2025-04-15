package service_test

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
	"github.com/sbonaiva/govm/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const (
	goBinDir = "/fake/home/.govm/go/bin"
	bashDir  = "/bin/bash"
	pathEnv  = "/usr/bin:/usr/local/bin"
)

var (
	fileModeType   = mock.AnythingOfType("fs.FileMode")
	arrayUInt8Type = mock.AnythingOfType("[]uint8")
)

type sharedServiceSuite struct {
	suite.Suite
	ctx          context.Context
	action       *domain.Action
	osGateway    *gateway.OsGatewayMock
	httpGateway  *gateway.HttpGatewayMock
	fileInfoMock *gateway.FileInfoMock
	sharedSvc    service.SharedService
}

func TestSharedService(t *testing.T) {
	suite.Run(t, new(sharedServiceSuite))
}

func (r *sharedServiceSuite) TearDownTest() {
	r.osGateway.AssertExpectations(r.T())
	r.httpGateway.AssertExpectations(r.T())
}

func (r *sharedServiceSuite) SetupTest() {
	r.ctx = context.Background()
	r.action = &domain.Action{
		Version:          "1.19.3",
		InstalledVersion: "1.19.2",
	}
	r.osGateway = new(gateway.OsGatewayMock)
	r.httpGateway = new(gateway.HttpGatewayMock)
	r.fileInfoMock = new(gateway.FileInfoMock)
	r.sharedSvc = service.NewShared(r.httpGateway, r.osGateway)
}

func (r *sharedServiceSuite) TestCheckUserHomeSuccess() {
	r.osGateway.On("GetUserHomeDir").Return("/home/fake", nil).Once()

	err := r.sharedSvc.CheckUserHome(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestCheckUserHomeError() {
	r.osGateway.On("GetUserHomeDir").Return("", errors.New("error")).Once()

	err := r.sharedSvc.CheckUserHome(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeCheckUserHome), err)
}

func (r *sharedServiceSuite) TestCheckVersionSuccess() {
	r.httpGateway.On("VersionExists", r.ctx, r.action.Version).Return(true, nil).Once()

	err := r.sharedSvc.CheckVersion(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestCheckVersionError() {
	r.httpGateway.On("VersionExists", r.ctx, r.action.Version).Return(false, errors.New("error")).Once()

	err := r.sharedSvc.CheckVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeCheckVersion), err)
}

func (r *sharedServiceSuite) TestCheckVersionNotExistsError() {
	r.httpGateway.On("VersionExists", r.ctx, r.action.Version).Return(false, nil).Once()

	err := r.sharedSvc.CheckVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewVersionNotAvailableError(r.action.Version), err)
}

func (r *sharedServiceSuite) TestDownloadVersionSuccess() {
	tempFile, _ := os.CreateTemp("", "")

	r.osGateway.On("RemoveDir", r.action.DownloadFile()).Return(nil).Once()
	r.osGateway.On("CreateFile", r.action.DownloadFile()).Return(tempFile, nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, r.action, tempFile).Return(nil).Once()

	err := r.sharedSvc.DownloadVersion(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestDownloadVersionRemoveDirError() {
	r.osGateway.On("RemoveDir", r.action.DownloadFile()).Return(errors.New("error")).Once()

	err := r.sharedSvc.DownloadVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeDownloadRemoveDir), err)
}

func (r *sharedServiceSuite) TestDownloadVersionCreateDirError() {
	var osNilFile *os.File

	r.osGateway.On("RemoveDir", r.action.DownloadFile()).Return(nil).Once()
	r.osGateway.On("CreateFile", r.action.DownloadFile()).Return(osNilFile, errors.New("error")).Once()

	err := r.sharedSvc.DownloadVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeDownloadCreateFile), err)
}

func (r *sharedServiceSuite) TestDownloadVersionError() {
	tempFile, _ := os.CreateTemp("", "")

	r.osGateway.On("RemoveDir", r.action.DownloadFile()).Return(nil).Once()
	r.osGateway.On("CreateFile", r.action.DownloadFile()).Return(tempFile, nil).Once()
	r.httpGateway.On("DownloadVersion", r.ctx, r.action, tempFile).Return(errors.New("error")).Once()

	err := r.sharedSvc.DownloadVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeDownloadVersion), err)
}

func (r *sharedServiceSuite) TestChecksumSuccess() {
	downloadFile, _ := os.CreateTemp("", "")
	checksumFile, _ := os.CreateTemp("", "")
	hash := sha256.New()
	io.Copy(hash, downloadFile)

	r.httpGateway.On("GetChecksum", r.ctx, r.action.Version).Return(fmt.Sprintf("%x", hash.Sum(nil)), nil).Once()
	r.osGateway.On("OpenFile", r.action.DownloadFile()).Return(checksumFile, nil).Once()

	err := r.sharedSvc.Checksum(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestChecksumDownloadError() {
	r.httpGateway.On("GetChecksum", r.ctx, r.action.Version).Return("", errors.New("error")).Once()

	err := r.sharedSvc.Checksum(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeChecksumDownload), err)
}

func (r *sharedServiceSuite) TestChecksumOpenFileError() {
	var osNilFile *os.File

	r.httpGateway.On("GetChecksum", r.ctx, r.action.Version).Return("", nil).Once()
	r.osGateway.On("OpenFile", r.action.DownloadFile()).Return(osNilFile, errors.New("error")).Once()

	err := r.sharedSvc.Checksum(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeChecksumOpenFile), err)
}

func (r *sharedServiceSuite) TestChecksumMismatchError() {
	r.httpGateway.On("GetChecksum", r.ctx, r.action.Version).Return("", nil).Once()
	r.osGateway.On("OpenFile", r.action.DownloadFile()).Return(os.CreateTemp("", "")).Once()

	err := r.sharedSvc.Checksum(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeChecksumMismatch), err)
}

func (r *sharedServiceSuite) TestRemoveVersionSuccess() {
	r.osGateway.On("RemoveDir", r.action.HomeGoDir()).Return(nil).Once()

	err := r.sharedSvc.RemoveVersion(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestRemoveVersionError() {
	r.osGateway.On("RemoveDir", r.action.HomeGoDir()).Return(errors.New("error")).Once()

	err := r.sharedSvc.RemoveVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeRemoveVersion), err)
}

func (r *sharedServiceSuite) TestUntarFilesSuccess() {
	r.osGateway.On("CreateDir", r.action.HomeGovmDir(), fileModeType).Return(nil).Once()
	r.osGateway.On("Untar", r.action.DownloadFile(), r.action.HomeGovmDir()).Return(nil).Once()
	r.osGateway.On("RemoveFile", r.action.DownloadFile()).Return(nil).Once()

	err := r.sharedSvc.UntarFiles(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestUntarFilesCreateDirError() {
	r.osGateway.On("CreateDir", mock.AnythingOfType("string"), fileModeType).Return(errors.New("error")).Once()

	err := r.sharedSvc.UntarFiles(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUntarCreateDir), err)
}

func (r *sharedServiceSuite) TestUntarFilesExtractError() {
	r.osGateway.On("CreateDir", r.action.HomeGovmDir(), fileModeType).Return(nil).Once()
	r.osGateway.On("Untar", r.action.DownloadFile(), r.action.HomeGovmDir()).Return(errors.New("error")).Once()

	err := r.sharedSvc.UntarFiles(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeUntarExtract), err)
}

func (r *sharedServiceSuite) TestAddToPathStatError() {
	r.fileInfoMock.On("IsDir").Return(true).Once()
	r.osGateway.On("GetEnv", "PATH").Return(pathEnv, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(r.fileInfoMock, errors.New("error")).Once()

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeAddToPathStat), err)
}

func (r *sharedServiceSuite) TestAddToPathReadFileError() {
	r.fileInfoMock.On("IsDir").Return(true).Once()
	r.osGateway.On("GetEnv", "PATH").Return(pathEnv, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(r.fileInfoMock, nil).Once()
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte{}, errors.New("error")).Once()

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeAddToPathRead), err)
}

func (r *sharedServiceSuite) TestAddToPathWriteFileError() {
	r.fileInfoMock.On("IsDir").Return(true).Once()
	r.osGateway.On("GetEnv", "PATH").Return(pathEnv, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(r.fileInfoMock, nil).Once()
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("content"), nil).Once()
	r.osGateway.On("WriteFile", mock.AnythingOfType("string"), arrayUInt8Type, fileModeType).Return(errors.New("error")).Once()

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeAddToPathWrite), err)
}

func (r *sharedServiceSuite) TestAddToPathNoShellsFoundError() {
	r.fileInfoMock.On("IsDir").Return(true).Once()
	r.osGateway.On("GetEnv", "PATH").Return(pathEnv, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(r.fileInfoMock, errors.New("error")).Times(8)

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.Equal(domain.NewUnexpectedError(domain.ErrCodeAddToPathNoShellsFound), err)
}

func (r *sharedServiceSuite) TestAddToPathWithEmptyShellEnvVarSuccess() {
	r.fileInfoMock.On("IsDir").Return(true).Once()
	r.osGateway.On("GetEnv", "PATH").Return(pathEnv, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(r.fileInfoMock, nil).Times(8)
	r.osGateway.On("WriteFile", mock.AnythingOfType("string"), arrayUInt8Type, fileModeType).Return(nil).Times(8)
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("export PATH=$PATH:/home/fake/go/bin"), nil).Times(8)

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestAddToPathWithFilledShellEnvVarSuccess() {
	r.fileInfoMock.On("IsDir").Return(true).Once()
	r.osGateway.On("GetEnv", "PATH").Return(pathEnv, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.AnythingOfType("string")).Return(r.fileInfoMock, nil).Once()
	r.osGateway.On("WriteFile", mock.AnythingOfType("string"), arrayUInt8Type, fileModeType).Return(nil).Once()
	r.osGateway.On("ReadFile", mock.AnythingOfType("string")).Return([]byte("export PATH=$PATH:/home/fake/go/bin"), nil).Once()

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestAddToPatWithGoAlreadyInPathSuccess() {
	r.osGateway.On("GetEnv", "PATH").Return("/home/fake/.govm/go/bin", nil).Once()

	err := r.sharedSvc.AddToPath(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestRemoveFromPathNoShellCommandsFoundError() {
	r.fileInfoMock.On("IsDir").Return(true)
	r.osGateway.On("GetEnv", "PATH").Return(goBinDir, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(r.fileInfoMock, errors.New("error")).Times(8)

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathNoShellsFound), err)
}

func (r *sharedServiceSuite) TestRemoveFromPathStatError() {
	r.fileInfoMock.On("IsDir").Return(true)
	r.osGateway.On("GetEnv", "PATH").Return(goBinDir, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(r.fileInfoMock, errors.New("error")).Once()

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathStat), err)
}

func (r *sharedServiceSuite) TestRemoveFromShellRunCommandsReadError() {
	r.fileInfoMock.On("IsDir").Return(true)
	r.osGateway.On("GetEnv", "PATH").Return(goBinDir, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(r.fileInfoMock, nil).Once()
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte{}, errors.New("error")).Once()

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathRead), err)
}

func (r *sharedServiceSuite) TestRemoveFromShellRunCommandsWriteError() {
	r.fileInfoMock.On("IsDir").Return(true)
	r.osGateway.On("GetEnv", "PATH").Return(goBinDir, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(r.fileInfoMock, nil).Once()
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte("content"), nil).Once()
	r.osGateway.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error")).Once()

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathWrite), err)
}

func (r *sharedServiceSuite) TestSuccessAlreadyRemovedFromPath() {
	r.osGateway.On("GetEnv", "PATH").Return("/usr/bin", nil).Once()

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestSuccessRemovingFromPathWithEmptyShellEnvVar() {
	r.fileInfoMock.On("IsDir").Return(true)
	r.osGateway.On("GetEnv", "PATH").Return(goBinDir, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return("", nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(r.fileInfoMock, nil).Times(8)
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte("path content"), nil).Times(8)
	r.osGateway.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(8)

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestSuccessRemovingFromPathWithFilledShellEnvVar() {
	r.fileInfoMock.On("IsDir").Return(true)
	r.osGateway.On("GetEnv", "PATH").Return(goBinDir, nil).Once()
	r.osGateway.On("GetEnv", "SHELL").Return(bashDir, nil).Once()
	r.osGateway.On("Stat", mock.Anything).Return(r.fileInfoMock, nil).Once()
	r.osGateway.On("ReadFile", mock.Anything).Return([]byte("path content"), nil).Once()
	r.osGateway.On("WriteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := r.sharedSvc.RemoveFromPath(r.ctx, r.action)

	r.NoError(err)
}

func (r *sharedServiceSuite) TestCheckInstalledVersionSuccess() {
	r.osGateway.On("GetInstalledGoVersion").Return("1.19.2", nil).Once()

	err := r.sharedSvc.CheckInstalledVersion(r.ctx, r.action)

	r.NoError(err)
	r.Equal("1.19.2", r.action.InstalledVersion)
}

func (r *sharedServiceSuite) TestCheckInstalledVersionError() {
	r.osGateway.On("GetInstalledGoVersion").Return("", errors.New("error")).Once()

	err := r.sharedSvc.CheckInstalledVersion(r.ctx, r.action)

	r.Error(err)
	r.Equal(domain.NewNoGoInstallationsFoundError(), err)
}

func (r *sharedServiceSuite) TestCheckAvailableUpdatesSuccess() {

	tests := []struct {
		updateStrategy  domain.UpdateStrategy
		expectedVersion string
		expectedError   error
	}{
		{
			updateStrategy:  domain.PatchStrategy,
			expectedVersion: "1.19.3",
		},
		{
			updateStrategy:  domain.MinorStrategy,
			expectedVersion: "1.24.2",
		},
		{
			updateStrategy:  domain.MajorStrategy,
			expectedVersion: "2.0.0",
		},
		{
			updateStrategy: domain.PatchStrategy,
			expectedError:  domain.NewNoUpdatesAvailableError(domain.PatchStrategy, "1.19.2"),
		},
		{
			updateStrategy: domain.MinorStrategy,
			expectedError:  domain.NewNoUpdatesAvailableError(domain.MinorStrategy, "1.19.2"),
		},
		{
			updateStrategy: domain.MajorStrategy,
			expectedError:  domain.NewNoUpdatesAvailableError(domain.MajorStrategy, "1.19.2"),
		},
	}

	for _, tc := range tests {
		r.Run(string(tc.updateStrategy), func() {
			action := &domain.Action{
				InstalledVersion: "1.19.2",
				UpdateStrategy:   tc.updateStrategy,
			}

			var availableVersions domain.VersionsResponse

			if tc.expectedError == nil {
				availableVersions = domain.VersionsResponse{
					Versions: []domain.VersionResponse{
						{
							Version: "2.0.0",
							Stable:  true,
						},
						{
							Version: "1.24.2",
							Stable:  true,
						},
						{
							Version: "1.23.8",
							Stable:  true,
						},
						{
							Version: "1.19.3",
							Stable:  true,
						},
						{
							Version: "1.19.2",
							Stable:  true,
						},
					},
				}
			} else {
				availableVersions = domain.VersionsResponse{
					Versions: []domain.VersionResponse{
						{
							Version: "1.19.2",
							Stable:  true,
						},
					},
				}
			}

			r.httpGateway.On("GetVersions", r.ctx).Return(availableVersions, nil).Once()

			err := r.sharedSvc.CheckAvailableUpdates(r.ctx, action)

			if tc.expectedError == nil {
				r.NoError(err)
				r.Equal(tc.expectedVersion, action.Version)
			} else {
				r.Error(err)
				r.Equal(tc.expectedError, err)
				r.Empty(action.Version)
			}
		})
	}
}

func (r *sharedServiceSuite) TestCheckAvailableUpdatesGetVersionsError() {

	action := &domain.Action{
		InstalledVersion: "1.19.2",
		UpdateStrategy:   domain.PatchStrategy,
	}

	r.httpGateway.On("GetVersions", r.ctx).Return(domain.VersionsResponse{}, errors.New("error")).Once()

	err := r.sharedSvc.CheckAvailableUpdates(r.ctx, action)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeListVersions), err)
	r.Empty(action.Version)
}

func (r *sharedServiceSuite) TestGetAvailableGoVersionsSuccess() {
	r.httpGateway.On("GetVersions", r.ctx).Return(domain.VersionsResponse{}, nil).Once()

	available, err := r.sharedSvc.GetAvailableGoVersions(r.ctx)

	r.NoError(err)
	r.Empty(available.Versions)
}

func (r *sharedServiceSuite) TestGetAvailableGoVersionsError() {
	r.httpGateway.On("GetVersions", r.ctx).Return(domain.VersionsResponse{}, errors.New("error")).Once()

	available, err := r.sharedSvc.GetAvailableGoVersions(r.ctx)

	r.Error(err)
	r.Equal(domain.NewUnexpectedError(domain.ErrCodeListVersions), err)
	r.Empty(available.Versions)
}

func (r *sharedServiceSuite) TestGetInstalledGoVersionSuccess() {
	r.osGateway.On("GetInstalledGoVersion").Return("1.2.2", nil).Once()

	installed, err := r.sharedSvc.GetInstalledGoVersion(r.ctx)

	r.NoError(err)
	r.Equal("1.2.2", installed)
}

func (r *sharedServiceSuite) TestGetInstalledGoVersionError() {
	r.osGateway.On("GetInstalledGoVersion").Return("", errors.New("error")).Once()

	installed, err := r.sharedSvc.GetInstalledGoVersion(r.ctx)

	r.Error(err)
	r.Equal(domain.NewNoGoInstallationsFoundError(), err)
	r.Empty(installed)
}
