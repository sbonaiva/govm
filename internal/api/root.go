package api

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/service"
	"github.com/spf13/cobra"
)

var (
	instance *cobra.Command
	once     sync.Once
)

func NewRootCmd(ctx context.Context, httpGateway gateway.HttpGateway, osGateway gateway.OsGateway) *cobra.Command {
	once.Do(func() {
		if instance == nil {
			instance = &cobra.Command{
				Use:     "govm",
				Short:   "::: Go Version Manager :::",
				Version: fmt.Sprintf("%s %s/%s", "0.0.2", runtime.GOOS, runtime.GOARCH),
			}

			sharedSvc := service.NewShared(httpGateway, osGateway)

			instance.AddCommand(
				NewListCmd(ctx, handler.NewList(sharedSvc)),
				NewInstallCmd(ctx, handler.NewInstall(sharedSvc)),
				NewUninstallCmd(ctx, handler.NewUninstall(sharedSvc)),
				NewUpdateCmd(ctx, handler.NewUpdate(sharedSvc)),
			)
		}
	})
	return instance
}
