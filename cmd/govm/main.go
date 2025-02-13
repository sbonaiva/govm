package main

import (
	"context"
	"log/slog"
	"os"
	"os/user"
	"path"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/util"
)

const (
	logDir  = ".govm"
	logFile = "govm.log"
)

func main() {

	user, err := user.Current()
	if err != nil {
		util.PrintError("Failed to get current user")
		os.Exit(1)
	}

	logPath := path.Join(user.HomeDir, logDir, logFile)

	if err := os.Remove(logPath); err != nil && !os.IsNotExist(err) {
		util.PrintError("Failed to remove log file")
		os.Exit(1)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		util.PrintError("Failed to create log file")
		os.Exit(1)
	}
	defer logFile.Close()

	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, nil)))

	ctx := context.Background()

	if err := api.NewRootCmd(ctx).Execute(); err != nil {
		util.PrintError(err.Error())
		os.Exit(1)
	}
}
