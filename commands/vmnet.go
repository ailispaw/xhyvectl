package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ailispaw/xhyvectl/vmnet"
)

var cmdVmnet = &cobra.Command{
	Use:   "vmnet <UUID>",
	Short: "Test vmnet",
	Long:  COMMAND + " vmnet - Test vmnet",
	Run:   testVmnet,
}

func init() {
}

func testVmnet(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <UUID> to test")
	}

	Sudo(ctx)

	mac, err := vmnet.GetMACAddressByUUID(args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MAC Address: %s\n", mac)

	ip, err := vmnet.GetIPAddressByMACAddress(mac)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("IP Address:  %s\n", ip)

	net, err := vmnet.GetIPNet()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("CIDR:        %s\n", net)
	fmt.Printf("Subnet Addr: %s\n", net.IP)
	fmt.Printf("Subnet Mask: %d.%d.%d.%d\n", net.Mask[0], net.Mask[1], net.Mask[2], net.Mask[3])
}
