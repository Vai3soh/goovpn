package dns

import (
	"fmt"
)

type DnsCore interface {
	CaseSetupDnsNotUseSystemd() error
	CaseSetupDnsWithUseSystemd() error
	GetCmdReturnResolver() (*string, error)
	CaseAddDnsAddress(text string)
}

type CmdCore interface {
	SplitCommand(cmd string) ([]string, error)
	RunProcessCmd([]string) error
}

type system struct {
	name       string
	useSystemd bool
	CmdCore
	DnsCore
}

func NewSystem(
	name string, cmdCore CmdCore, dnsCore DnsCore, useSystemd bool,
) (obj *system, err error) {
	obj = &system{
		name:       name,
		useSystemd: useSystemd,
		CmdCore:    cmdCore,
		DnsCore:    dnsCore,
	}
	return
}

func (s *system) Name() string {
	return s.name
}

type Names struct {
	goos []system
	dns  map[string]map[string]func() error
	DnsCore
}

func NewNames(dnscore DnsCore) (obj *Names, err error) {
	obj = &Names{

		goos:    make([]system, 0),
		dns:     make(map[string]map[string]func() error),
		DnsCore: dnscore,
	}
	return
}

func (n *Names) SetGoos(sys ...system) {
	n.goos = append(n.goos, sys...)
}

func (n *Names) ConfigureDns(key string) {

	m := make(map[string]func() error)

	for _, sys := range n.goos {
		if key == `revert` {
			m[`linux`] = sys.LinuxRevertDns
			m[`windows`] = sys.WindowsRevertDns
			n.dns[key] = m
		} else {
			m[`linux`] = sys.LinuxSetupDns
			m[`windows`] = sys.WindowsSetupDns
			n.dns[key] = m
		}
	}
}

func (n *Names) SetupDns(key string) error {
	for _, sys := range n.goos {
		return n.dns[key][sys.Name()]()
	}
	return fmt.Errorf(`not found OS`)
}

func (s *system) WindowsSetupDns() error {
	return nil
}

func (s *system) WindowsRevertDns() error {
	return nil
}

func (s *system) LinuxSetupDns() error {
	if !s.useSystemd {
		err := s.DnsCore.CaseSetupDnsNotUseSystemd()
		if err != nil {
			return fmt.Errorf("setup dns err [%w]", err)
		}
	} else {
		err := s.DnsCore.CaseSetupDnsWithUseSystemd()
		if err != nil {
			return fmt.Errorf("setup dns err [%w]", err)
		}
	}
	return nil
}

func (s *system) LinuxRevertDns() error {
	if s.useSystemd {
		return nil
	}
	cmd, err := s.DnsCore.GetCmdReturnResolver()
	if err != nil {
		return fmt.Errorf("don't get cmd return resolver [%w]", err)
	}
	cmdArg, err := s.CmdCore.SplitCommand(*cmd)
	if err != nil {
		return fmt.Errorf("don't split cmd: [%w]", err)
	}
	err = s.CmdCore.RunProcessCmd(cmdArg)
	if err != nil {
		return fmt.Errorf("process don't running: [%w]", err)
	}
	return nil
}

func (n *Names) AddDnsAddrs(text string) {
	n.DnsCore.CaseAddDnsAddress(text)
}
