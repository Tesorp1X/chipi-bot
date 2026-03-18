package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/db"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/fsmopt"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/v2/dispatcher"

	tele "gopkg.in/telebot.v4"
)

type Application struct {
	tgBot      *tele.Bot
	botGroup   *tele.Group
	dispatcher *dispatcher.Dispatcher
	fsmManager *fsm.Manager
	conf       *config.Config
	dbService  *db.DBService
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

func initBot(apiKey string, debug bool) (*tele.Bot, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:     apiKey,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose:   debug,
		ParseMode: tele.ModeHTML,
	})

	if err != nil {
		return nil, fmt.Errorf(
			"error in initBot(): couldn't init a new bot: %v",
			err,
		)
	}

	return b, nil
}

func NewApplication(conf *config.Config) (*Application, error) {
	b, err := initBot(conf.ApiKey, conf.VerboseDebug)
	if err != nil {
		return nil, fmt.Errorf(
			"error in NewApplication: couldn't initialize a new bot: %v",
			err,
		)
	}
	group := b.Group()
	d := dispatcher.NewDispatcher(group)
	manager := fsm.New(memory.NewStorage())

	dbs, err := db.MakeNewDBService(conf)
	if err != nil {
		return nil, fmt.Errorf(
			"error in application.NewApplication: couldn't initialize DB-Service: %v",
			err,
		)
	}

	app := &Application{
		tgBot:      b,
		botGroup:   group,
		dispatcher: d,
		fsmManager: manager,
		conf:       conf,
		dbService:  dbs,
	}

	app.registerCommands(allCommandsAndActions...)
	app.registerActions(allAllowedActionsAndStates...)

	return app, nil
}

type Command string

func (c *Command) String() string {
	return string(*c)
}

//type HandlerFunc = func(tele.Context, fsm.Context) error

type commandWithStates struct {
	command Command
	states  []fsm.State
}

func (app *Application) registerCommands(commandsWithHandlers ...commandWithStates) {
	for _, chs := range commandsWithHandlers {
		app.dispatcher.Dispatch(
			app.fsmManager.New(
				fsmopt.On(chs.command.String()),
				fsmopt.OnStates(chs.states...),
				fsmopt.Do(app.HandleAnyCommand),
			),
		)
	}
}

type Action string

func (a *Action) String() string {
	return string(*a)
}

type actionsWithStates struct {
	action Action
	states []fsm.State
}

func (app *Application) registerActions(actionsWithHandlers ...actionsWithStates) {
	for _, ahs := range actionsWithHandlers {
		app.fsmManager.Bind(
			app.dispatcher,
			fsmopt.On(ahs.action.String()),
			fsmopt.OnStates(ahs.states...),
			fsmopt.Do(app.HandleAnyAction),
		)
	}
}
