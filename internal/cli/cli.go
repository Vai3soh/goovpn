package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/shlex"
	cmdchain "github.com/rainu/go-command-chain"
)

type cmd struct {
	proc    *exec.Cmd
	command string
	stdout  io.ReadCloser
	stderr  io.ReadCloser
}

type CmdSeacher interface {
	SetCommand(cmd string)
	getBinaryPath() (*string, error)
}

type Option func(*cmd)

func NewCmd(opts ...Option) *cmd {
	c := &cmd{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *cmd) SplitCmd() ([]string, error) {
	args, err := shlex.Split(c.command)
	if err != nil {
		return nil, fmt.Errorf("shlex command failed: [%w]", err)
	}
	return args, nil

}

func (c *cmd) PassArgumentsToExec(args []string) *exec.Cmd {
	return exec.Command(args[0], args[1:]...)
}

func (c *cmd) Proc() *exec.Cmd {
	return c.proc
}

func (c *cmd) SetToProc(execCmd *exec.Cmd) {
	c.proc = execCmd
}

func (c *cmd) SetCommand(command string) {
	c.command = command
}

func (c *cmd) SetStdout(stdout io.ReadCloser) {
	c.stdout = stdout
}

func (c *cmd) SetStderr(stderr io.ReadCloser) {
	c.stderr = stderr
}

func (c *cmd) Stdout() (io.ReadCloser, error) {
	stdout, err := c.proc.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("pipe err: [%w]", err)
	}
	return stdout, nil
}

func (c *cmd) ReadFromStdout() ([]byte, error) {
	b, err := io.ReadAll(c.stdout)
	if err != nil {
		return nil, fmt.Errorf("read all err: [%w]", err)
	}
	return b, nil
}

func (c *cmd) ReadFromStdErr() ([]byte, error) {
	b, err := io.ReadAll(c.stderr)

	if err != nil {
		return nil, fmt.Errorf("read all err: [%w]", err)
	}
	return b, nil
}

func (c *cmd) Stderr() (io.ReadCloser, error) {
	stderr, err := c.proc.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("pipe err: [%w]", err)
	}
	return stderr, nil
}

func (c *cmd) StartProc() error {
	if err := c.proc.Start(); err != nil {
		return fmt.Errorf("process start err: [%w]", err)
	}
	return nil
}

func (c *cmd) WaitProc() error {
	if err := c.proc.Wait(); err != nil {
		return fmt.Errorf("process wait err: [%w]", err)
	}
	return nil
}

func (c *cmd) Scanner() *bufio.Scanner {
	return bufio.NewScanner(io.MultiReader(c.stdout, c.stderr))
}

func (c *cmd) RunCmdWithPipe(args1, args2 []string) error {
	err := cmdchain.Builder().
		Join(args1[0], args1[1:]...).
		Join(args2[0], args2[1:]...).
		Finalize().WithError().Run()
	if err != nil {
		return fmt.Errorf("run cmd with pipe err: [%w]", err)
	}
	return nil
}

func (c *cmd) getBinaryPath() (*string, error) {
	if runtime.GOOS == `linux` {
		os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin")
	}
	full, err := exec.LookPath(c.command)
	if err != nil {
		return nil, fmt.Errorf("%s path not found: [%w]", c.command, err)
	}
	return &full, nil
}
