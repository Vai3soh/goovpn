package session

import (
	"context"

	"github.com/Vai3soh/goovpn/pkg/openvpn3"
)

type OpenvpnClient struct {
	cread *Cread
	*openvpn3.Config
	*openvpn3.Session
	ui interface{}
}

type Cread struct {
	Username string
	Password string
}

type Option func(*OpenvpnClient)

func NewOpenvpnClient(opts ...Option) *OpenvpnClient {
	op := &OpenvpnClient{
		cread:   &Cread{},
		Config:  &openvpn3.Config{},
		Session: &openvpn3.Session{},
	}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func WithConfig(config string) Option {
	return func(op *OpenvpnClient) {
		op.ProfileContent = config
	}
}

func WithCompressionMode(mode string) Option {
	return func(op *OpenvpnClient) {
		op.CompressionMode = mode
	}
}

func WithTimeout(timeout int) Option {
	return func(op *OpenvpnClient) {
		op.ConnTimeout = timeout
	}
}

func WithDisableClientCert(b bool) Option {
	return func(op *OpenvpnClient) {
		op.DisableClientCert = b
	}
}

func WithUi(callbacks interface{}) Option {
	return func(op *OpenvpnClient) {
		op.ui = callbacks
	}
}

func (op *OpenvpnClient) SetConfig(config string) {
	op.ProfileContent = config
}

func (op *OpenvpnClient) SetCread(user, pwd string) {
	op.cread = &Cread{Username: user, Password: pwd}
}

func (op *OpenvpnClient) SetSession() {

	op.Session = openvpn3.NewSession(
		*op.Config,
		openvpn3.UserCredentials{
			Username: op.cread.Username,
			Password: op.cread.Password,
		},
		op.ui,
	)
}

func (op *OpenvpnClient) StartSession(ctx context.Context) {
	op.Session.Start(ctx)

}

func (op *OpenvpnClient) StopSession() {
	op.Session.Stop()
}
