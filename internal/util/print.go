package util

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintSuccess(message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	fmt.Println(color.GreenString(message))
}

func PrintError(message string) {
	fmt.Println(color.RedString(message))
}
