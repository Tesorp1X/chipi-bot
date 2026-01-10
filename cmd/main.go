package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/Tesorp1X/chipi-bot/application"
	"github.com/Tesorp1X/chipi-bot/config"
)

var debug = flag.Bool("debug", false, "log debug info")

func main() {
	cfg, err := config.InitConfig(*debug)
	if err != nil {
		panic(err)
	}

	app, err := application.NewApplication(cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	stopped := app.RunBot(ctx)
	<-stopped

}
