package commands

import (
	"fmt"
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

	out, err := exec.Command("sudo", os.Args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", string(out))

	os.Exit(0)
}
