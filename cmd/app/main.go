package main

import (
	"log"

	embedfile "github.com/Vai3soh/goovpn"
	"github.com/Vai3soh/goovpn/internal/app"
	"github.com/caarlos0/env/v6"
)

type config struct {
	LogLvl string `env:"LVL_GOOVPN" envDefault:"debug"`
	PathDB string `env:"DB_GOOVPN"  envDefault:"goovpn.db"`
}

func main() {
	embd := embedfile.NewData()

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	app.Run_wails(
		*embd.Fs, embd.AppIcon, embd.ConnectIcon,
		embd.DisconnectIcon, cfg.LogLvl, cfg.PathDB,
	)
}
