package api

import (
	"context"

	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewListCmd(ctx context.Context, handler handler.ListHandler) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all Go versions",
		Long:    "List all Go versions",
		Example: "govm list",
		Run: func(cmd *cobra.Command, args []string) {
			if err := handler.Handle(ctx); err != nil {
				util.PrintError(err.Error())
			}
		},
	}
}
