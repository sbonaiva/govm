package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/spf13/cobra"
)

func main() {

	os.Remove("govm.log")

	logFile, err := os.OpenFile("govm.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, nil)))

	ctx := context.Background()

	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(
		api.NewListCmd(ctx),
		api.NewInstallCmd(ctx),
		api.NewUninstallCmd(ctx),
	)

	if err := rootCmd.Execute(); err != nil {
		err.Error()
	}
}
