package main

import (
	"fmt"
	"os"

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
