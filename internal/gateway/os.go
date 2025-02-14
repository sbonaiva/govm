package gateway

import (
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
