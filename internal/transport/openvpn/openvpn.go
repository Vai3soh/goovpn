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
	SetSessionCread(u, p string) error
	GetOvpnAuthPathFileName() string
	SetPathToFile(path string)
	ReadFile() ([]byte, error)
	SetBodyToCfg(path string)
	GetVpnCread() (string, string)
	CheckOvpnUseAuthUserPass() bool
	SetProfileBody(profileBody string)
	RunSession(ctx context.Context) error
	DestroyVpnClient()
	ExitSession()
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
	core core, pcore ProfileCore, dnsManager DnsManager, l Logger,
) *TransportOvpnClient {

	return &TransportOvpnClient{
		configsPath: configsPath,
		stopTimeout: stopTimeout,
		core:        core,
		dnsManager:  dnsManager,
		pcore:       pcore,
		l:           l,
	}
}

func (t *TransportOvpnClient) setVpnCread(username, password string) error {
	return t.core.SetSessionCread(username, password)
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
	err := t.setVpnCread(username, password)
	if err != nil {
		t.l.Fatal(err)
	}
}

func (t *TransportOvpnClient) caseRunSessionOpenvpn(ctx context.Context, profile string) {
	t.core.SetProfileBody(profile)
	t.initSession()
	err := t.core.RunSession(ctx)
	if err != nil {
		t.l.Fatal(err)
	}
}

func (t *TransportOvpnClient) Connect(ctx context.Context, countReccon int) func() error {

	f := func() error {
		go t.readLogsFromChan(countReccon)
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
			} else {
				os.Chdir(t.configsPath)
				err := t.pcore.SaveProfileWithCfgFile(cfg)
				if err != nil {
					return err
				}
				profile := t.pcore.GetProfileFromCache(cfg)
				t.caseRunSessionOpenvpn(ctx, profile.Body)
			}
		} else {
			t.caseRunSessionOpenvpn(ctx, profile.Body)
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

func (t *TransportOvpnClient) readLogsFromChan(countReccon int) {

	logChan := t.core.GetChanVpnLog()
	iterSkip := false
	toBreak := false
	count := 0
	for text := range logChan {
		t.core.CaseSetLogsInTextWidget(text)
		if toBreak {
			t.core.DestroyVpnClient()
			break
		}

		if strings.Contains(text, `Server poll timeout, trying next remote entry...`) {
			count++
			if count == countReccon {
				t.core.TurnOnConfigsBox()
				t.core.EnableConnectButton()
				toBreak = true
			}
		}

		if strings.Contains(text, `event name: CONNECTING`) {
			t.core.EnableDisconnectButton()
		}

		if strings.Contains(text, `DNS Servers:`) {
			t.dnsManager.AddDnsAddrs(text)
		}

		if strings.Contains(text, `event name: DISCONNECTED`) {
			t.core.EnableConnectButton()
			t.core.TurnOnConfigsBox()
			if !iterSkip {
				t.dnsManager.ConfigureDns(`revert`)
				err := t.dnsManager.SetupDns(`revert`)
				if err != nil {
					t.l.Fatalf("don't manager dns: [%w]\n", err)
				}
			}
			t.core.DisableDisconnectButton()
			err := t.core.TraySetImageDisconnect()
			if err != nil {
				t.l.Info("don't set image: [%w]\n", err)
			}
			toBreak = true
		}

		if strings.Contains(text, `event name: CONNECTED`) {
			err := t.core.TraySetImageConnect()
			if err != nil {
				t.l.Info("don't set image: [%w]\n", err)
			}
			continue
		}

		if strings.Contains(text, `event name: RECONNECTING`) {
			err := t.core.CaseFlickeringIcon()
			if err != nil {
				t.l.Fatalf("case flick icon fatal err: [%w]", err)
			}
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
