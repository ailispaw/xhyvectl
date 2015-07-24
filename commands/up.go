package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ailispaw/xhyvectl/vmnet"
)

var (
	cpus   int
	memory string
)

var cmdUp = &cobra.Command{
	Use:     "up <VM NAME>",
	Aliases: []string{"start"},
	Short:   "Start a VM",
	Long:    COMMAND + " up - Start a VM",
	Run:     up,
}

func init() {
	flags := cmdUp.Flags()
	flags.IntVarP(&cpus, "cpus", "c", 1, "# of CPUs")
	flags.StringVarP(&memory, "memory", "m", "1G", "Memory Size (ex. 512M, 2G)")
}

func up(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument to start <VM NAME>")
	}

	Sudo(ctx)

	vmName := args[0]
	vmDir := os.ExpandEnv(VM_ROOT_DIR + vmName)

	if _, err := os.Lstat(vmDir); err != nil {
		log.Fatalf("%s has not initialized yet. Please initialize a VM first.", vmName)
	}

	configPath := vmDir + string(filepath.Separator) + CONFIG_FILE_NAME

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	uuidPath := vmDir + string(filepath.Separator) + UUID_FILE_NAME

	uuidByte, err := ioutil.ReadFile(uuidPath)
	if err != nil {
		log.Fatal(err)
	}
	uuid := string(uuidByte)

	var mac string
	macPath := vmDir + string(filepath.Separator) + MAC_FILE_NAME

	macByte, err := ioutil.ReadFile(macPath)
	if err != nil {
		macFile, err := os.Create(macPath)
		if err != nil {
			log.Fatal(err)
		}
		defer macFile.Close()

		mac, err = vmnet.GetMACAddressByUUID(uuid)
		if err != nil {
			log.Fatal(err)
		}

		macFile.WriteString(mac)
		macFile.Close()
	} else {
		mac = string(macByte)
	}

	var xArgs []string

	xArgs = append(xArgs, "-U", uuid)

	if config.VMConfig.ACPI {
		xArgs = append(xArgs, "-A")
	}
	xArgs = append(xArgs, "-c", fmt.Sprintf("%d", cpus))
	xArgs = append(xArgs, "-m", memory)
	for _, hdd := range config.VMConfig.HDD {
		params := strings.Split(hdd, ",")
		params[len(params)-1] = vmDir + string(filepath.Separator) + params[len(params)-1]
		xArgs = append(xArgs, "-s", strings.Join(params, ","))
	}
	for _, net := range config.VMConfig.Net {
		xArgs = append(xArgs, "-s", net)
	}
	for _, pci := range config.VMConfig.PCI {
		xArgs = append(xArgs, "-s", pci)
	}
	for _, lpc := range config.VMConfig.LPC {
		xArgs = append(xArgs, "-l", lpc)
	}

	if config.Type == "linux" {
		xArgs = append(xArgs, "-f", fmt.Sprintf("kexec,%s,%s,%s",
			vmDir+string(filepath.Separator)+config.BootConfig.Kernel,
			vmDir+string(filepath.Separator)+config.BootConfig.Initrd,
			config.BootConfig.CommandLine))
	} else if config.Type == "freebsd" {
		xArgs = append(xArgs, "-f", fmt.Sprintf("fbsd,%s,%s,%s",
			vmDir+string(filepath.Separator)+config.BootConfig.UserBoot,
			vmDir+string(filepath.Separator)+config.BootConfig.BootVolume,
			config.BootConfig.KernelENV))
	}

	for {
		cmd := exec.Command("xhyve", xArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			break
		}
	}
}
