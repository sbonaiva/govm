package api

import (
	"context"
	"sync"

	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

var (
	instance *cobra.Command
	once     sync.Once
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	once.Do(func() {
		if instance == nil {
			instance = &cobra.Command{
				Use:     "govm",
				Short:   "::: Go Version Manager :::",
				Version: util.GoVersionManager,
			}

			instance.AddCommand(
				NewListCmd(ctx, handler.NewList()),
				NewInstallCmd(ctx, handler.NewInstall()),
				NewUninstallCmd(ctx, handler.NewUninstall()),
			)
		}
	})
	return instance
}
