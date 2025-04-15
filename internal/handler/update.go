package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/service"
)

type UpdateHandler interface {
	Handle(ctx context.Context, update *domain.Action) (string, error)
}

type updateHandler struct {
	sharedSvc service.SharedService
}

func NewUpdate(sharedSvc service.SharedService) UpdateHandler {
	return &updateHandler{
		sharedSvc: sharedSvc,
	}
}

func (r *updateHandler) Handle(ctx context.Context, update *domain.Action) (string, error) {
	slog.InfoContext(ctx, "Updating Go version", slog.String("UpdateHandler", "Handle"))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	defer spn.Stop()
	spn.Start()

	steps := []struct {
		message string
		action  func() error
	}{
		{" Checking update strategy...", func() error { return update.CheckUpdateStrategy() }},
		{" Checking installed version...", func() error { return r.sharedSvc.CheckInstalledVersion(ctx, update) }},
		{" Checking available updates...", func() error { return r.sharedSvc.CheckAvailableUpdates(ctx, update) }},
		{" Getting home...", func() error { return r.sharedSvc.CheckUserHome(ctx, update) }},
		{" Checking version...", func() error { return r.sharedSvc.CheckVersion(ctx, update) }},
		{" Downloading files...", func() error { return r.sharedSvc.DownloadVersion(ctx, update) }},
		{" Verifying checksum...", func() error { return r.sharedSvc.Checksum(ctx, update) }},
		{" Removing previous version...", func() error { return r.sharedSvc.RemoveVersion(ctx, update) }},
		{" Extracting files...", func() error { return r.sharedSvc.UntarFiles(ctx, update) }},
		{" Adding to path...", func() error { return r.sharedSvc.AddToPath(ctx, update) }},
	}

	for _, step := range steps {
		spn.Suffix = step.message
		if err := step.action(); err != nil {
			return "", err
		}
	}

	return update.Version, nil
}
