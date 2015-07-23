package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	BOX_ROOT_DIR = "$HOME/.xhyvectl/boxes/"
)

var cmdBox = &cobra.Command{
	Use:   "box [command]",
	Short: "Manage boxes",
	Long:  COMMAND + " box - Manage boxes",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

var cmdBoxList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed boxes",
	Long:    COMMAND + " list - List installed boxes",
	Run:     listBox,
}

var cmdBoxInstall = &cobra.Command{
	Use:     "install <BOX NAME> <BOX FILE|URL>",
	Aliases: []string{"add"},
	Short:   "Install a box",
	Long:    COMMAND + " install - Install a box",
	Run:     installBox,
}

var cmdBoxUpdate = &cobra.Command{
	Use:   "update <BOX NAME> <BOX FILE|URL>",
	Short: "Update a box",
	Long:  COMMAND + " update - Update a box",
	Run:   updateBox,
}

var cmdBoxUninstall = &cobra.Command{
	Use:     "uninstall <BOX NAME>...",
	Aliases: []string{"remove", "rm"},
	Short:   "Uninstall box(es)",
	Long:    COMMAND + " uninstall - Uninstall box(es)",
	Run:     uninstallBoxes,
}

func init() {
	cmdBox.AddCommand(cmdBoxList)

	flags := cmdBoxInstall.Flags()
	flags.BoolVarP(&boolForce, "force", "f", false, "Override a box even if exists")
	cmdBox.AddCommand(cmdBoxInstall)

	cmdBox.AddCommand(cmdBoxUpdate)

	cmdBox.AddCommand(cmdBoxUninstall)
}

func listBox(ctx *cobra.Command, args []string) {
	boxes, _ := filepath.Glob(os.ExpandEnv(BOX_ROOT_DIR) + "*" + BOX_FILE_EXTENSION)
	for _, box := range boxes {
		boxFile := filepath.Base(box)
		fmt.Println(boxFile[0 : len(boxFile)-len(BOX_FILE_EXTENSION)])
	}
}

func installBox(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments to install <BOX FILE|URL> as <BOX NAME>")
	}

	boxName := args[0]
	dstPath := os.ExpandEnv(BOX_ROOT_DIR + boxName + BOX_FILE_EXTENSION)

	if _, err := os.Lstat(dstPath); err == nil {
		if !boolForce {
			log.Fatalf("%s has already installed", boxName)
		}
		if err := os.Remove(dstPath); err != nil {
			log.Fatalf("Could not remove %s, error %s", boxName, err)
		}
	}

	if err := os.MkdirAll(os.ExpandEnv(BOX_ROOT_DIR), 0755); err != nil {
		log.Fatal(err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Could not create a file, error %s", err)
	}
	defer dstFile.Close()

	srcPath := args[1]

	if strings.HasPrefix(srcPath, "http://") || strings.HasPrefix(srcPath, "https://") {
		resp, err := http.Get(srcPath)
		if err != nil {
			log.Fatalf("Could not download: %s, error %s", srcPath, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			log.Fatalf("Could not download: %s, HTTP status %s", srcPath, resp.Status)
		}

		if _, err := io.Copy(dstFile, resp.Body); err != nil {
			log.Fatalf("Could not download: %s, error %s", srcPath, err)
		}
	} else {
		srcPath = filepath.Clean(srcPath)
		srcFile, err := os.Open(srcPath)
		if err != nil {
			log.Fatalf("Could not open: %s, error %s", srcPath, err)
		}
		defer srcFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			log.Fatalf("Could not download: %s, error %s", srcPath, err)
		}
	}

	fmt.Printf("%s: Installed %s\n", COMMAND, boxName)
}

func updateBox(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments to update <BOX NAME> with <BOX FILE|URL>")
	}

	boxName := args[0]
	dstPath := os.ExpandEnv(BOX_ROOT_DIR + boxName + BOX_FILE_EXTENSION)

	if _, err := os.Lstat(dstPath); err != nil {
		log.Fatalf("%s has not installed yet", boxName)
	}

	boolForce = true
	installBox(ctx, args)
}

func uninstallBoxes(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <BOX NAME> at least to remove")
	}

	var gotError = false
	for _, boxName := range args {
		boxPath := os.ExpandEnv(BOX_ROOT_DIR + boxName + BOX_FILE_EXTENSION)
		if err := os.Remove(boxPath); err != nil {
			log.Error(err)
			gotError = true
			continue
		}
		fmt.Printf("%s: Removed %s\n", COMMAND, boxName)
	}
	if gotError {
		log.Fatal("Error: failed to remove one or more images")
	}
}
