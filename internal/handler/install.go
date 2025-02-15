package handler

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
)

type InstallHandler interface {
	Handle(ctx context.Context, install *domain.Install) error
}

type installHandler struct {
	httpGateway gateway.HttpGateway
	osGateway   gateway.OsGateway
}

func NewInstall(httpGateway gateway.HttpGateway, osGateway gateway.OsGateway) InstallHandler {
	return &installHandler{
		httpGateway: httpGateway,
		osGateway:   osGateway,
	}
}

func (r *installHandler) Handle(ctx context.Context, install *domain.Install) error {
	slog.InfoContext(ctx, "Installing Go version", slog.String("Install", "Execute"), slog.String("version", install.Version))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	defer spn.Stop()
	spn.Start()

	steps := []struct {
		message string
		action  func() error
	}{
		{" Getting home...", func() error { return r.checkUserHome(ctx, install) }},
		{" Checking version...", func() error { return r.checkVersion(ctx, install) }},
		{" Downloading files...", func() error { return r.downloadVersion(ctx, install) }},
		{" Verifying checksum...", func() error { return r.checksum(ctx, install) }},
		{" Removing previous version...", func() error { return r.removePreviousVersion(ctx, install) }},
		{" Extracting files...", func() error { return r.untarFiles(ctx, install) }},
		{" Adding to path...", func() error { return r.addToPath(ctx, install) }},
	}

	for _, step := range steps {
		spn.Suffix = step.message
		if err := step.action(); err != nil {
			return err
		}
	}

	return nil
}

func (r *installHandler) checkUserHome(ctx context.Context, install *domain.Install) error {
	homeDir, err := r.osGateway.GetUserHomeDir()
	if err != nil {
		slog.ErrorContext(ctx, "Getting current user", slog.String("Install", "home"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallCheckUserHome)
	}
	install.HomeDir = homeDir
	return nil
}

func (r *installHandler) checkVersion(ctx context.Context, install *domain.Install) error {
	ok, err := r.httpGateway.VersionExists(ctx, install.Version)
	if err != nil {
		slog.ErrorContext(ctx, "Checking version", slog.String("Install", "checkVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallCheckVersion)
	}

	if !ok {
		return domain.NewVersionNotAvailableError(install.Version)
	}
	return nil
}

func (r *installHandler) downloadVersion(ctx context.Context, install *domain.Install) error {
	if err := r.osGateway.RemoveDir(install.DownloadFile()); err != nil {
		slog.ErrorContext(ctx, "Removing previous download", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallDownloadRemoveDir)
	}

	file, err := r.osGateway.CreateFile(install.DownloadFile())
	if err != nil {
		slog.ErrorContext(ctx, "Allocating resources", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallDownloadCreateFile)
	}
	defer file.Close()

	if err := r.httpGateway.DownloadVersion(ctx, *install, file); err != nil {
		slog.ErrorContext(ctx, "Downloading version", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallDownloadVersion)
	}

	return nil
}

func (r *installHandler) checksum(ctx context.Context, install *domain.Install) error {
	expectedChecksum, err := r.httpGateway.GetChecksum(ctx, install.Version)
	if err != nil {
		slog.ErrorContext(ctx, "Getting checksum", slog.String("Install", "checksum"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallChecksumDownload)
	}

	file, err := r.osGateway.OpenFile(install.DownloadFile())
	if err != nil {
		slog.ErrorContext(ctx, "Opening file", slog.String("Install", "checksum"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallChecksumOpenFile)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		slog.ErrorContext(ctx, "Calculating checksum", slog.String("Install", "checksum"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallChecksumCopy)
	}

	if expectedChecksum != fmt.Sprintf("%x", hash.Sum(nil)) {
		slog.ErrorContext(ctx, "Checksum does not match", slog.String("Install", "checksum"))
		return domain.NewUnexpectedError(domain.ErrCodeInstallChecksumMismatch)
	}

	return nil
}

func (r *installHandler) removePreviousVersion(ctx context.Context, install *domain.Install) error {
	if err := r.osGateway.RemoveDir(install.HomeGoDir()); err != nil {
		slog.ErrorContext(ctx, "Removing previous version", slog.String("Install", "removePreviousVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallRemovePreviousVersion)
	}
	return nil
}

func (r *installHandler) untarFiles(ctx context.Context, install *domain.Install) error {
	if err := r.osGateway.CreateDir(install.HomeGovmDir(), 0755); err != nil {
		slog.ErrorContext(ctx, "Creating directory", slog.String("Install", "untar"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallUntarCreateDir)
	}

	if err := r.osGateway.Untar(install.DownloadFile(), install.HomeGovmDir()); err != nil {
		slog.ErrorContext(ctx, "Extracting files", slog.String("Install", "untar"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallUntarExtract)
	}

	defer r.osGateway.RemoveFile(install.DownloadFile())

	return nil
}

func (r *installHandler) addToPath(ctx context.Context, install *domain.Install) error {
	if path := r.osGateway.GetEnv("PATH"); strings.Contains(path, install.HomeGoBinDir()) {
		slog.InfoContext(ctx, "Go is already in PATH", slog.String("Install", "addToPath"))
		return nil
	}

	if shell := r.osGateway.GetEnv("SHELL"); shell != "" {
		if rcf, exists := domain.ShellRunCommandsFiles[shell]; exists {
			return r.addToShellRunCommands(ctx, install, rcf)
		}
	}

	succeded := 0
	for _, rcf := range domain.ShellRunCommandsFiles {
		err := r.addToShellRunCommands(ctx, install, rcf)
		if err == nil {
			succeded++
		}
	}

	if succeded == 0 {
		slog.ErrorContext(ctx, "No shell rc file found", slog.String("Install", "addToPath"))
		return domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathNoShellsFound)
	}

	return nil
}

func (r *installHandler) addToShellRunCommands(ctx context.Context, install *domain.Install, rcf string) error {
	rcfPath := filepath.Join(install.HomeDir, rcf)

	if _, err := r.osGateway.Stat(rcfPath); err != nil {
		slog.ErrorContext(ctx, "Checking file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathStat)
	}

	oldContent, err := r.osGateway.ReadFile(rcfPath)
	if err != nil {
		slog.ErrorContext(ctx, "Reading file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathRead)
	}

	newContent := []byte(fmt.Sprintf("%s\n%s", string(oldContent), install.Export()))

	if err := r.osGateway.WriteFile(rcfPath, newContent, 0644); err != nil {
		slog.ErrorContext(ctx, "Writing file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeInstallAddToPathWrite)
	}

	return nil
}
