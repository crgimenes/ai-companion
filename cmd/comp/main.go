package main

import (
	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/json"
)

type Config struct {
	Prompt string `json:"prompt" cfg:"prompt" cfgDefault:"-"`
}

func main() {
	var cfg Config
	config.File = "config.json"
	config.Parse(&cfg)

	println(cfg.Prompt)
}
