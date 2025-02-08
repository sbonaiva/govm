package api

import (
	"context"

	"github.com/sbonaiva/govm/internal/service"
	"github.com/spf13/cobra"
)

func NewListCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all Go versions",
		Long:    "List all Go versions",
		Example: "govm list",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.NewList().Execute(ctx); err != nil {
				err.Error()
			}
		},
	}
}
