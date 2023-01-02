package gui

import (
	"context"
	"fmt"
	rt "runtime"

	"github.com/Vai3soh/goovpn/entity"
	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Gui struct {
	ctx              context.Context
	fileManager      FilesManager
	sessionService   SessionService
	dbService        DbService
	logService       LogService
	transportService TransportService
}

type TransportService interface {
	Disconnect()
}

type FilesManager interface {
	FilesInDir(path string) ([]string, error)
	ChangeWorkingDir(path string) error
}

type SessionService interface {
	SetConnTimeout(t int)
	SetProxyAllowCleartextAuth(b bool)
	SetAllowUnusedAddrFamilies(s string)
	SetClockTickMS(uint)
	SetProxyHost(s string)
	SetExternalPkiAlias(s string)
	SetGremlinConfig(string)
	SetGuiVersion(string)
	SetHwAddrOverride(string)
	SetPlatformVersion(string)
	SetPortOverride(string)
	SetPrivateKeyPassword(string)
	SetProtoOverride(string)
	SetProtoVersionOverride(int)
	SetProxyUsername(string)
	SetProxyPassword(string)
	SetProxyPort(string)
	SetServerOverride(string)
	SetSsoMethods(string)
	SetTlsCertProfileOverride(string)
	SetTlsCipherList(string)
	SetTlsCiphersuitesList(string)
	SetTlsVersionMinOverride(string)
	SetDefaultKeyDirection(int)
	SetTunPersist(bool)
	SetEnableLegacyAlgorithms(bool)
	SetEnableNonPreferredDCAlgorithms(bool)
	SetDisableClientCert(bool)
	SetRetryOnAuthFailed(bool)
	SetAllowLocalDnsResolvers(bool)
	SetAllowLocalLanAccess(bool)
	SetAltProxy(bool)
	SetAutologinSessions(bool)
	SetDco(bool)
	SetEcho(bool)
	SetGenerate_tun_builder_capture_event(bool)
	SetGoogleDnsFallback(bool)
	SetInfo(bool)
	SetSynchronousDnsLookup(bool)
	SetWintun(bool)

	SetSslDebugLevel(int)
	SetCompressionMode(string)
}

type LogService interface {
	Fatal(...interface{})
	Debugf(string, ...interface{})
	Fatalf(string, ...interface{})
	Info(...interface{})
}

type DbService interface {
	ReOpen() error
	CreateBucket(name string) error
	Store(key, value string) error
	DeleteKey(key string) error
	StoreBulk(result []entity.Message) error
	GetValueFromBucket(key string) error
	GetAllValue()
	Message() []entity.Message
	SetNameBucket(nameBucket string)
	BucketIsCreate() bool
}

func NewGui(
	f FilesManager, s SessionService,
	d DbService, l LogService, t TransportService) *Gui {
	return &Gui{
		sessionService:   s,
		fileManager:      f,
		dbService:        d,
		logService:       l,
		transportService: t,
	}
}

func (g *Gui) Startup(ctx context.Context) {
	g.ctx = ctx
}

func (g *Gui) DomReady(ctx context.Context) {

}

func (g *Gui) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

func (g *Gui) Shutdown(ctx context.Context) {
	g.transportService.Disconnect()
}

func (g *Gui) Run() {

}

func (g *Gui) Stoped() {

}

func (g *Gui) IsWindows() bool {
	return rt.GOOS == `windows`
}

func (g *Gui) GetData(name string) []entity.Message {
	g.dbMustOpen()
	g.dbService.SetNameBucket(name)
	g.dbService.GetAllValue()
	return g.dbService.Message()
}

func (g *Gui) SaveToFrontendParams() {

	data := g.GetData(`general_configure`)
	runtime.EventsEmit(g.ctx, "rcv:save_from_db_general_configure", data)

	data = g.GetData(`general_openvpn_library`)
	err := g.setSessionParamsGeneralOvpnLib(data)
	if err != nil {
		g.logService.Fatal(err)
	}
	runtime.EventsEmit(g.ctx, "rcv:save_from_db_checkbox", data)

	data = g.GetData(`ssl_cmp`)
	g.setSessionParamsDebugSslAndCompression(data)

	runtime.EventsEmit(g.ctx, "rcv:save_from_db_select", data)

	data = g.GetData(`other_options`)
	err = g.setSessionParamsOtherOpt(data)
	if err != nil {
		g.logService.Fatal(err)
	}
	runtime.EventsEmit(g.ctx, "rcv:save_from_db_input", data)

}

func (g *Gui) GetConfigsAndChangeCWD() ([]string, error) {
	g.dbMustOpen()
	g.dbService.SetNameBucket(`general_configure`)
	err := g.dbService.GetValueFromBucket(`config_dir_path`)
	if err != nil {
		return nil, err
	}
	path := g.dbService.Message()
	files, err := g.fileManager.FilesInDir(path[0].Value)
	if err != nil {
		return nil, err
	}
	err = g.fileManager.ChangeWorkingDir(path[0].Value)
	if err != nil {
		return nil, err
	}
	return files, err
}

func (g *Gui) dbMustOpen() {
	err := g.dbService.ReOpen()
	if err != nil {
		g.logService.Fatalf("don't reopen db [%s]\n", err)
	}
}

func (g *Gui) DeleteData(data map[string]interface{}, bucketName string) {
	result, err := g.DecodeData(data)
	if err != nil {
		g.logService.Fatalf("don't decode data [%s]\n", err)
	}
	g.dbMustOpen()
	err = g.dbService.CreateBucket(bucketName)
	if err != nil {
		g.logService.Fatalf("don't create storage [%s]\n", err)
	}
	g.dbMustOpen()
	err = g.dbService.DeleteKey(result.AtrId)
	if err != nil {
		g.logService.Fatalf("don't create storage [%s]\n", err)
	}
	g.setSessionFromAnyBucket(result, bucketName)
}

func (g *Gui) setSessionFromAnyBucket(result *entity.Message, bucketName string) {
	inObj := []entity.Message{*result}
	if bucketName == `general_openvpn_library` {
		g.setSessionParamsGeneralOvpnLib(inObj)
	} else if bucketName == `ssl_cmp` {
		g.setSessionParamsDebugSslAndCompression(inObj)
	} else if bucketName == `other_options` {
		g.setSessionParamsOtherOpt(inObj)
	}
}

func (g *Gui) SaveData(data map[string]interface{}, bucketName string) {

	result, err := g.DecodeData(data)
	if err != nil {
		g.logService.Fatalf("don't decode data [%s]\n", err)
	}
	g.dbMustOpen()
	err = g.dbService.CreateBucket(bucketName)
	if err != nil {
		g.logService.Fatalf("don't create storage [%s]\n", err)
	}
	g.dbMustOpen()
	err = g.dbService.Store(result.AtrId, result.Value)
	if err != nil {
		g.logService.Fatalf("don't save value [%s] [%s]\n", result.Value, err)
	}
	g.setSessionFromAnyBucket(result, bucketName)
}

func (g *Gui) DecodeData(data map[string]interface{}) (*entity.Message, error) {
	var result entity.Message
	err := mapstructure.Decode(data, &result)
	if err != nil {
		return nil, fmt.Errorf("decode failed %s", err)
	}
	return &result, nil
}
