package api

import (
	"context"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(ctx context.Context, handler handler.UpdateHandler) *cobra.Command {
	var updateStrategyParam domain.UpdateStrategy

	updateCmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"update"},
		Short:   "Update Go version",
		Long:    "Update Go version to latest major, minor or patch version",
		Example: "govm update [patch|minor|major]",
		Run: func(cmd *cobra.Command, args []string) {
			v, err := handler.Handle(ctx, &domain.Action{UpdateStrategy: updateStrategyParam})
			if err != nil {
				util.PrintError(err.Error())
				return
			}
			util.PrintSuccess("Go updated to version \"%s\" successfully!", v)
			util.PrintWarning("Please, reopen your terminal to start using new version.")
		},
	}

	updateCmd.Flags().StringVarP(
		(*string)(&updateStrategyParam),
		"strategy",
		"s",
		string(domain.PatchStrategy),
		"Update strategy to use (patch, minor, major)",
	)

	return updateCmd
}
