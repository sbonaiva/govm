package api

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewInstallCmd(ctx context.Context, handler handler.InstallHandler) *cobra.Command {
	return &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Install a Go version",
		Long:    "Install a Go version",
		Example: "govm install [version]",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.Handle(
				ctx,
				&domain.Install{
					Version: args[0],
				},
			); err != nil {
				util.PrintError(err.Error())
				return
			}
			util.PrintSuccess("Go version \"%s\" installed successfully!", args[0])
			util.PrintWarning("Please, reopen your terminal to start using new version.")
		},
	}
}
