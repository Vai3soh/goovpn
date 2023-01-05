package embedfile

import (
	"embed"
)

var (
	//go:embed frontend/src
	assets embed.FS
	//go:embed  embedfile/assets/app.png
	icon []byte
	//go:embed  embedfile/assets/connecting.png
	connectIcon []byte
	//go:embed  embedfile/assets/disconnect.png
	disconnectIcon []byte
)

type Data struct {
	Fs             *embed.FS
	AppIcon        []byte
	ConnectIcon    []byte
	DisconnectIcon []byte
}

func NewData() *Data {
	return &Data{
		Fs:             &assets,
		AppIcon:        icon,
		ConnectIcon:    connectIcon,
		DisconnectIcon: disconnectIcon,
	}
}
