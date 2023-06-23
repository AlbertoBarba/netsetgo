package network

import (
	"net"
)

type Config struct {
	BridgeName     string
	BridgeIP       net.IP
	ContainerIP    net.IP
	Subnet         *net.IPNet
	VethNamePrefix string
}