package handler

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/sbonaiva/govm/internal/service"
)

type ListHandler interface {
	Handle(ctx context.Context) error
}

type listHandler struct {
	sharedSvc service.SharedService
}

func NewList(sharedSvc service.SharedService) ListHandler {
	return &listHandler{
		sharedSvc: sharedSvc,
	}
}

func (r *listHandler) Handle(ctx context.Context) error {

	slog.InfoContext(ctx, "Listing all Go versions", slog.String("ListHandler", "Handle"))

	availableVersions, err := r.sharedSvc.GetAvailableGoVersions(ctx)
	if err != nil {
		return err
	}

	installedVersion, _ := r.sharedSvc.GetInstalledGoVersion(ctx)

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("Available Go versions for %s/%s \n", runtime.GOOS, runtime.GOARCH)
	fmt.Println(strings.Repeat("=", 100))

	numCols := 6
	maxRows := (len(availableVersions.Versions) + numCols - 1) / numCols

	for i := 0; i < maxRows; i++ {
		var row []string
		for j := 0; j < numCols; j++ {
			idx := i + j*maxRows
			if idx < len(availableVersions.Versions) {
				row = append(row, fmt.Sprintf("%-15s", availableVersions.Versions[idx].String(installedVersion)))
			}
		}
		fmt.Println(strings.Join(row, ""))
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("* currently in use")
	fmt.Println(strings.Repeat("=", 100))

	return nil
}
