package util

import (
	"errors"
	"os/exec"
	"strings"
)

func GetInstalledGoVersion() (string, error) {
	outputBytes, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		return "", err
	}

	outputParts := strings.Split(string(outputBytes), " ")
	if len(outputParts) < 3 {
		return "", errors.New("unexpected go version command output")
	}

	return outputParts[2], nil
}
