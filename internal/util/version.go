package util

import (
	"errors"
	"os/exec"
	"strings"
)

func GetInstalledGoVersion() (string, error) {
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
