package gateway

import (
	"io/fs"
	"os"
	"time"

	"github.com/stretchr/testify/mock"
)

type OsGatewayMock struct {
	mock.Mock
}

func (m *OsGatewayMock) GetUserHomeDir() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *OsGatewayMock) CreateDir(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *OsGatewayMock) RemoveDir(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *OsGatewayMock) CreateFile(path string) (*os.File, error) {
	args := m.Called(path)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *OsGatewayMock) OpenFile(path string) (*os.File, error) {
	args := m.Called(path)
	return args.Get(0).(*os.File), args.Error(1)
}

func (m *OsGatewayMock) Stat(path string) (os.FileInfo, error) {
	args := m.Called(path)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

func (m *OsGatewayMock) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *OsGatewayMock) WriteFile(path string, data []byte, perm os.FileMode) error {
	args := m.Called(path, data, perm)
	return args.Error(0)
}

func (m *OsGatewayMock) RemoveFile(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *OsGatewayMock) GetEnv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func (m *OsGatewayMock) Untar(source string, target string) error {
	args := m.Called(source, target)
	return args.Error(0)
}

func (m *OsGatewayMock) GetInstalledGoVersion() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type FileInfoMock struct {
	mock.Mock
}

func (f *FileInfoMock) IsDir() bool {
	args := f.Called()
	return args.Bool(0)
}

func (f *FileInfoMock) ModTime() time.Time {
	args := f.Called()
	return args.Get(0).(time.Time)
}

func (f *FileInfoMock) Mode() fs.FileMode {
	args := f.Called()
	return args.Get(0).(fs.FileMode)
}

func (f *FileInfoMock) Name() string {
	args := f.Called()
	return args.String(0)
}

func (f *FileInfoMock) Size() int64 {
	args := f.Called()
	return args.Get(0).(int64)
}

func (f *FileInfoMock) Sys() any {
	args := f.Called()
	return args.Get(0)
}
