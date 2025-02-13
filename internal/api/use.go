package api

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewUseCmd(ctx context.Context, handler handler.UseHandler) *cobra.Command {
	return &cobra.Command{
		Use:     "use",
		Aliases: []string{"use"},
		Short:   "Use a Go version",
		Long:    "Use a Go version",
		Example: "govm use [version]",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.Handle(
				ctx,
				&domain.Use{
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
