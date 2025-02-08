package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/sbonaiva/govm/internal/domain"
)

type Uninstall interface {
	Execute(ctx context.Context) error
}

type uninstall struct {
}

func NewUninstall() Uninstall {
	return &uninstall{}
}

func (r *uninstall) Execute(ctx context.Context) error {
	slog.Info("Uninstalling current Go version", slog.String("Uninstall", "Execute"))

	if _, err := os.Stat("/usr/local/go"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	//TODO: Implement prompt to confirm uninstall

	if err := r.removeCurrentVersion(); err != nil {
		return err
	}

	if err := r.removeFromPath(); err != nil {
		return err
	}

	return nil
}

func (r *uninstall) removeCurrentVersion() error {
	if err := exec.Command("sudo", "rm", "-rf", "/usr/local/go").Run(); err != nil {
		return err
	}
	return nil
}

func (r *uninstall) removeFromPath() error {
	if path := os.Getenv("PATH"); !strings.Contains(path, "/usr/local/go/bin") {
		return nil
	}

	if err := r.removeFromShellRunCommands(); err != nil {
		return err
	}

	return nil
}

func (r *uninstall) removeFromShellRunCommands() error {

	succeded := 0

	for _, rcf := range domain.ShellRunCommandsFiles {

		if _, err := os.Stat(rcf); err != nil {
			continue
		}

		oldContent, err := os.ReadFile(rcf)
		if err != nil {
			return err
		}

		// TODO: Implement prompt to confirm uninstall
		newContent := []byte(strings.ReplaceAll(string(oldContent), "export", ""))

		if err := os.WriteFile(rcf, newContent, 0644); err != nil {
			return err
		}

		if err := exec.Command("source", rcf).Run(); err != nil {
			return err
		}

		succeded++
	}

	if succeded == 0 {
		return fmt.Errorf("Could not find any shell run commands file")
	}

	return nil
}
