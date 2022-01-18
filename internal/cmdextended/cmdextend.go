package cmdextended

import (
	"io"
	"io/ioutil"
	"os/exec"

	"bufio"

	"github.com/google/shlex"
	cmdchain "github.com/rainu/go-command-chain"
)

type Cmd struct {
	proc    *exec.Cmd
	command string
	stdout  io.ReadCloser
	stderr  io.ReadCloser
}

type Option func(*Cmd)

func NewCmd(opts ...Option) *Cmd {
	c := &Cmd{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Cmd) SplitCmd() ([]string, error) {
	args, err := shlex.Split(c.command)
	if err != nil {
		return nil, err
	}
	return args, nil

}

func (c *Cmd) PassArgumentsToExec(args []string) *exec.Cmd {
	return exec.Command(args[0], args[1:]...)
}

func (c *Cmd) Proc() *exec.Cmd {
	return c.proc
}

func (c *Cmd) SetToProc(execCmd *exec.Cmd) {
	c.proc = execCmd
}

func (c *Cmd) SetCommand(command string) {
	c.command = command
}

func (c *Cmd) SetStdout(stdout io.ReadCloser) {
	c.stdout = stdout
}

func (c *Cmd) SetStderr(stderr io.ReadCloser) {
	c.stderr = stderr
}

func (c *Cmd) Stdout() (io.ReadCloser, error) {
	stdout, err := c.proc.StdoutPipe()
	if err != nil {
		return nil, err
	}
	return stdout, nil
}

func (c *Cmd) ReadFromStdout() ([]byte, error) {
	b, err := ioutil.ReadAll(c.stdout)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Cmd) ReadFromStdErr() ([]byte, error) {
	b, err := ioutil.ReadAll(c.stderr)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Cmd) Stderr() (io.ReadCloser, error) {
	stderr, err := c.proc.StderrPipe()
	if err != nil {
		return nil, err
	}
	return stderr, nil
}

func (c *Cmd) StartProc() error {
	if err := c.proc.Start(); err != nil {
		return err
	}
	return nil
}

func (c *Cmd) WaitProc() error {
	if err := c.proc.Wait(); err != nil {
		return err
	}
	return nil
}

func (c *Cmd) Scanner() *bufio.Scanner {
	return bufio.NewScanner(io.MultiReader(c.stdout, c.stderr))
}

func (c *Cmd) RunCmdWithPipe(args1, args2 []string) error {
	err := cmdchain.Builder().
		Join(args1[0], args1[1:]...).
		Join(args2[0], args2[1:]...).
		Finalize().WithError().Run()
	if err != nil {
		return err
	}
	return nil
}
