package main

import (
	"flag"
	"log"

	"github.com/Vai3soh/goovpn/config"
	"github.com/Vai3soh/goovpn/internal/app"
	"github.com/ilyakaznacheev/cleanenv"
)

var Config = flag.String(
	"config", "./config/config.yml",
	"path to configuration file",
)

func init() {
	flag.Parse()
}

func main() {
	var cfg config.Config

	err := cleanenv.ReadConfig(*Config, &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(&cfg)
}
