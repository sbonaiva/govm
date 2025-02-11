package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/integration"
)

type Install interface {
	Execute(ctx context.Context, install *domain.Install) error
}

type install struct {
	goDevClient integration.GoDevClient
}

func NewInstall() Install {
	return &install{
		goDevClient: integration.NewGoDevClient(),
	}
}

func (r *install) Execute(ctx context.Context, install *domain.Install) error {
	slog.InfoContext(ctx, "Installing Go version", slog.String("Install", "Execute"), slog.String("version", install.Version))

	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
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

func (r *install) checkUserHome(ctx context.Context, install *domain.Install) error {
	usr, err := user.Current()
	if err != nil {
		slog.ErrorContext(ctx, "Getting current user", slog.String("Install", "home"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	install.HomeDir = usr.HomeDir
	return nil
}

func (r *install) checkVersion(ctx context.Context, install *domain.Install) error {
	ok, err := r.goDevClient.VersionExists(ctx, install.Version)
	if err != nil {
		slog.ErrorContext(ctx, "Checking version", slog.String("Install", "checkVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	if !ok {
		return fmt.Errorf("Go version \"%s\" is not available", install.Version)
	}
	return nil
}

func (r *install) downloadVersion(ctx context.Context, install *domain.Install) error {
	if err := os.Remove(install.DownloadDir()); err != nil && !os.IsNotExist(err) {
		slog.ErrorContext(ctx, "Removing previous download", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	file, err := os.Create(install.DownloadDir())
	if err != nil {
		slog.ErrorContext(ctx, "Allocating resources", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	defer file.Close()

	if err := r.goDevClient.DownloadVersion(ctx, *install, file); err != nil {
		slog.ErrorContext(ctx, "Downloading version", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *install) checksum(ctx context.Context, install *domain.Install) error {
	expectedChecksum, err := r.goDevClient.GetChecksum(ctx, install.Version)
	if err != nil {
		slog.ErrorContext(ctx, "Getting checksum", slog.String("Install", "checksum"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	file, err := os.Open(install.DownloadDir())
	if err != nil {
		slog.ErrorContext(ctx, "Opening file", slog.String("Install", "checksum"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		slog.ErrorContext(ctx, "Calculating checksum", slog.String("Install", "checksum"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	if expectedChecksum != fmt.Sprintf("%x", hash.Sum(nil)) {
		slog.ErrorContext(ctx, "Checksum does not match", slog.String("Install", "checksum"))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *install) removePreviousVersion(ctx context.Context, install *domain.Install) error {
	if err := os.RemoveAll(install.HomeGoDir()); err != nil && !os.IsNotExist(err) {
		slog.ErrorContext(ctx, "Removing previous version", slog.String("Install", "removePreviousVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	return nil
}

func (r *install) untarFiles(ctx context.Context, install *domain.Install) error {
	if err := os.Mkdir(install.HomeGovmDir(), 0755); err != nil && !os.IsExist(err) {
		slog.ErrorContext(ctx, "Creating directory", slog.String("Install", "untar"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	cmd := exec.Command("tar", "-C", install.HomeGovmDir(), "-xzf", install.DownloadDir())
	if err := cmd.Run(); err != nil {
		slog.ErrorContext(ctx, "Extracting files", slog.String("Install", "untar"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	defer os.Remove(install.DownloadDir())

	return nil
}

func (r *install) addToPath(ctx context.Context, install *domain.Install) error {
	if path := os.Getenv("PATH"); strings.Contains(path, install.HomeGoBinDir()) {
		slog.InfoContext(ctx, "Go is already in PATH", slog.String("Install", "addToPath"))
		return nil
	}

	if shell := os.Getenv("SHELL"); shell != "" {
		if rcf, exists := domain.ShellRunCommandsFiles[shell]; exists {
			return r.addToShellRunCommands(ctx, install, rcf, shell)
		}
	}

	succeded := 0
	for shell, rcf := range domain.ShellRunCommandsFiles {
		err := r.addToShellRunCommands(ctx, install, rcf, shell)
		if err == nil {
			succeded++
		}
	}

	if succeded == 0 {
		slog.ErrorContext(ctx, "No shell rc file found", slog.String("Install", "addToPath"))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *install) addToShellRunCommands(ctx context.Context, install *domain.Install, rcf string, shell string) error {
	rcfPath := filepath.Join(install.HomeDir, rcf)

	if _, err := os.Stat(rcfPath); err != nil {
		slog.ErrorContext(ctx, "Checking file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	oldContent, err := os.ReadFile(rcfPath)
	if err != nil {
		slog.ErrorContext(ctx, "Reading file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	newContent := []byte(fmt.Sprintf("%s\n%s", string(oldContent), install.Export()))

	if err := os.WriteFile(rcfPath, newContent, 0644); err != nil {
		slog.ErrorContext(ctx, "Writing file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	return nil
}
