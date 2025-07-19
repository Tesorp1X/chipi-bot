package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Tesorp1X/chipi-bot/handlers"
	"github.com/Tesorp1X/chipi-bot/middlewares"
	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/Tesorp1X/chipi-bot/util"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/v2/dispatcher"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

var debug = flag.Bool("debug", false, "log debug info")

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Couldn't load .env")
		panic(err)
	}

	pref := tele.Settings{
		Token:     os.Getenv("API_KEY"),
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose:   *debug,
		ParseMode: tele.ModeHTML,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	adminsList := util.ExtractAdminsIDs(os.Getenv("ADMINS"))
	g := bot.Group()
	g.Use(middleware.Whitelist(adminsList...))
	m := fsm.New(memory.NewStorage())

	// Bind to bot group for call before filters.
	// WrapContext adds fsm context into telebot context
	// It helps make less allocations
	g.Use(m.WrapContext)

	g.Use(middlewares.AutoSessionAssigner)

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
			fsmopt.On("/current"),         // set endpoint
			fsmopt.OnStates(fsm.AnyState), // set state filter
			fsmopt.Do(handlers.ShowCurrentTotalCommand),
		),
	)

	dp.Dispatch(
		m.New(
			fsmopt.On("/finish"),          // set endpoint
			fsmopt.OnStates(fsm.AnyState), // set state filter
			fsmopt.Do(handlers.FinishSessionCommand),
		),
	)

	dp.Dispatch(
		m.New(
			fsmopt.On("/show"),            // set endpoint
			fsmopt.OnStates(fsm.AnyState), // set state filter
			fsmopt.Do(handlers.ShowCommand),
		),
	)

	dp.Dispatch(
		m.New(
			fsmopt.On("/halo"),
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
		fsmopt.Do(handlers.CheckNameResponseHandler),
	)

	m.Bind(
		dp,
		fsmopt.On(tele.OnText),
		fsmopt.OnStates(models.StateWaitForNewCheckName),
		fsmopt.Do(handlers.NewCheckNameResponseHandler),
	)

	m.Bind(
		dp,
		fsmopt.On(tele.OnCallback),
		fsmopt.OnStates(fsm.AnyState),
		fsmopt.Do(handlers.HandleCallbackAction),
	)

	m.Bind(
		dp,
		fsmopt.On(tele.OnText),
		fsmopt.OnStates(models.StateWaitForItemName),
		fsmopt.Do(handlers.ItemNameResponseHandler),
	)

	m.Bind(
		dp,
		fsmopt.On(tele.OnText),
		fsmopt.OnStates(models.StateWaitForItemPrice),
		fsmopt.Do(handlers.ItemPriceResponseHandler),
	)

	bot.Start()
}
