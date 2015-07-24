package commands

import (
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
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
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func SetupProgressBar(size int, prefix string) *pb.ProgressBar {
	bar := pb.New(size)
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond * 10)
	bar.ShowCounters = false
	bar.ShowFinalTime = false
	bar.SetWidth(100)
	bar.Prefix(prefix)
	return bar
}
