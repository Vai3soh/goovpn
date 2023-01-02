package config

import (
	"strconv"

	"github.com/Vai3soh/goovpn/entity"
)

type dbService interface {
	ReOpen() error
	SetNameBucket(nameBucket string)
	GetValueFromBucket(searchParam string) error
	Message() []entity.Message
}

type logService interface {
	Fatalf(msg string, args ...interface{})
}

type paramsDefault struct {
	dbService
	logService
	useSystemd               bool
	disableCert              bool
	tunPersist               bool
	legacyAlgorithms         bool
	nonPreferredDCAlgorithms bool
	sslDebug                 int
	compressionMode          string
	configsPath              string
	paramsOtherOpt
}

type paramsOtherOpt struct {
	connTimeout int
}

func NewParamsDefault(dbService dbService, logService logService) *paramsDefault {
	return &paramsDefault{
		dbService:  dbService,
		logService: logService,
	}
}

func (p *paramsDefault) UseSystemd() bool {
	return p.useSystemd
}

func (p *paramsDefault) DisableCert() bool {
	return p.disableCert
}

func (p *paramsDefault) TunPersist() bool {
	return p.tunPersist
}

func (p *paramsDefault) LegacyAlgorithms() bool {
	return p.legacyAlgorithms
}

func (p *paramsDefault) NonPreferredDCAlgorithms() bool {
	return p.nonPreferredDCAlgorithms
}

func (p *paramsDefault) SslDebug() int {
	return p.sslDebug
}

func (p *paramsDefault) CompressionMode() string {
	return p.compressionMode
}

func (p *paramsDefault) ConfigsPath() string {
	return p.configsPath
}

func (p *paramsDefault) ConnTimeout() int {
	return p.connTimeout
}

func (p *paramsDefault) SetUseSystemd(b bool) {
	p.useSystemd = b
}

func (p *paramsDefault) SetCompressionMode(c string) {
	p.compressionMode = c
}

func (p *paramsDefault) SetConnTimeout(t int) {
	p.connTimeout = t
}

func (p *paramsDefault) SetDisableCert(b bool) {
	p.disableCert = b
}

func (p *paramsDefault) SetTunPersist(b bool) {
	p.tunPersist = b
}

func (p *paramsDefault) SetLegacyAlgo(b bool) {
	p.legacyAlgorithms = b
}

func (p *paramsDefault) SetConfigsPath(path string) {
	p.configsPath = path
}

func (p *paramsDefault) SetNonPreferredDCAlgorithms(e bool) {
	p.nonPreferredDCAlgorithms = e
}

func (p *paramsDefault) SetSslDebug(i int) {
	p.sslDebug = i
}

func (p *paramsDefault) mustOpen() {
	err := p.dbService.ReOpen()
	if err != nil {
		p.logService.Fatalf("don't open [%s]\n", err)
	}
}

func (p *paramsDefault) getValue(nameBucket, searchParam string) *string {
	p.mustOpen()
	p.dbService.SetNameBucket(nameBucket)
	err := p.dbService.GetValueFromBucket(searchParam)
	if err != nil {

		return nil
	}
	return &p.dbService.Message()[0].Value
}

func getBool(v *string) bool {
	if v != nil {
		return *v != ""
	}
	return false
}

func (p *paramsDefault) GetParamIfStoreInDb() *paramsDefault {

	sslCmp := map[string]any{
		`#ssl`: p.SetSslDebug,
		`#cmp`: p.SetCompressionMode,
	}
	ovpnLib := map[string]any{
		`tun_persist`:         p.SetTunPersist,
		`disable_client_cert`: p.SetDisableCert,
		`legacy_algo`:         p.SetLegacyAlgo,
		`preferred_dc_algo`:   p.SetNonPreferredDCAlgorithms,
	}

	configure := map[string]any{
		`config_dir_path`: p.SetConfigsPath,
		`use_systemd`:     p.SetUseSystemd,
	}

	other := map[string]any{
		`with_conn_timeout`: p.SetConnTimeout,
	}

	m := map[string][]map[string]any{

		`general_openvpn_library`: {
			ovpnLib,
		},

		`general_configure`: {
			configure,
		},

		`ssl_cmp`: {
			sslCmp,
		},

		`other_options`: {
			other,
		},
	}

	for k, v := range m {
		for _, e := range v {
			for id, f := range e {
				w := p.getValue(k, id)
				if c, ok := f.(func(bool)); ok {
					c(getBool(w))
				} else if c, ok := f.(func(string)); ok {
					if w != nil {
						c(*w)
					}
				} else if c, ok := f.(func(int)); ok {
					if w != nil {
						t, _ := strconv.Atoi(*w)
						c(t)
					}
				}
			}
		}
	}
	return p
}
