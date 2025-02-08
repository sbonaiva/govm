package service

import (
	"context"
	"fmt"
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
	spn := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	defer spn.Stop()
	spn.Start()

	steps := []struct {
		message string
		action  func() error
	}{
		{" Getting home...", func() error { return r.home(install) }},
		{" Checking version...", func() error { return r.check(ctx, install) }},
		{" Downloading files...", func() error { return r.download(ctx, install) }},
		{" Removing previous version...", func() error { return r.remove(install) }},
		{" Extracting files...", func() error { return r.untar(install) }},
		{" Adding to path...", func() error { return r.path(install) }},
	}

	for _, step := range steps {
		spn.Suffix = step.message
		if err := step.action(); err != nil {
			return err
		}
	}

	return nil
}

func (r *install) home(install *domain.Install) error {
	usr, err := user.Current()
	if err != nil {
		slog.Error("Getting current user", slog.String("Install", "home"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	install.HomeDir = usr.HomeDir
	return nil
}

func (r *install) check(ctx context.Context, install *domain.Install) error {
	ok, err := r.goDevClient.VersionExists(ctx, install.Version)
	if err != nil {
		slog.Error("Checking version", slog.String("Install", "checkVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	if !ok {
		return fmt.Errorf("Go version \"%s\" is not available", install.Version)
	}
	return nil
}

func (r *install) download(ctx context.Context, install *domain.Install) error {
	if err := os.Remove(install.DownloadDir()); err != nil && !os.IsNotExist(err) {
		slog.Error("Removing previous download", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	file, err := os.Create(install.DownloadDir())
	if err != nil {
		slog.Error("Allocating resources", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	defer file.Close()

	if err := r.goDevClient.DownloadVersion(ctx, *install, file); err != nil {
		slog.Error("Downloading version", slog.String("Install", "downloadVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *install) remove(install *domain.Install) error {
	if err := os.RemoveAll(install.HomeGovmDir()); err != nil && !os.IsNotExist(err) {
		slog.Error("Removing previous version", slog.String("Install", "removePreviousVersion"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}
	return nil
}

func (r *install) untar(install *domain.Install) error {
	if err := os.Mkdir(install.HomeGovmDir(), 0755); err != nil && !os.IsExist(err) {
		slog.Error("Creating directory", slog.String("Install", "untar"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	cmd := exec.Command("tar", "-C", install.HomeGovmDir(), "-xzf", install.DownloadDir())
	if err := cmd.Run(); err != nil {
		slog.Error("Extracting files", slog.String("Install", "untar"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	defer os.Remove(install.DownloadDir())

	return nil
}

func (r *install) path(install *domain.Install) error {
	if path := os.Getenv("PATH"); strings.Contains(path, install.HomeGoBinDir()) {
		slog.Info("Go is already in PATH", slog.String("Install", "addToPath"))
		return nil
	}

	if shell := os.Getenv("SHELL"); shell != "" {
		if rcf, exists := domain.ShellRunCommandsFiles[shell]; exists {
			return r.addToShellRunCommands(install, rcf, shell)
		}
	}

	succeded := 0
	for shell, rcf := range domain.ShellRunCommandsFiles {
		err := r.addToShellRunCommands(install, rcf, shell)
		if err == nil {
			succeded++
		}
	}

	if succeded == 0 {
		slog.Error("No shell rc file found", slog.String("Install", "addToPath"))
		return domain.ErrUnexpected
	}

	return nil
}

func (r *install) addToShellRunCommands(install *domain.Install, rcf string, shell string) error {
	rcfPath := filepath.Join(install.HomeDir, rcf)

	if _, err := os.Stat(rcfPath); err != nil {
		slog.Error("Checking file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	oldContent, err := os.ReadFile(rcfPath)
	if err != nil {
		slog.Error("Reading file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	newContent := []byte(fmt.Sprintf("%s\n%s", string(oldContent), install.Export()))

	if err := os.WriteFile(rcfPath, newContent, 0644); err != nil {
		slog.Error("Writing file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	cmd := exec.Command(shell, "-c", fmt.Sprintf("source %s", rcfPath))
	if err := cmd.Run(); err != nil {
		slog.Error("Sourcing file", slog.String("Install", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.ErrUnexpected
	}

	return nil
}
