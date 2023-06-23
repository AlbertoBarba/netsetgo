package netsetgo

import (
	"fmt"
	"net"

	"code.cloudfoundry.org/guardian/kawasaki/netns"
	"github.com/AlbertoBarba/netsetgo/configurer"
	"github.com/AlbertoBarba/netsetgo/device"
	"github.com/AlbertoBarba/netsetgo/network"
)

//go:generate counterfeiter . Configurer
type Configurer interface {
	Apply(netConfig network.Config, pid int) error
}

type Netset struct {
	HostConfigurer      Configurer
	ContainerConfigurer Configurer
}

func New(hostConfigurer, containerConfigurer Configurer) *Netset {
	return &Netset{
		HostConfigurer:      hostConfigurer,
		ContainerConfigurer: containerConfigurer,
	}
}

func (n *Netset) ConfigureHost(netConfig network.Config, pid int) error {
	return n.HostConfigurer.Apply(netConfig, pid)
}

func (n *Netset) ConfigureContainer(netConfig network.Config, pid int) error {
	return n.ContainerConfigurer.Apply(netConfig, pid)
}

type OptionFunc func(*NetSetGoOptions) error

func SetBridgeName(bridgeName string) OptionFunc {
	return func(c *NetSetGoOptions) error {
		c.bridgeName = bridgeName
		return nil
	}
}

func SetBridgeAddress(bridgeAddress string) OptionFunc {
	return func(c *NetSetGoOptions) error {
		c.bridgeAddress = bridgeAddress
		return nil
	}
}

func SetVethNamePrefix(vethNamePrefix string) OptionFunc {
	return func(c *NetSetGoOptions) error {
		c.vethNamePrefix = vethNamePrefix
		return nil
	}
}

func SetContainerAddress(containerAddress string) OptionFunc {
	return func(c *NetSetGoOptions) error {
		c.containerAddress = containerAddress
		return nil
	}
}

type NetSetGoOptions struct {
	bridgeName       string
	bridgeAddress    string
	vethNamePrefix   string
	containerAddress string
}

func ConfigureForPid(pid int, opts ...OptionFunc) error {
	if pid == 0 {
		return fmt.Errorf("netsetgo needs a pid")
	}

	options := &NetSetGoOptions{
		bridgeName:       "brg0",
		bridgeAddress:    "10.10.10.1/24",
		vethNamePrefix:   "veth",
		containerAddress: "10.10.10.2/24",
	}

	// Run the options on it
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return err
		}
	}

	bridgeCreator := device.NewBridge()
	vethCreator := device.NewVeth()
	netnsExecer := &netns.Execer{}

	hostConfigurer := configurer.NewHostConfigurer(bridgeCreator, vethCreator)
	containerConfigurer := configurer.NewContainerConfigurer(netnsExecer)
	netset := New(hostConfigurer, containerConfigurer)

	bridgeIP, bridgeSubnet, err := net.ParseCIDR(options.bridgeAddress)
	if err != nil {
		return err
	}

	containerIP, _, err := net.ParseCIDR(options.containerAddress)
	if err != nil {
		return err
	}

	netConfig := network.Config{
		BridgeName:     options.bridgeName,
		BridgeIP:       bridgeIP,
		ContainerIP:    containerIP,
		Subnet:         bridgeSubnet,
		VethNamePrefix: options.vethNamePrefix,
	}

	if err := netset.ConfigureHost(netConfig, pid); err != nil {
		return err
	}

	return netset.ConfigureContainer(netConfig, pid)
}
