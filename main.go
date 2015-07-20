package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/ailispaw/xhyvectl/vmnet"
)

const (
	COMMAND = "xhyvectl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <UUID>\n", os.Args[0])
		os.Exit(1)
	}

	user, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: Could not get the current user\n", COMMAND)
		os.Exit(1)
	}
	if user.Uid != "0" {
		fmt.Fprintf(os.Stderr, "%s: Requires root permissions\n", COMMAND)
		fmt.Fprintf(os.Stderr, "sudo")
		for _, arg := range os.Args {
			fmt.Fprintf(os.Stderr, " %s", arg)
		}
		fmt.Fprintf(os.Stderr, "\n")
		out, err := exec.Command("sudo", os.Args...).Output()
		if err != nil {
			os.Exit(1)
		}
		fmt.Printf("%s", string(out))
		os.Exit(0)
	}

	mac, err := vmnet.GetMACAddressFromUUID(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", COMMAND, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", mac)

	ip, err := vmnet.GetIPAddressFromMACAddress(mac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", COMMAND, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", ip)

	net, err := vmnet.GetIPNet()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", COMMAND, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", net)
	fmt.Printf("%s\n", net.IP)
	fmt.Printf("%d.%d.%d.%d\n", net.Mask[0], net.Mask[1], net.Mask[2], net.Mask[3])

	os.Exit(0)
}
