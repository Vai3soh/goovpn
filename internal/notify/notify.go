package notify

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Notify struct {
	ctx         context.Context
	connIcon    []byte
	disconnIcon []byte
}

func NewNotify(c, d []byte) *Notify {
	return &Notify{
		connIcon:    c,
		disconnIcon: d,
	}
}

func (n *Notify) SetContext(ctx context.Context) {
	n.ctx = ctx
}

func (n *Notify) DisconnectNotify() error {
	err := runtime.SendNotification(n.ctx, runtime.NotificationOptions{
		AppID:   "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\msinfo32.exe",
		AppIcon: n.disconnIcon,
		Title:   "Goovpn connection",
		Message: "Openvpn connection is lost",
		Timeout: 2 * time.Second,
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *Notify) ConnectNotify() error {
	err := runtime.SendNotification(n.ctx, runtime.NotificationOptions{
		AppID:   "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\msinfo32.exe",
		AppIcon: n.connIcon,
		Title:   "Goovpn connection",
		Message: "Openvpn connection established",
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	return nil
}
