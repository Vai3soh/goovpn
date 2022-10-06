package usecasedns

import "os/exec"

type СmdSetters interface {
	SetCommand(cmd string)
}

type CmdToolsManager interface {
	SplitCmd() ([]string, error)
	PassArgumentsToExec(cmdArg []string) *exec.Cmd
}

type CmdResolver interface {
	CmdResolvConfAndPrintf() (*string, *string, error)
	CmdSystemdResolv() (*string, error)
	CmdDownResolvConf() (*string, error)
}

type ProcessManager interface {
	RunCmdWithPipe(args1, args2 []string) error
	SetToProc(cmdExec *exec.Cmd)
	StartProc() error
}

type DnsSetters interface {
	SetAddress(matches []string)
}
