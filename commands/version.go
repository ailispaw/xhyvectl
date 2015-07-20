package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show the version information",
	Long:    COMMAND + " version - Show the version information",
	Run:     showVersion,
}

func showVersion(ctx *cobra.Command, args []string) {
	fmt.Printf("%s: version %s, build %s\n", COMMAND, VERSION, GITSHA1)
}
