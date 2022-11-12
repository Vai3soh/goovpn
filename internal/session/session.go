package session

import (
	"context"
	"fmt"

	"github.com/Vai3soh/ovpncli"
)

type OpenvpnClient struct {
	ovpncli.Client
	ovpncli.ClientAPI_Config
	OverwriteClient
}

type OpenvpnClientCred struct {
	ovpncli.ClientAPI_ProvideCreds
}

type OverwriteClient struct {
	ovpncli.ClientAPI_OpenVPNClient
	vpnLog chan string
}

func NewOverwriteClient() *OverwriteClient {
	return &OverwriteClient{
		vpnLog: make(chan string),
	}
}

func (ocl *OverwriteClient) Log(li ovpncli.ClientAPI_LogInfo) {
	ocl.vpnLog <- li.GetText()

}

func (ocl *OverwriteClient) ChanVpnLog() chan string {
	return ocl.vpnLog
}

func (ocl *OverwriteClient) CloseChanVpnLog() {
	close(ocl.vpnLog)
}

func (ocl *OverwriteClient) Event(ev ovpncli.ClientAPI_Event) {
	ocl.vpnLog <- fmt.Sprintf("event name: %s", ev.GetName())
	ocl.vpnLog <- fmt.Sprintf("event info: %s", ev.GetInfo())
}

func (ocl *OverwriteClient) Remote_override_enabled() {

}

func (ocl *OverwriteClient) Socket_protect() {

}

type Option func(*OpenvpnClient)

func NewOpenvpnClientCred(optsCred ...ovpncli.OptionCred) *OpenvpnClientCred {
	op := &OpenvpnClientCred{
		ClientAPI_ProvideCreds: ovpncli.NewClientCreds(optsCred...),
	}
	return op
}

func NewOpenvpnClient(opts ...ovpncli.Option) *OpenvpnClient {
	ocl := &OverwriteClient{vpnLog: make(chan string)}
	op := &OpenvpnClient{
		Client:           ovpncli.NewClient(ocl),
		ClientAPI_Config: ovpncli.NewClientConfig(opts...),
		OverwriteClient:  *ocl,
	}
	return op
}

func (op *OpenvpnClient) SetClient(c ovpncli.Client) {
	op.Client = c
}

func (op *OpenvpnClient) DestroyClient() {
	ovpncli.DeleteClient(op.Client)
}

func (op *OpenvpnClient) GetOverwriteClient() *OverwriteClient {
	return &op.OverwriteClient
}

func (op *OpenvpnClient) StartSession(ctx context.Context) error {
	ev := op.Eval_config(op.ClientAPI_Config)
	if ev.GetError() {
		return fmt.Errorf("config eval failed [%s]", ev.GetMessage())
	}
	op.StartConnection(ctx)
	return nil
}

func (op *OpenvpnClient) StopSession() {
	op.StopConnection()
}

func (op *OpenvpnClient) SetConfig(config string) {
	op.ClientAPI_Config.SetContent(config)
}

func (op *OpenvpnClient) SetCread(u, p string) error {
	creds := NewOpenvpnClientCred(ovpncli.WithPassword(p), ovpncli.WithUsername(u))
	status := op.Provide_creds(creds)
	if status.GetError() {
		return fmt.Errorf("provide cred failed [%s]", status.GetMessage())
	}
	return nil
}
