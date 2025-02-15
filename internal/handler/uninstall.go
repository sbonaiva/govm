package handler

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
)

type UninstallHandler interface {
	Handle(ctx context.Context, uninstall *domain.Uninstall) error
}

type uninstallHandler struct {
	osGateway gateway.OsGateway
}

func NewUninstall(osGateway gateway.OsGateway) UninstallHandler {
	return &uninstallHandler{
		osGateway: osGateway,
	}
}

func (r *uninstallHandler) Handle(ctx context.Context, uninstall *domain.Uninstall) error {
	slog.InfoContext(ctx, "Uninstalling Go version", slog.String("Uninstall", "Execute"))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
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
	usrHomeDir, err := r.osGateway.GetUserHomeDir()
	if err != nil {
		slog.ErrorContext(ctx, "Getting current user", slog.String("Uninstall", "checkUserHome"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallCheckUserHome)
	}
	uninstall.HomeDir = usrHomeDir
	return nil
}

func (r *uninstallHandler) checkVersion(ctx context.Context, uninstall *domain.Uninstall) error {

	fs, err := r.osGateway.Stat(uninstall.HomeGoDir())
	if err != nil {
		slog.ErrorContext(ctx, "Checking version", slog.String("Uninstall", "checkVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallCheckVersionStat)
	}

	if !fs.IsDir() {
		slog.ErrorContext(ctx, "Checking version", slog.String("Uninstall", "checkVersion"), slog.String("error", "Go directory not found"))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallCheckVersionNotDir)
	}

	return nil
}

func (r *uninstallHandler) removeCurrentVersion(ctx context.Context, uninstall *domain.Uninstall) error {
	if err := r.osGateway.RemoveDir(uninstall.HomeGoDir()); err != nil {
		slog.ErrorContext(ctx, "Removing current version", slog.String("Uninstall", "removeCurrentVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveDir)
	}
	return nil
}

func (r *uninstallHandler) removeFromPath(ctx context.Context, uninstall *domain.Uninstall) error {
	if path := r.osGateway.GetEnv("PATH"); !strings.Contains(path, uninstall.HomeGoBinDir()) {
		slog.InfoContext(ctx, "Go is already removed from PATH", slog.String("Uninstall", "removeFromPath"))
		return nil
	}

	if shell := r.osGateway.GetEnv("SHELL"); shell != "" {
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
		return domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathNoShellsFound)
	}

	return nil
}

func (r *uninstallHandler) removeFromShellRunCommands(ctx context.Context, uninstall *domain.Uninstall, rcf string) error {
	rcfPath := filepath.Join(uninstall.HomeDir, rcf)

	if _, err := r.osGateway.Stat(rcfPath); err != nil {
		slog.ErrorContext(ctx, "Checking file", slog.String("Uninstall", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathStat)
	}

	oldContent, err := r.osGateway.ReadFile(rcfPath)
	if err != nil {
		slog.ErrorContext(ctx, "Reading file", slog.String("Uninstall", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathRead)
	}

	newContent := strings.ReplaceAll(string(oldContent), uninstall.Export(), "")

	if err := r.osGateway.WriteFile(rcfPath, []byte(newContent), 0644); err != nil {
		slog.ErrorContext(ctx, "Writing file", slog.String("Uninstall", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUninstallRemoveFromPathWrite)
	}

	return nil
}
