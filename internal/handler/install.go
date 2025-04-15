package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/service"
)

type InstallHandler interface {
	Handle(ctx context.Context, install *domain.Action) error
}

type installHandler struct {
	sharedSvc service.SharedService
}

func NewInstall(sharedHandler service.SharedService) InstallHandler {
	return &installHandler{
		sharedSvc: sharedHandler,
	}
}

func (r *installHandler) Handle(ctx context.Context, install *domain.Action) error {
	slog.InfoContext(ctx, "Installing Go version", slog.String("InstallHandler", "Handle"), slog.String("version", install.Version))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	defer spn.Stop()
	spn.Start()

	steps := []struct {
		message string
		action  func() error
	}{
		{" Getting home...", func() error { return r.sharedSvc.CheckUserHome(ctx, install) }},
		{" Checking version...", func() error { return r.sharedSvc.CheckVersion(ctx, install) }},
		{" Downloading files...", func() error { return r.sharedSvc.DownloadVersion(ctx, install) }},
		{" Verifying checksum...", func() error { return r.sharedSvc.Checksum(ctx, install) }},
		{" Removing previous version...", func() error { return r.sharedSvc.RemoveVersion(ctx, install) }},
		{" Extracting files...", func() error { return r.sharedSvc.UntarFiles(ctx, install) }},
		{" Adding to path...", func() error { return r.sharedSvc.AddToPath(ctx, install) }},
	}

	for _, step := range steps {
		spn.Suffix = step.message
		if err := step.action(); err != nil {
			return err
		}
	}

	return nil
}
