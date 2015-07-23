package commands

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	CONFIG_FILE_NAME   = "config.yml"
	BOX_FILE_EXTENSION = ".xhv"
)

var (
	boolForce bool
)

var cmdPackage = &cobra.Command{
	Use:     "package <BOX NAME> <DIR>",
	Aliases: []string{"pack"},
	Short:   "Package files into a box for xhyve",
	Long:    COMMAND + " package - Package files into a box for xhyve",
	Run:     doPackage,
}

func init() {
	flags := cmdPackage.Flags()
	flags.BoolVarP(&boolForce, "force", "f", false, "Override a box even if exists")
}

func doPackage(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments to package <DIR>/* into <BOX NAME>.xhv")
	}

	boxName := args[0] + BOX_FILE_EXTENSION

	if boxFi, err := os.Lstat(boxName); err == nil {
		if boxFi.Mode().IsDir() {
			log.Fatalf("%s is a directory", boxName)
		}
		if !boolForce {
			log.Fatalf("%s exists already", boxName)
		}
	}

	srcPath := filepath.Clean(args[1])

	if srcFi, err := os.Lstat(srcPath); err != nil {
		log.Fatalf("Can't get a file info: %s, error: %s", srcPath, err)
	} else if !srcFi.Mode().IsDir() {
		log.Fatalf("%s is not a directory", srcPath)
	}

	configPath := srcPath + string(filepath.Separator) + CONFIG_FILE_NAME
	if _, err := os.Lstat(configPath); err != nil {
		log.Fatalf("%s must be in the %s", CONFIG_FILE_NAME, srcPath)
	}

	dir, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	boxFile, err := os.Create(boxName)
	if err != nil {
		log.Fatal(err)
	}
	defer boxFile.Close()

	var fileWriter io.WriteCloser = boxFile
	fileWriter = gzip.NewWriter(boxFile)
	defer fileWriter.Close()

	boxFileWriter := tar.NewWriter(fileWriter)
	defer boxFileWriter.Close()

	var (
		nFiles uint64 = 0
		nTotal uint64 = 0
	)

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()
		filePath := srcPath + string(filepath.Separator) + fileName
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Could not open a file: %s, error: %s", filePath, err)
		}
		defer file.Close()

		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			log.Fatalf("Could not get a file info header to archive: %s, error: %s", filePath, err)
		}

		header.Name = fileInfo.Name()

		size := uint64(fileInfo.Size())

		if err := boxFileWriter.WriteHeader(header); err != nil {
			log.Fatalf("Could not archive a file header: %s, error: %s", filePath, err)
		}

		bar := SetupProgressBar(int(size),
			fmt.Sprintf("%-20s %6s ", fileInfo.Name(), humanize.Bytes(size)))
		bar.Start()

		barWriter := io.MultiWriter(boxFileWriter, bar)

		if _, err = io.Copy(barWriter, file); err != nil {
			log.Fatalf("Could not archive a file: %s, error: %s", filePath, err)
		}

		bar.Finish()
		file.Close()

		nFiles++
		nTotal += size
	}

	boxFi, err := os.Lstat(boxName)
	if err != nil {
		log.Fatalf("Could not make a package %s", boxName)
	}

	fmt.Printf("%-20s %6s => %s\n",
		fmt.Sprintf("Total %d files", nFiles),
		humanize.Bytes(nTotal), humanize.Bytes(uint64(boxFi.Size())))

	fmt.Printf("%s: Packaged %s/* into %s\n", COMMAND, srcPath, boxName)
}
