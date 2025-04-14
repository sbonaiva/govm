package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/util"
)

var (
	shellRunCommandFiles = map[string]string{
		"/bin/bash":     ".bashrc",
		"/usr/bin/bash": ".bashrc",
		"/bin/zsh":      ".zshrc",
		"/usr/bin/zsh":  ".zshrc",
		"/bin/ksh":      ".kshrc",
		"/usr/bin/ksh":  ".kshrc",
		"/bin/fish":     ".config/fish/config.fish",
		"/usr/bin/fish": ".config/fish/config.fish",
	}
)

type SharedService interface {
	CheckUserHome(ctx context.Context, action *domain.Action) error
	CheckVersion(ctx context.Context, action *domain.Action) error
	DownloadVersion(ctx context.Context, action *domain.Action) error
	Checksum(ctx context.Context, action *domain.Action) error
	RemoveVersion(ctx context.Context, action *domain.Action) error
	UntarFiles(ctx context.Context, action *domain.Action) error
	AddToPath(ctx context.Context, action *domain.Action) error
	RemoveFromPath(ctx context.Context, action *domain.Action) error
	CheckInstalledVersion(ctx context.Context, action *domain.Action) error
	CheckAvailableUpdates(ctx context.Context, action *domain.Action) error
}

type sharedService struct {
	httpGateway gateway.HttpGateway
	osGateway   gateway.OsGateway
}

type version struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

func NewShared(httpGateway gateway.HttpGateway, osGateway gateway.OsGateway) SharedService {
	return &sharedService{
		httpGateway: httpGateway,
		osGateway:   osGateway,
	}
}

func (r *sharedService) CheckUserHome(ctx context.Context, action *domain.Action) error {
	homeDir, err := r.osGateway.GetUserHomeDir()
	if err != nil {
		slog.ErrorContext(ctx, "Getting current user", slog.String("SharedService", "CheckUserHome"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeCheckUserHome)
	}
	action.HomeDir = homeDir
	return nil
}

func (r *sharedService) CheckVersion(ctx context.Context, action *domain.Action) error {
	ok, err := r.httpGateway.VersionExists(ctx, action.Version)
	if err != nil {
		slog.ErrorContext(ctx, "Checking version", slog.String("SharedService", "CheckVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeCheckVersion)
	}

	if !ok {
		return domain.NewVersionNotAvailableError(action.Version)
	}
	return nil
}

func (r *sharedService) DownloadVersion(ctx context.Context, action *domain.Action) error {
	if err := r.osGateway.RemoveDir(action.DownloadFile()); err != nil {
		slog.ErrorContext(ctx, "Removing previous download", slog.String("SharedService", "DownloadVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeDownloadRemoveDir)
	}

	file, err := r.osGateway.CreateFile(action.DownloadFile())
	if err != nil {
		slog.ErrorContext(ctx, "Allocating resources", slog.String("SharedService", "DownloadVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeDownloadCreateFile)
	}
	defer file.Close()

	if err := r.httpGateway.DownloadVersion(ctx, action, file); err != nil {
		slog.ErrorContext(ctx, "Downloading version", slog.String("SharedService", "DownloadVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeDownloadVersion)
	}

	return nil
}

func (r *sharedService) Checksum(ctx context.Context, action *domain.Action) error {
	expectedChecksum, err := r.httpGateway.GetChecksum(ctx, action.Version)
	if err != nil {
		slog.ErrorContext(ctx, "Getting checksum", slog.String("SharedService", "Checksum"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeChecksumDownload)
	}

	file, err := r.osGateway.OpenFile(action.DownloadFile())
	if err != nil {
		slog.ErrorContext(ctx, "Opening file", slog.String("SharedService", "Checksum"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeChecksumOpenFile)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		slog.ErrorContext(ctx, "Calculating checksum", slog.String("SharedService", "Checksum"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeChecksumCopy)
	}

	if expectedChecksum != fmt.Sprintf("%x", hash.Sum(nil)) {
		slog.ErrorContext(ctx, "Checksum does not match", slog.String("SharedService", "Checksum"))
		return domain.NewUnexpectedError(domain.ErrCodeChecksumMismatch)
	}

	return nil
}

func (r *sharedService) RemoveVersion(ctx context.Context, action *domain.Action) error {
	if err := r.osGateway.RemoveDir(action.HomeGoDir()); err != nil {
		slog.ErrorContext(ctx, "Removing version", slog.String("SharedService", "RemoveVersion"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeRemoveVersion)
	}
	return nil
}

func (r *sharedService) UntarFiles(ctx context.Context, action *domain.Action) error {
	if err := r.osGateway.CreateDir(action.HomeGovmDir(), 0755); err != nil {
		slog.ErrorContext(ctx, "Creating directory", slog.String("SharedService", "UntarFiles"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUntarCreateDir)
	}

	if err := r.osGateway.Untar(action.DownloadFile(), action.HomeGovmDir()); err != nil {
		slog.ErrorContext(ctx, "Extracting files", slog.String("SharedService", "UntarFiles"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeUntarExtract)
	}

	defer r.osGateway.RemoveFile(action.DownloadFile())

	return nil
}

func (r *sharedService) AddToPath(ctx context.Context, action *domain.Action) error {
	if path := r.osGateway.GetEnv("PATH"); strings.Contains(path, action.HomeGoBinDir()) {
		slog.InfoContext(ctx, "Go is already in PATH", slog.String("SharedService", "AddToPath"))
		return nil
	}

	if shell := r.osGateway.GetEnv("SHELL"); shell != "" {
		if rcf, exists := shellRunCommandFiles[shell]; exists {
			return r.addToShellRunCommands(ctx, action, rcf)
		}
	}

	succeded := 0
	for _, rcf := range shellRunCommandFiles {
		err := r.addToShellRunCommands(ctx, action, rcf)
		if err == nil {
			succeded++
		}
	}

	if succeded == 0 {
		slog.ErrorContext(ctx, "No shell rc file found", slog.String("SharedService", "AddToPath"))
		return domain.NewUnexpectedError(domain.ErrCodeAddToPathNoShellsFound)
	}

	return nil
}

func (r *sharedService) RemoveFromPath(ctx context.Context, action *domain.Action) error {
	if path := r.osGateway.GetEnv("PATH"); !strings.Contains(path, action.HomeGoBinDir()) {
		slog.InfoContext(ctx, "Go is already removed from PATH", slog.String("SharedService", "RemoveFromPath"))
		return nil
	}

	if shell := r.osGateway.GetEnv("SHELL"); shell != "" {
		if rcf, exists := shellRunCommandFiles[shell]; exists {
			return r.removeFromShellRunCommands(ctx, action, rcf)
		}
	}

	succeded := 0
	for _, rcf := range shellRunCommandFiles {
		err := r.removeFromShellRunCommands(ctx, action, rcf)
		if err == nil {
			succeded++
		}
	}

	if succeded == 0 {
		slog.ErrorContext(ctx, "No shell rc file found", slog.String("SharedService", "RemoveFromPath"))
		return domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathNoShellsFound)
	}

	return nil
}

func (r *sharedService) addToShellRunCommands(ctx context.Context, action *domain.Action, rcf string) error {
	rcfPath := filepath.Join(action.HomeDir, rcf)

	if _, err := r.osGateway.Stat(rcfPath); err != nil {
		slog.ErrorContext(ctx, "Checking file", slog.String("SharedService", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeAddToPathStat)
	}

	oldContent, err := r.osGateway.ReadFile(rcfPath)
	if err != nil {
		slog.ErrorContext(ctx, "Reading file", slog.String("SharedService", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeAddToPathRead)
	}

	newContent := []byte(fmt.Sprintf("%s\n%s", string(oldContent), action.Export()))

	if err := r.osGateway.WriteFile(rcfPath, newContent, 0644); err != nil {
		slog.ErrorContext(ctx, "Writing file", slog.String("SharedService", "addToShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeAddToPathWrite)
	}

	return nil
}

func (r *sharedService) removeFromShellRunCommands(ctx context.Context, action *domain.Action, rcf string) error {
	rcfPath := filepath.Join(action.HomeDir, rcf)

	if _, err := r.osGateway.Stat(rcfPath); err != nil {
		slog.ErrorContext(ctx, "Checking file", slog.String("SharedService", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathStat)
	}

	oldContent, err := r.osGateway.ReadFile(rcfPath)
	if err != nil {
		slog.ErrorContext(ctx, "Reading file", slog.String("SharedService", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathRead)
	}

	newContent := strings.ReplaceAll(string(oldContent), action.Export(), "")

	if err := r.osGateway.WriteFile(rcfPath, []byte(newContent), 0644); err != nil {
		slog.ErrorContext(ctx, "Writing file", slog.String("SharedService", "removeFromShellRunCommands"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeRemoveFromPathWrite)
	}

	return nil
}

func (r *sharedService) CheckInstalledVersion(ctx context.Context, action *domain.Action) error {
	installedVersion, err := util.GetInstalledGoVersion()
	if err != nil {
		slog.ErrorContext(ctx, "Get installed version", slog.String("SharedService", "CheckInstalledVersion"), slog.String("error", err.Error()))
		return domain.NewNoGoInstallationsFoundError()
	}
	action.InstalledVersion = installedVersion
	return nil
}

func (r *sharedService) CheckAvailableUpdates(ctx context.Context, action *domain.Action) error {

	availableVersions, err := r.httpGateway.GetVersions(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Get versions", slog.String("SharedService", "CheckAvailableUpdates"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeListVersions)
	}

	patch, minor, major := r.findLatestVersions(action.InstalledVersion, availableVersions.StringSlice())

	switch action.UpdateStrategy {
	case domain.MajorStrategy:
		if major.Raw == action.InstalledVersion {
			return domain.NewNoUpdatesAvailableError(domain.MajorStrategy, action.InstalledVersion)
		}
		action.Version = major.Raw
	case domain.MinorStrategy:
		if minor.Raw == action.InstalledVersion {
			return domain.NewNoUpdatesAvailableError(domain.MinorStrategy, action.InstalledVersion)
		}
		action.Version = minor.Raw
	case domain.PatchStrategy:
		if patch.Raw == action.InstalledVersion {
			return domain.NewNoUpdatesAvailableError(domain.PatchStrategy, action.InstalledVersion)
		}
		action.Version = patch.Raw
	}

	return nil
}

func (r *sharedService) findLatestVersions(cv string, av []string) (patch, minor, major *version) {
	currentVersion := r.parseVersion(cv)

	availableVersions := make([]version, len(av))
	for i, v := range av {
		availableVersions[i] = r.parseVersion(v)
	}

	sort.Slice(availableVersions, func(i, j int) bool {
		return r.compareVersionDesc(availableVersions[i], availableVersions[j])
	})

	for _, v := range availableVersions {
		if patch == nil && v.Major == currentVersion.Major && v.Minor == currentVersion.Minor {
			patch = &v
		}

		if minor == nil && v.Major == currentVersion.Major {
			minor = &v
		}

		if major == nil {
			major = &v
		}
		if patch != nil && minor != nil && major != nil {
			break
		}
	}
	return
}

func (r *sharedService) compareVersionDesc(v1, v2 version) bool {
	if v1.Major != v2.Major {
		return v1.Major > v2.Major
	}

	if v1.Minor != v2.Minor {
		return v1.Minor > v2.Minor
	}

	return v1.Patch > v2.Patch
}

func (r *sharedService) parseVersion(s string) version {
	parts := strings.Split(s, ".")
	v := version{Raw: s}

	if len(parts) > 0 {
		v.Major, _ = strconv.Atoi(parts[0])
	}
	if len(parts) > 1 {
		v.Minor, _ = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		v.Patch, _ = strconv.Atoi(parts[2])
	}
	return v
}
