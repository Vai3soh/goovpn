package openvpn

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Vai3soh/goovpn/entity"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logger interface {
	Fatal(...interface{})
	Debugf(string, ...interface{})
	Fatalf(string, ...interface{})
	Info(...interface{})
}

type sessionService interface {
	ReloadClient(context.Context)
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
	RunSession() error
	DestroyVpnClient()
	ExitSession()
	GetChanVpnLog() chan string
	SetPhyseInterface(iface string)
}

type DBService interface {
	SetNameBucket(name string)
	GetValueFromBucket(key string) error
	ReOpen() error
	Message() []entity.Message
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
	ctx         context.Context
	configsPath string
	core        core
	dnsManager  DnsManager
	pcore       ProfileCore
	l           Logger

	dbService DBService
	sessionService
}

func New(
	configsPath string, core core,
	pcore ProfileCore, dnsManager DnsManager,
	l Logger, db DBService, s sessionService,
) *TransportOvpnClient {

	return &TransportOvpnClient{
		configsPath:    configsPath,
		core:           core,
		dnsManager:     dnsManager,
		pcore:          pcore,
		l:              l,
		dbService:      db,
		sessionService: s,
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

func (t *TransportOvpnClient) caseRunSessionOpenvpn(profile string) {
	t.core.SetProfileBody(profile)
	t.initSession()
	err := t.core.RunSession()
	if err != nil {
		_, err := runtime.MessageDialog(t.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: err.Error(),
		})
		if err != nil {
			t.l.Fatal(err)
		}
	}
}

func (t *TransportOvpnClient) Context() context.Context {
	return t.ctx
}

func (t *TransportOvpnClient) SetContext(ctx context.Context) {
	t.ctx = ctx
}

func (t *TransportOvpnClient) Connect(cfgName string) {

	t.sessionService.ReloadClient(t.Context())
	err := t.dbService.ReOpen()
	if err != nil {
		t.l.Fatal(err)
	}
	t.dbService.SetNameBucket(`general_configure`)
	err = t.dbService.GetValueFromBucket(`reconn_count`)
	if err != nil {
		t.l.Fatal(err)
	}
	message := t.dbService.Message()
	v, _ := strconv.Atoi(message[0].Value)
	go t.readLogsFromChan(v)
	cfg := t.configsPath + cfgName
	profile := t.pcore.GetProfileFromCache(cfg)
	if profile.Body == "" {
		err := t.pcore.SaveProfileWithoutCfgFile(cfg)
		if err != nil {
			t.l.Fatalf("pcore failed: [%w]\n", err)
		}
		if !t.pcore.CheckUseCfgFile() {
			profile := t.pcore.GetProfileFromCache(cfg)
			os.Chdir(t.configsPath)
			t.caseRunSessionOpenvpn(profile.Body)

		} else {
			os.Chdir(t.configsPath)
			err := t.pcore.SaveProfileWithCfgFile(cfg)
			if err != nil {
				t.l.Fatalf("pcore failed: [%w]\n", err)
			}
			profile := t.pcore.GetProfileFromCache(cfg)
			t.caseRunSessionOpenvpn(profile.Body)
		}
	} else {
		t.caseRunSessionOpenvpn(profile.Body)
	}
}

func (t *TransportOvpnClient) Disconnect() {
	t.core.ExitSession()
	t.dnsManager.ConfigureDns(`revert`)
	err := t.dnsManager.SetupDns(`revert`)
	if err != nil {
		t.l.Fatalf("don't manager dns: [%w]\n", err)
	}
}

func (t *TransportOvpnClient) readLogsFromChan(countReccon int) {

	logChan := t.core.GetChanVpnLog()
	iterSkip := false
	toBreak := false
	count := 0
	for text := range logChan {
		runtime.EventsEmit(t.ctx, "rcv:read_log", text, count, countReccon)

		if toBreak {
			time.Sleep(time.Duration(10) * time.Millisecond)
			t.core.DestroyVpnClient()
			break
		}

		if strings.Contains(text, "Server poll timeout, trying next remote entry...") {
			count++
			if count == countReccon {
				t.core.ExitSession()
			}
		}

		if strings.Contains(text, `DNS Servers:`) {
			t.dnsManager.AddDnsAddrs(text)
		}

		if strings.Contains(text, "UNKNOWN/UNSUPPORTED OPTIONS") {
			toBreak = true
		}

		if strings.Contains(text, `event name: DISCONNECTED`) {

			if !iterSkip {
				t.dnsManager.ConfigureDns(`revert`)
				err := t.dnsManager.SetupDns(`revert`)
				if err != nil {
					t.l.Fatalf("don't manager dns: [%w]\n", err)
				}
			}
			toBreak = true
		}

		if strings.Contains(text, `event name: CONNECTED`) {
			continue
		}

		if strings.Contains(text, `event name: RECONNECTING`) {

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

			iterSkip = true
			continue

		} else {

		}
	}
}
