package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type cmdResolver struct {
	iface   string
	address []string
	CmdSeacher
}

type OptionResolver func(*cmdResolver)

func NewResolver(opts ...OptionResolver) *cmdResolver {
	f := &cmdResolver{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithInterface(iface string) OptionResolver {
	return func(c *cmdResolver) {
		c.iface = iface
	}
}

func WithAddress(address []string) OptionResolver {
	return func(c *cmdResolver) {
		c.address = address
	}
}

func WithCliResolver(cli CmdSeacher) OptionResolver {
	return func(c *cmdResolver) {
		c.CmdSeacher = cli
	}
}

func (c *cmdResolver) SetAddress(addr []string) {
	c.address = addr
}

func (c *cmdResolver) SetInterface(iface string) {
	c.iface = iface
}

func (c *cmdResolver) CmdDownResolvConf() (*string, error) {
	c.CmdSeacher.SetCommand("resolvconf")
	full, err := c.CmdSeacher.getBinaryPath()
	if err != nil {
		return nil, fmt.Errorf("cmd down resolvconf err: [%w]", err)
	}
	cmdDown := *full + " -d " + c.iface + ".openvpn"
	return &cmdDown, nil
}

func (c *cmdResolver) CmdSystemdResolv() (*string, error) {
	c.CmdSeacher.SetCommand("systemd-resolve")
	full, err := c.CmdSeacher.getBinaryPath()
	if err != nil {
		return nil, fmt.Errorf("cmd systemd-resolve err: [%w]", err)
	}
	cmd := *full + " " + "-i " + c.iface + " --set-Resolver=" +
		strings.Join(c.address[:], " --set-Resolver=") + " --set-domain=~"
	return &cmd, nil
}

func (c *cmdResolver) CmdResolvConfAndPrintf() (*string, *string, error) {
	c.CmdSeacher.SetCommand("resolvconf")
	full, err := c.CmdSeacher.getBinaryPath()
	if err != nil {
		return nil, nil, fmt.Errorf("cmd resolvconf err: [%w]", err)
	}
	c.CmdSeacher.SetCommand("printf")
	printf, err := c.CmdSeacher.getBinaryPath()
	if err != nil {
		return nil, nil, fmt.Errorf("cmd resolvconf err: [%w]", err)
	}
	for i, value := range c.address {
		c.address[i] = value + "\\n"
	}
	cmdPrintf := *printf + " " + strconv.Quote("nameserver "+strings.Join(c.address[:], "nameserver "))
	cmdResolv := *full + " -a " + c.iface + ".openvpn -m 0 -x"
	return &cmdPrintf, &cmdResolv, nil
}
