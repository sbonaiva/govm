package handler

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/sbonaiva/govm/internal/gateway"
)

type ListHandler interface {
	Handle(ctx context.Context) error
}

type listHandler struct {
	httpGateway gateway.HttpGateway
}

func NewList() ListHandler {
	return &listHandler{
		httpGateway: gateway.NewHttpGateway(),
	}
}

func (r *listHandler) Handle(ctx context.Context) error {

	slog.InfoContext(ctx, "Listing all Go versions", slog.String("List", "Execute"))

	versions, err := r.httpGateway.GetVersions(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error while getting versions", slog.String("List", "Execute"), slog.String("error", err.Error()))
		return err
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("Available Go versions for %s/%s \n", runtime.GOOS, runtime.GOARCH)
	fmt.Println(strings.Repeat("=", 100))

	numCols := 6
	maxRows := (len(versions) + numCols - 1) / numCols

	for i := 0; i < maxRows; i++ {
		var row []string
		for j := 0; j < numCols; j++ {
			idx := i + j*maxRows
			if idx < len(versions) {
				row = append(row, fmt.Sprintf("%-15s", versions[idx].String()))
			}
		}
		fmt.Println(strings.Join(row, ""))
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("* currently in use")
	fmt.Println(strings.Repeat("=", 100))

	return nil
}
