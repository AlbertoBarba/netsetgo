package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlbertoBarba/netsetgo"
)

func main() {
	var bridgeName, bridgeAddress, containerAddress, vethNamePrefix string
	var pid int

	flag.StringVar(&bridgeName, "bridgeName", "brg0", "Name to assign to bridge device")
	flag.StringVar(&bridgeAddress, "bridgeAddress", "10.10.10.1/24", "Address to assign to bridge device (CIDR notation)")
	flag.StringVar(&vethNamePrefix, "vethNamePrefix", "veth", "Name prefix for veth devices")
	flag.StringVar(&containerAddress, "containerAddress", "10.10.10.2/24", "Address to assign to the container (CIDR notation)")
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	check(netsetgo.ConfigureForPid(
		pid,
		netsetgo.SetBridgeName(bridgeName),
		netsetgo.SetBridgeAddress(bridgeAddress),
		netsetgo.SetVethNamePrefix(vethNamePrefix),
		netsetgo.SetContainerAddress(containerAddress),
	))
}

func check(err error) {
	if err != nil {
		fmt.Printf("ERROR - %s\n", err.Error())
		os.Exit(1)
	}
}
