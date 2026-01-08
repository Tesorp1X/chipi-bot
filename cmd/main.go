package main

import (
	"flag"

	"github.com/Tesorp1X/chipi-bot/config"
)

var debug = flag.Bool("debug", false, "log debug info")

func main() {
	cfg, err := config.InitConfig(*debug)
	if err != nil {
		panic(err)
	}
}
