package commands

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	VM_ROOT_DIR = "$HOME/.xhyvectl/vms/"
)

var cmdInit = &cobra.Command{
	Use:   "init <VM NAME> <BOX NAME>",
	Short: "Initialize a VM with a box",
	Long:  COMMAND + " install - Initialize a VM with a box",
	Run:   doInit,
}

func init() {
	flags := cmdInit.Flags()
	flags.BoolVarP(&boolForce, "force", "f", false, "Override a box even if exists")
}

func doInit(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments to initialize <VM NAME> with <BOX NAME>")
	}

	vmName := args[0]
	vmDir := os.ExpandEnv(VM_ROOT_DIR + vmName)

	if _, err := os.Lstat(vmDir); err == nil {
		if !boolForce {
			log.Fatalf("%s has already initialized", vmName)
		}
	}

	boxName := args[1]
	boxPath := os.ExpandEnv(BOX_ROOT_DIR + boxName + BOX_FILE_EXTENSION)

	if _, err := os.Lstat(boxPath); err != nil {
		log.Fatalf("%s has not installed yet. Please install a box first.", boxName)
	}

	if err := os.MkdirAll(vmDir, 0755); err != nil {
		log.Fatal(err)
	}

	boxFile, err := os.Open(boxPath)
	if err != nil {
		log.Fatal(err)
	}
	defer boxFile.Close()

	var fileReader io.ReadCloser = boxFile
	if fileReader, err = gzip.NewReader(boxFile); err != nil {
		log.Fatal(err)
	}
	defer fileReader.Close()

	boxFileReader := tar.NewReader(fileReader)

	for {
		header, err := boxFileReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		fileName := header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			filePath := vmDir + string(filepath.Separator) + fileName

			file, err := os.Create(filePath)
			if err != nil {
				log.Fatalf("Could not create a file: %s, error: %s", filePath, err)
			}
			defer file.Close()

			size := uint64(header.Size)

			bar := SetupProgressBar(int(size),
				fmt.Sprintf("%-20s %6s ", fileName, humanize.Bytes(size)))
			bar.Start()

			barWriter := io.MultiWriter(file, bar)

			if _, err = io.Copy(barWriter, boxFileReader); err != nil {
				log.Fatalf("Could not extract a file: %s, error: %s", filePath, err)
			}

			bar.Finish()

			if err := os.Chmod(filePath, os.FileMode(header.Mode)); err != nil {
				log.Fatalf("Could not extract a file: %s, error: %s", filePath, err)
			}

			file.Close()
		}
	}

	uuidPath := vmDir + string(filepath.Separator) + UUID_FILE_NAME

	if _, err := os.Lstat(uuidPath); err != nil {
		uuidFile, err := os.Create(uuidPath)
		if err != nil {
			log.Fatal(err)
		}
		defer uuidFile.Close()

		out, err := exec.Command("uuidgen").Output()
		if err != nil {
			log.Fatal(err)
		}

		uuidFile.WriteString(strings.TrimSpace(string(out)))
		uuidFile.Close()
	}

	fmt.Printf("%s: Initialized %s with %s\n", COMMAND, vmName, boxName)
}
