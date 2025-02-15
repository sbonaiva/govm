package gateway

import (
	"os"
	"os/user"
)

type osClient struct{}

func (o *osClient) GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func (o *osClient) CreateDir(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o *osClient) RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func (o *osClient) CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

func (o *osClient) OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

func (o *osClient) RemoveFile(path string) error {
	return os.Remove(path)
}

func (o *osClient) GetEnv(key string) string {
	return os.Getenv(key)
}
