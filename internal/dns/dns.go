package dns

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Dns struct {
	iface   string
	address []string
}

type Option func(*Dns)

func NewDns(opts ...Option) *Dns {
	f := &Dns{}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithInterface(iface string) Option {
	return func(d *Dns) {
		d.iface = iface
	}
}

func WithAddress(address []string) Option {
	return func(d *Dns) {
		d.address = address
	}
}

func (d *Dns) SetAddress(addr []string) {
	d.address = addr
}

func (d *Dns) SetInterface(iface string) {
	d.iface = iface
}

func (d *Dns) CmdSystemdResolv() (*string, error) {
	os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin")
	full, err := exec.LookPath("systemd-resolve")
	if err != nil {
		return nil, err
	}
	cmd := full + " " + "-i " + d.iface + " --set-dns=" +
		strings.Join(d.address[:], " --set-dns=") + " --set-domain=~"
	return &cmd, nil
}

func (d *Dns) CmdResolvConf() (*string, *string, error) {
	os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin")
	full, err := exec.LookPath("resolvconf")
	if err != nil {
		return nil, nil, err
	}
	printf, err := exec.LookPath("printf")
	if err != nil {
		return nil, nil, err
	}
	for i, value := range d.address {
		d.address[i] = value + "\\n"
	}
	cmdPrintf := printf + " " + strconv.Quote("nameserver "+strings.Join(d.address[:], "nameserver "))
	cmdResolv := full + " -a " + d.iface + ".openvpn -m 0 -x"
	return &cmdPrintf, &cmdResolv, nil
}

func (d *Dns) CmdDownResolvConf() (*string, error) {
	os.Setenv("PATH", "/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin")
	full, err := exec.LookPath("resolvconf")
	if err != nil {
		return nil, err
	}
	cmdDown := full + " -d " + d.iface + ".openvpn"
	return &cmdDown, nil
}
