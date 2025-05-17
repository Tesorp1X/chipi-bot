package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Tesorp1X/chipi-bot/handlers"
	"github.com/Tesorp1X/chipi-bot/models"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/v2/dispatcher"

	//"github.com/vitaliy-ukiru/telebot-filter/v2/routing"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

var debug = flag.Bool("debug", false, "log debug info")

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Couldn't load .env")
		panic(err)
	}

	pref := tele.Settings{
		Token:   os.Getenv("API_KEY"),
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: *debug,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	g := bot.Group()
	m := fsm.New(memory.NewStorage())

	// Bind to bot group for call before filters.
	// WrapContext adds fsm context into telebot context
	// It helps make less allocations
	g.Use(m.WrapContext)

	dp := dispatcher.NewDispatcher(g)

	dp.Dispatch(
		m.New(
			fsmopt.On("/cancel"),          // set endpoint
			fsmopt.OnStates(fsm.AnyState), // set state filter
			fsmopt.Do(handlers.CancelHandler),
		),
	)

	dp.Dispatch(
		m.New(
			fsmopt.On("/hello"),
			fsmopt.OnStates(fsm.AnyState),
			fsmopt.Do(handlers.HelloHandler),
		),
	)

	dp.Dispatch(
		m.New(
			fsmopt.On("/newcheck"),
			fsmopt.OnStates(fsm.AnyState),
			fsmopt.Do(handlers.NewCheckHandler),
		),
	)

	m.Bind(
		dp,
		fsmopt.On(tele.OnText),
		fsmopt.OnStates(models.StateWaitForCheckName),
		fsmopt.Do(handlers.ChecknameResponseHandler),
	)

	m.Bind(
		dp,
		fsmopt.On(tele.OnCallback),
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.Do(handlers.HandleCallbackAction),
	)

	bot.Start()
	log.Println("Bot is operational!")
}
