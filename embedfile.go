package embedfile

import (
	"embed"
)

var (
	//go:embed frontend/src
	assets embed.FS
	//go:embed  build/appicon.png
	icon []byte
)

type Data struct {
	Fs   *embed.FS
	Icon []byte
}

func NewData() *Data {
	return &Data{
		Fs:   &assets,
		Icon: icon,
	}
}
