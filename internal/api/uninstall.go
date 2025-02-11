package api

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/service"
	"github.com/sbonaiva/govm/internal/util"
	"github.com/spf13/cobra"
)

func NewUninstallCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"u"},
		Short:   "Uninstall a Go version",
		Long:    "Uninstall a Go version",
		Example: "govm uninstall",
		Run: func(cmd *cobra.Command, args []string) {
			// Uninstall a Go version
			proceed := false
			reader := bufio.NewReader(os.Stdin)
			for {
				if !proceed {
					fmt.Print("Confirm uninstall current Go version? (y/n): ")
					confirmation, _ := reader.ReadString('\n')
					confirmation = strings.TrimSpace(confirmation)

					if confirmation == "y" {
						proceed = true
						break
					}

					if confirmation == "n" {
						util.PrintWarning("Uninstall aborted by user")
						return
					}

					util.PrintError("Invalid option, please type 'y' or 'n'")
					continue
				}
			}
			if err := service.NewUninstall().Execute(
				ctx,
				&domain.Uninstall{},
			); err != nil {
				util.PrintError(err.Error())
				return
			}
			util.PrintSuccess("Go uninstalled successfully!")
		},
	}
}
