package handler

import (
	"context"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
)

type UninstallHandler interface {
	Handle(ctx context.Context, uninstall *domain.Uninstall) error
}

type uninstallHandler struct{}

func NewUninstall() UninstallHandler {
	return &uninstallHandler{}
}

func (r *uninstallHandler) Handle(ctx context.Context, uninstall *domain.Uninstall) error {
	slog.InfoContext(ctx, "Uninstalling Go version", slog.String("Uninstall", "Execute"))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	defer spn.Stop()
	spn.Start()

	steps := []struct {
		message string
		action  func() error
	}{
		{" Getting home...", func() error { return r.checkUserHome(ctx, uninstall) }},
		{" Checking version...", func() error { return r.checkVersion(ctx, uninstall) }},
		{" Removing current version...", func() error { return r.removeCurrentVersion(ctx, uninstall) }},
		{" Removing from path...", func() error { return r.removeFromPath(ctx, uninstall) }},
	}

	for _, step := range steps {
		spn.Suffix = step.message
		if err := step.action(); err != nil {
			return err
		}
	}

	return nil
}

func (r *uninstallHandler) checkUserHome(ctx context.Context, uninstall *domain.Uninstall) error {
	usr, err := user.Current()
	if err != nil {
		slog.ErrorContext(ctx, "Getting current user", slog.String("Uninstall", "checkUserHome"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	uninstall.HomeDir = usr.HomeDir
	return nil
}

func (r *uninstallHandler) checkVersion(ctx context.Context, uninstall *domain.Uninstall) error {

	fs, err := os.Stat(uninstall.HomeGoDir())
	if err != nil {
		slog.ErrorContext(ctx, "Checking version", slog.String("Uninstall", "checkVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	if !fs.IsDir() {
		slog.ErrorContext(ctx, "Checking version", slog.String("Uninstall", "checkVersion"), slog.String("error", "Go directory not found"))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *uninstallHandler) removeCurrentVersion(ctx context.Context, uninstall *domain.Uninstall) error {
	if err := os.RemoveAll(uninstall.HomeGoDir()); err != nil {
		slog.ErrorContext(ctx, "Removing current version", slog.String("Uninstall", "removeCurrentVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	return nil
}

func (r *uninstallHandler) removeFromPath(ctx context.Context, uninstall *domain.Uninstall) error {
	if path := os.Getenv("PATH"); !strings.Contains(path, uninstall.HomeGoBinDir()) {
		slog.InfoContext(ctx, "Go is already removed from PATH", slog.String("Uninstall", "removeFromPath"))
		return nil
	}

	if shell := os.Getenv("SHELL"); shell != "" {
		if rcf, exists := domain.ShellRunCommandsFiles[shell]; exists {
			return r.removeFromShellRunCommands(ctx, uninstall, rcf)
		}
	}

	succeded := 0
	for _, rcf := range domain.ShellRunCommandsFiles {
		err := r.removeFromShellRunCommands(ctx, uninstall, rcf)
		if err == nil {
			succeded++
		}
	}

	if succeded == 0 {
		slog.ErrorContext(ctx, "No shell rc file found", slog.String("Uninstall", "removeFromPath"))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *uninstallHandler) removeFromShellRunCommands(ctx context.Context, uninstall *domain.Uninstall, rcf string) error {
	rcfPath := filepath.Join(uninstall.HomeDir, rcf)

	if _, err := os.Stat(rcfPath); err != nil {
		slog.ErrorContext(ctx, "Checking file", slog.String("Uninstall", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	oldContent, err := os.ReadFile(rcfPath)
	if err != nil {
		slog.ErrorContext(ctx, "Reading file", slog.String("Uninstall", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	newContent := strings.ReplaceAll(string(oldContent), uninstall.Export(), "")

	if err := os.WriteFile(rcfPath, []byte(newContent), 0644); err != nil {
		slog.ErrorContext(ctx, "Writing file", slog.String("Uninstall", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	return nil
}
