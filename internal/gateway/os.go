package gateway

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

type OsGateway interface {
	GetUserHomeDir() (string, error)
	Stat(path string) (os.FileInfo, error)
	CreateDir(path string, perm os.FileMode) error
	RemoveDir(path string) error
	CreateFile(path string) (*os.File, error)
	OpenFile(path string) (*os.File, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	RemoveFile(path string) error
	GetEnv(key string) string
	Untar(source string, target string) error
	GetInstalledGoVersion() (string, error)
}

type osClient struct{}

func NewOsGateway() OsGateway {
	return &osClient{}
}

func (o *osClient) GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func (o *osClient) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (o *osClient) CreateDir(path string, perm os.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (o *osClient) RemoveDir(path string) error {
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (o *osClient) CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

func (o *osClient) OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

func (o *osClient) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (o *osClient) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (o *osClient) RemoveFile(path string) error {
	return os.Remove(path)
}

func (o *osClient) GetEnv(key string) string {
	return os.Getenv(key)
}

func (o *osClient) Untar(source string, target string) error {
	return exec.Command("tar", "-C", target, "-xzf", source).Run()
}

func (o *osClient) GetInstalledGoVersion() (string, error) {
	goPath, err := exec.LookPath("go")
	if err != nil {
		return "", err
	}

	outputBytes, err := exec.Command(goPath, "version").CombinedOutput()
	if err != nil {
		return "", err
	}

	outputParts := strings.Split(string(outputBytes), " ")
	if len(outputParts) < 3 {
		return "", errors.New("unexpected go version command output")
	}

	return outputParts[2], nil
}
