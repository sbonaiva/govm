package api

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewLogCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "log",
		Short:   "Show log info",
		Long:    "Show log info",
		Example: "govm log",
		Run: func(cmd *cobra.Command, args []string) {
			logFilePath := path.Join(os.TempDir(), "govm.log")
			if _, err := os.Stat(logFilePath); err == nil {
				fmt.Println(strings.Repeat("=", 100))
				fmt.Println("Log file path:", logFilePath)
				fmt.Println(strings.Repeat("=", 100))

				logFile, err := os.Open(logFilePath)
				if err != nil {
					util.PrintError("Failed to open log file")
					os.Exit(1)
				}

				scanner := bufio.NewScanner(logFile)
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}

				return
			}

			util.PrintWarning("No log entries available.")
		},
	}
}
