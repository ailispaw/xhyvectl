package commands

import (
	"os"
	"os/exec"
	"os/user"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ErrorExit(ctx *cobra.Command, message string) {
	if message != "" {
		log.Error(message)
	}
	ctx.Usage()
	os.Exit(1)
}

func Sudo(ctx *cobra.Command) {
	user, err := user.Current()
	if err != nil {
		ErrorExit(ctx, "Could not get the current user")
	}
	if user.Uid == "0" {
		return
	}

	log.Warn("Requires root permissions")
	log.Warnf("Executes sudo %s", strings.Join(os.Args, " "))

	cmd := exec.Command("sudo", os.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
