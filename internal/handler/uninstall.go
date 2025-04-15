package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/service"
)

type UninstallHandler interface {
	Handle(ctx context.Context, uninstall *domain.Action) error
}

type uninstallHandler struct {
	sharedSvc service.SharedService
}

func NewUninstall(sharedHandler service.SharedService) UninstallHandler {
	return &uninstallHandler{
		sharedSvc: sharedHandler,
	}
}

func (r *uninstallHandler) Handle(ctx context.Context, uninstall *domain.Action) error {
	slog.InfoContext(ctx, "Uninstalling Go version", slog.String("UninstallHandler", "Handle"))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	defer spn.Stop()
	spn.Start()

	steps := []struct {
		message string
		action  func() error
	}{
		{" Getting home...", func() error { return r.sharedSvc.CheckUserHome(ctx, uninstall) }},
		{" Checking if Go is installed...", func() error { return r.checkIfGoIsInstalled(ctx) }},
		{" Removing current version...", func() error { return r.sharedSvc.RemoveVersion(ctx, uninstall) }},
		{" Removing from path...", func() error { return r.sharedSvc.RemoveFromPath(ctx, uninstall) }},
	}

	for _, step := range steps {
		spn.Suffix = step.message
		if err := step.action(); err != nil {
			return err
		}
	}

	return nil
}

func (r *uninstallHandler) checkIfGoIsInstalled(ctx context.Context) error {
	slog.InfoContext(ctx, "Checking uninstall", slog.String("UninstallHandler", "checkIfGoInstalled"))

	v, err := r.sharedSvc.GetInstalledGoVersion(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting installed Go version", slog.String("UninstallHandler", "checkIfGoInstalled"), slog.String("error", err.Error()))
		return err
	}

	if v == "" {
		slog.ErrorContext(ctx, "No Go installations found", slog.String("UninstallHandler", "checkIfGoInstalled"))
		return domain.NewNoGoInstallationsFoundError()
	}
	return nil
}
