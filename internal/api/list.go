package api

import (
	"context"
	"sync"

	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

var (
	listHandler     handler.ListHandler
	onceListHandler sync.Once
)

func getListHandler() handler.ListHandler {
	onceListHandler.Do(func() {
		listHandler = handler.NewList()
	})
	return listHandler
}

func NewListCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all Go versions",
		Long:    "List all Go versions",
		Example: "govm list",
		Run: func(cmd *cobra.Command, args []string) {
			if err := getListHandler().Handle(ctx); err != nil {
				util.PrintError(err.Error())
			}
		},
	}
}
