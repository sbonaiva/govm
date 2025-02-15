package api

import (
	"context"
	"sync"

	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

var (
	instance *cobra.Command
	once     sync.Once
)

func NewRootCmd(
	ctx context.Context,
	httpGateway gateway.HttpGateway,
	osGateway gateway.OsGateway,
) *cobra.Command {
	once.Do(func() {
		if instance == nil {
			instance = &cobra.Command{
				Use:     "govm",
				Short:   "::: Go Version Manager :::",
				Version: util.GoVersionManager,
			}

			instance.AddCommand(
				NewListCmd(ctx, handler.NewList(httpGateway)),
				NewInstallCmd(ctx, handler.NewInstall(httpGateway, osGateway)),
				NewUninstallCmd(ctx, handler.NewUninstall(osGateway)),
				NewUseCmd(ctx, handler.NewUse(httpGateway)),
			)
		}
	})
	return instance
}
