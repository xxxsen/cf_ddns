package provider

import (
	"fmt"
	"net"
	"strings"

	"github.com/xxxsen/common/utils"
)

type machineProvider struct {
	typ          string
	eth          string
	allowPrivate bool
}

func (p *machineProvider) Name() string {
	return ProviderMachine
}

func (p *machineProvider) Get() (string, error) {
	switch p.typ {
	case "ipv4":
		return p.getV4()
	case "ipv6":
		return p.getV6()
	}
	return "", fmt.Errorf("invalid type")
}

func (p *machineProvider) getAll() ([]*net.IPNet, []*net.IPNet, error) {
	intf, err := net.InterfaceByName(p.eth)
	if err != nil {
		return nil, nil, fmt.Errorf("get addr failed")
	}
	addrs, err := intf.Addrs()
	if err != nil {
		return nil, nil, err
	}
	_, ipv6Unicast, _ := net.ParseCIDR("2000::/3")
	var v4 []*net.IPNet
	var v6 []*net.IPNet
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if !ipnet.IP.IsGlobalUnicast() {
			continue
		}
		if !p.allowPrivate && ipnet.IP.IsPrivate() {
			continue
		}
		_, bits := ipnet.Mask.Size()
		if bits == 128 && ipv6Unicast.Contains(ipnet.IP) {
			v6 = append(v6, ipnet)
		}
		if bits == 32 {
			v4 = append(v4, ipnet)
		}
	}
	return v4, v6, nil
}

func (p *machineProvider) getV6() (string, error) {
	_, v6, err := p.getAll()
	if err != nil {
		return "", err
	}
	if len(v6) == 0 {
		return "", fmt.Errorf("ipv6 not found")
	}
	return v6[0].IP.String(), nil
}

func (p *machineProvider) getV4() (string, error) {
	v4, _, err := p.getAll()
	if err != nil {
		return "", err
	}
	if len(v4) == 0 {
		return "", fmt.Errorf("ipv4 not found")
	}
	return v4[0].IP.String(), nil
}

func createMachineProvider(args interface{}) (IProvider, error) {
	c := &MachineConfig{}
	if err := utils.ConvStructJson(args, c); err != nil {
		return nil, err
	}
	c.Type = strings.ToLower(c.Type)
	if c.Type != "ipv4" && c.Type != "ipv6" {
		return nil, fmt.Errorf("invalid ip type:%s", c.Type)
	}
	if len(c.Interface) == 0 {
		return nil, fmt.Errorf("empty interface")
	}
	return &machineProvider{typ: c.Type, eth: c.Interface, allowPrivate: c.AllowPrivate}, nil

}

func init() {
	Register(ProviderMachine, createMachineProvider)
}
