package openvpn

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/Vai3soh/goovpn/entity"
)

type Logger interface {
	Fatal(...interface{})
	Debugf(string, ...interface{})
	Fatalf(string, ...interface{})
	Info(...interface{})
}

type core interface {
	SetSessionCread(u, p string)
	GetOvpnAuthPathFileName() string
	SetPathToFile(path string)
	ReadFile() ([]byte, error)
	SetBodyToCfg(path string)
	GetVpnCread() (string, string)
	CheckOvpnUseAuthUserPass() bool
	SetProfileBody(profileBody string)
	NewSession()
	RunSession(ctx context.Context)
	OffComboBoxAndClear()
	DisableConnectionButton()
	GetTextFromComboBox() string
	EnableDisconnectButton()
	GetChanVpnLog() chan string
	CaseSetLogsInTextWidget(text string)
	EnableConnectButton()
	TurnOnConfigsBox()
	DisableDisconnectButton()
	TraySetImageDisconnect() error
	TraySetImageConnect() error
	SetPhyseInterface(iface string)
	CaseFlickeringIcon() error
}

type DnsManager interface {
	SetupDns(key string) error
	ConfigureDns(key string)
	AddDnsAddrs(s string)
}

type ProfileCore interface {
	CheckUseCfgFile() bool
	GetProfileFromCache(cfg string) entity.Profile
	SearchFileAbsolutePath(file string) (*string, error)
	SaveProfileWithoutCfgFile(cfg string) error
	SaveProfileWithCfgFile(cfg string) error
}

type TransportOvpnClient struct {
	configsPath string
	stopTimeout int
	core        core
	dnsManager  DnsManager
	pcore       ProfileCore
	l           Logger
}

func New(
	configsPath string, stopTimeout int,
	core core, pcore ProfileCore, dnsManager DnsManager,
) *TransportOvpnClient {

	return &TransportOvpnClient{
		configsPath: configsPath,
		stopTimeout: stopTimeout,
		core:        core,
		dnsManager:  dnsManager,
		pcore:       pcore,
	}
}

func (t *TransportOvpnClient) setVpnCread(username, password string) {
	t.core.SetSessionCread(username, password)
}

func (t *TransportOvpnClient) getVpnCread(ok bool) (string, string) {
	if ok {
		fileAuth := t.core.GetOvpnAuthPathFileName()
		if fileAuth == "" {
			msg := `not get path from directive` +
				` auth-user-pass, edit config - auth-user-pass auth.txt`
			t.l.Fatal(msg)
		}
		absPathCredFile, err := t.pcore.SearchFileAbsolutePath(fileAuth)
		if err != nil {
			t.l.Fatalf("don't get file absolute path: [%w]\n", err)
		}
		t.core.SetPathToFile(*absPathCredFile)
		CredFileBody, err := t.core.ReadFile()
		if err != nil {
			t.l.Fatal(err)
		}
		t.core.SetBodyToCfg(string(CredFileBody))
		username, password := t.core.GetVpnCread()
		return username, password
	}
	return "", ""
}

func (t *TransportOvpnClient) initSession() {
	ok := t.core.CheckOvpnUseAuthUserPass()
	username, password := t.getVpnCread(ok)
	t.setVpnCread(username, password)
	t.core.NewSession()
}

func (t *TransportOvpnClient) caseRunSessionOpenvpn(ctx context.Context, profile string) {
	t.core.SetProfileBody(profile)
	t.initSession()
	t.core.RunSession(ctx)
}

func (t *TransportOvpnClient) Connect(ctx context.Context) func() error {

	f := func() error {
		go t.readLogsFromChan()
		t.core.OffComboBoxAndClear()
		t.core.DisableConnectionButton()
		cfg := t.configsPath + t.core.GetTextFromComboBox()
		profile := t.pcore.GetProfileFromCache(cfg)
		if profile.Body == "" {
			err := t.pcore.SaveProfileWithoutCfgFile(cfg)
			if err != nil {
				return err
			}
			if !t.pcore.CheckUseCfgFile() {
				profile := t.pcore.GetProfileFromCache(cfg)
				os.Chdir(t.configsPath)
				t.caseRunSessionOpenvpn(ctx, profile.Body)
				t.core.EnableDisconnectButton()
			} else {
				os.Chdir(t.configsPath)
				err := t.pcore.SaveProfileWithCfgFile(cfg)
				if err != nil {
					return err
				}
				profile := t.pcore.GetProfileFromCache(cfg)
				t.caseRunSessionOpenvpn(ctx, profile.Body)
				t.core.EnableDisconnectButton()
			}
		} else {
			t.caseRunSessionOpenvpn(ctx, profile.Body)
			t.core.EnableDisconnectButton()
		}

		return nil
	}
	return f
}

func (t *TransportOvpnClient) Disconnect(stop context.CancelFunc) func() {
	return func() {
		stop()
		t.dnsManager.ConfigureDns(`revert`)
		err := t.dnsManager.SetupDns(`revert`)
		if err != nil {
			t.l.Fatalf("don't manager dns: [%w]\n", err)
		}
	}
}

func (t *TransportOvpnClient) readLogsFromChan() {
	logChan := t.core.GetChanVpnLog()
	iterSkip := false
	for text := range logChan {
		t.core.CaseSetLogsInTextWidget(text)

		if strings.Contains(text, `DNS Servers:`) {
			t.dnsManager.AddDnsAddrs(text)
		}

		if strings.Contains(text, `Openvpn3 session ended`) {
			t.core.EnableConnectButton()
			t.core.TurnOnConfigsBox()
			t.dnsManager.ConfigureDns(`revert`)
			err := t.dnsManager.SetupDns(`revert`)
			if err != nil {
				t.l.Fatalf("don't manager dns: [%w]\n", err)
			}
			t.core.DisableDisconnectButton()
			err = t.core.TraySetImageDisconnect()
			if err != nil {
				t.l.Info("don't set image: [%w]\n", err)
			}
			break
		}

		if strings.Contains(text, "Connected via") {
			err := t.core.TraySetImageConnect()
			if err != nil {
				t.l.Info("don't set image: [%w]\n", err)
			}
			continue
		}

		if iterSkip {
			continue
		}

		r := `\w+[0-9]\s+opened|\w+\s+opened`
		if regexp.MustCompile(r).MatchString(text) {

			r := regexp.MustCompile(`(?P<int>\w+[0-9]|\w+)`)
			iface := r.FindStringSubmatch(text)

			t.core.SetPhyseInterface(iface[0])
			t.dnsManager.ConfigureDns(`setup`)
			err := t.dnsManager.SetupDns(`setup`)
			if err != nil {
				t.l.Fatalf("setup dns err [%w]\n", err)
			}
			err = t.core.TraySetImageConnect()
			if err != nil {
				t.l.Info("don't set image: [%w]\n", err)
			}
			iterSkip = true
			continue

		} else {
			err := t.core.CaseFlickeringIcon()
			if err != nil {
				t.l.Fatalf("case flick icon fatal err: [%w]", err)
			}
		}
	}
}
