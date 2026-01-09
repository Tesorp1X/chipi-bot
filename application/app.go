package application

import (
	"context"
	"fmt"
	"log"

	"github.com/Tesorp1X/chipi-bot/config"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/telebot-filter/dispatcher"

	tele "gopkg.in/telebot.v4"
)

type Application struct {
	tgBot      *tele.Bot
	botGroup   *tele.Group
	dispatcher *dispatcher.Dispatcher
	fsmManager *fsm.Manager
	conf       *config.Config
}

func (app *Application) RunBot(ctx context.Context) <-chan struct{} {
	closedCh := make(chan struct{})
	go app.tgBot.Start()
	fmt.Println("started: from RunBot (after bot.Start())")

	go func() {
		defer func() {
			closedCh <- struct{}{}
			close(closedCh)

		}()
		<-ctx.Done()
		app.tgBot.Stop()
		log.Println("stopped: from RunBot (after bot.Stop())")
	}()

	return closedCh
}
