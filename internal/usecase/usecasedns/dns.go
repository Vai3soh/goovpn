package usecasedns

import (
	"regexp"
	"strings"
)

type DnsUseCase struct {
	cmdSetters      СmdSetters
	cmdToolsManager CmdToolsManager
	cmdResolver     CmdResolver
	processManager  ProcessManager
	dnsSetters      DnsSetters
}

func NewDnsUseCase(
	cmdSetters СmdSetters, cmdToolsManager CmdToolsManager,
	cmdResolver CmdResolver, processManager ProcessManager, dnsSetters DnsSetters,
) (obj *DnsUseCase, err error) {
	obj = &DnsUseCase{
		cmdSetters:      cmdSetters,
		cmdToolsManager: cmdToolsManager,
		cmdResolver:     cmdResolver,
		processManager:  processManager,
		dnsSetters:      dnsSetters,
	}
	return
}

func (dns *DnsUseCase) SplitCommand(cmd string) ([]string, error) {
	dns.cmdSetters.SetCommand(cmd)
	cmdArg, err := dns.cmdToolsManager.SplitCmd()
	if err != nil {
		return nil, err
	}
	return cmdArg, nil
}

func (dns *DnsUseCase) CaseSetupDnsNotUseSystemd() error {
	cmdPrintf, cmdResolv, err := dns.cmdResolver.CmdResolvConfAndPrintf()
	if err != nil {
		return err
	}
	cmdArg1, err := dns.SplitCommand(*cmdPrintf)
	if err != nil {
		return err
	}
	cmdArg2, err := dns.SplitCommand(*cmdResolv)
	if err != nil {
		return err
	}
	if err := dns.processManager.RunCmdWithPipe(cmdArg1, cmdArg2); err != nil {
		return err
	}
	return nil
}

func (dns *DnsUseCase) CaseSetupDnsWithUseSystemd() error {
	cmdDns, CmdDomain, err := dns.cmdResolver.CmdSystemdResolv()
	if err != nil {
		return err
	}
	cmdArgDns, err := dns.SplitCommand(*cmdDns)
	if err != nil {
		return err
	}
	err = dns.RunProcessCmd(cmdArgDns)
	if err != nil {
		return err
	}
	cmdArgDomain, err := dns.SplitCommand(*CmdDomain)
	if err != nil {
		return err
	}
	err = dns.RunProcessCmd(cmdArgDomain)
	if err != nil {
		return err
	}
	return nil
}

func (dns *DnsUseCase) RunProcessCmd(cmdArg []string) error {
	cmdExec := dns.cmdToolsManager.PassArgumentsToExec(cmdArg)
	dns.processManager.SetToProc(cmdExec)
	err := dns.processManager.StartProc()
	if err != nil {
		return err
	}
	return nil
}

func (dns *DnsUseCase) GetCmdReturnResolver() (*string, error) {
	cmd, err := dns.cmdResolver.CmdDownResolvConf()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (dns *DnsUseCase) CaseAddDnsAddress(text string) {
	reg := regexp.MustCompile(`(?m)DNS Servers:\s+[^*]+\d$`)
	matches := reg.FindAllString(text, -1)
	if len(matches) != 0 {
		matches = strings.Split(matches[0], "\n")
		matches = append(matches[:0], matches[0+1:]...)
		for i := range matches {
			matches[i] = strings.TrimSpace(matches[i])
		}
		dns.dnsSetters.SetAddress(matches)
	} else {
		matches = []string{"1.1.1.1"}
		dns.dnsSetters.SetAddress(matches)
	}
}
