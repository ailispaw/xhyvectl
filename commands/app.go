package commands

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	COMMAND = "xhyvectl"
)

var (
	VERSION string
	GITSHA1 string
)

var app = &cobra.Command{
	Use:   COMMAND,
	Short: "A management tool for xhyve",
	Long:  COMMAND + " - A management tool for xhyve",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

func init() {
	cobra.OnInitialize(Initialize)
}

func Initialize() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
}

func Execute() {
	app.AddCommand(cmdVersion)
	app.AddCommand(cmdVmnet)

	app.SetOutput(os.Stdout)
	app.Execute()
}
