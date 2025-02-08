package api

import (
	"context"

	"github.com/sbonaiva/govm/internal/service"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewUninstallCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"u"},
		Short:   "Uninstall a Go version",
		Long:    "Uninstall a Go version",
		Example: "govm uninstall [version]",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.NewUninstall().Execute(ctx); err != nil {
				util.PrintError(err.Error())
				return
			}
			util.PrintSuccess("Go uninstalled successfully", args[0])
		},
	}
}
