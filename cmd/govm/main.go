package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/sbonaiva/govm/internal/util"
)

const (
	logCmd  = "log"
	logFile = "govm.log"
	goVersionsURL = "https://go.dev/dl/?mode=json&include=all"
	goDownloadURL = "https://go.dev/dl/%s"
)

var (
	Version = "dev"
)

func main() {
	ctx := context.Background()

	if !(len(os.Args) > 1 && os.Args[1] == logCmd) {
		logFilePath := path.Join(os.TempDir(), logFile)
		if err := os.Remove(logFilePath); err != nil && !os.IsNotExist(err) {
			util.PrintError("Failed to remove log file")
			os.Exit(1)
		}

		logFile, err := os.Create(logFilePath)
		if err != nil {
			util.PrintError("Failed to create log file")
			fmt.Println(err)
			os.Exit(1)
		}
		defer logFile.Close()

		slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, nil)))
	}

	httpGateway := gateway.NewHttpGateway(&gateway.HttpConfig{
		GoVersionURL:  goVersionsURL,
		GoDownloadURL: goDownloadURL,
	})
	osGateway := gateway.NewOsGateway()
	rootCmd := api.NewRootCmd(ctx, Version, httpGateway, osGateway)

	if err := rootCmd.Execute(); err != nil {
		util.PrintError(err.Error())
		os.Exit(1)
	}
}
