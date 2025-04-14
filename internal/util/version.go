package util

import (
	"errors"
	"os/exec"
	"strings"
)

func GetInstalledGoVersion() (string, error) {
	outputBytes, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", errors.New("no go version found")
	}

	outputParts := strings.Split(string(outputBytes), " ")
	if len(outputParts) < 3 {
		return "", errors.New("unexpected output format")
	}

	return outputParts[2], nil
}
