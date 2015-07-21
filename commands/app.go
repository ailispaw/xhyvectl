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

	boolVerbose, boolDebug bool
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
	flags := app.PersistentFlags()
	flags.BoolVarP(&boolVerbose, "verbose", "V", false, "Print verbose messages")
	flags.BoolVarP(&boolDebug, "debug", "D", false, "Print debug messages")

	cobra.OnInitialize(Initialize)
}

func Initialize() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	if boolVerbose {
		log.SetLevel(log.InfoLevel)
	}
	if boolDebug {
		log.SetLevel(log.DebugLevel)
	}
}

func Execute() {
	app.AddCommand(cmdPackage)
	app.AddCommand(cmdVersion)
	app.AddCommand(cmdVmnet)

	app.SetOutput(os.Stdout)
	app.Execute()
}
