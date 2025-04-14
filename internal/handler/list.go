package handler

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/gateway"
)

type ListHandler interface {
	Handle(ctx context.Context) error
}

type listHandler struct {
	httpGateway gateway.HttpGateway
}

func NewList(httpGateway gateway.HttpGateway) ListHandler {
	return &listHandler{
		httpGateway: httpGateway,
	}
}

func (r *listHandler) Handle(ctx context.Context) error {

	slog.InfoContext(ctx, "Listing all Go versions", slog.String("ListHandler", "Handle"))

	res, err := r.httpGateway.GetVersions(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error while getting versions", slog.String("ListHandler", "Handle"), slog.String("error", err.Error()))
		return domain.NewUnexpectedError(domain.ErrCodeListVersions)
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Printf("Available Go versions for %s/%s \n", runtime.GOOS, runtime.GOARCH)
	fmt.Println(strings.Repeat("=", 100))

	numCols := 6
	maxRows := (len(res.Versions) + numCols - 1) / numCols

	for i := 0; i < maxRows; i++ {
		var row []string
		for j := 0; j < numCols; j++ {
			idx := i + j*maxRows
			if idx < len(res.Versions) {
				row = append(row, fmt.Sprintf("%-15s", res.Versions[idx].String()))
			}
		}
		fmt.Println(strings.Join(row, ""))
	}

	fmt.Println(strings.Repeat("=", 100))
	fmt.Println("* currently in use")
	fmt.Println(strings.Repeat("=", 100))

	return nil
}
