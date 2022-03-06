package session

import (
	"time"

	"github.com/Vai3soh/goovpn/internal/usecase"
	"github.com/Vai3soh/goovpn/pkg/openvpn3"
)

type Openvpn struct {
	openvpn3.Config
	openvpn3.Session
	usecase.UiInteractor
}

type Option func(*Openvpn)

func NewOpenvpn(opts ...Option) *Openvpn {
	op := &Openvpn{}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func WithConfig(config string) Option {
	return func(op *Openvpn) {
		op.ProfileContent = config
	}
}

func WithCompressionMode(mode string) Option {
	return func(op *Openvpn) {
		op.CompressionMode = mode
	}
}

func WithTimeout(timeout int) Option {
	return func(op *Openvpn) {
		op.ConnTimeout = timeout
	}
}

func WithDisableClientCert(b bool) Option {
	return func(op *Openvpn) {
		op.DisableClientCert = b
	}
}

func WithUi(ui usecase.UiInteractor) Option {
	return func(op *Openvpn) {
		op.UiInteractor = ui
	}
}

func (op *Openvpn) SetConfig(config string) {
	op.ProfileContent = config
}

func (op *Openvpn) SetSession(username, password string) {
	op.Session = *openvpn3.NewSession(
		op.Config,
		openvpn3.UserCredentials{
			Username: username,
			Password: password,
		},
		op.UiInteractor,
	)
}

func (op *Openvpn) StartSession() {
	op.Session.Start()
}

func (op *Openvpn) StopSession() {
	op.Session.Stop()
}

func (op *Openvpn) StopSessionWithTimeout(timeout int) {
	op.StopSession()
	time.Sleep(time.Duration(timeout) * time.Millisecond)
}

func (op *Openvpn) SessionIsClose() bool {
	return op.Session.IsClose()
}
