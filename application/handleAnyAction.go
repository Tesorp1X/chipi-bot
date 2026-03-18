package application

import (
	"fmt"

	"github.com/Tesorp1X/chipi-bot/application/callbacks"
	"github.com/Tesorp1X/chipi-bot/application/handlers"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	tele "gopkg.in/telebot.v4"
)

var allAllowedActionsAndStates []actionsWithStates = []actionsWithStates{
	{
		action: tele.OnDocument,
		// todo: change allowed states
		states: []fsm.State{fsm.AnyState},
	},
	{
		action: tele.OnCallback,
		// todo: change allowed states
		states: []fsm.State{fsm.AnyState},
	},
	{
		action: tele.OnText,
		// todo: change allowed states
		states: []fsm.State{fsm.AnyState},
	},
}

// default function
func (app *Application) HandleAnyAction(c tele.Context, state fsm.Context) error {
	switch {
	case c.Callback() != nil:
		if err := callbacks.HandleAnyCallback(app.conf, c, state); err != nil {
			return fmt.Errorf("error in application.HandleAnyAction(), action 'OnCallback': %v", err)
		}
	case c.Message().Document != nil:
		if err := handlers.OnDocumentActionHandler(app.conf, c, state); err != nil {
			return fmt.Errorf("error in application.HandleAnyAction(), action 'OnDocument': %v", err)
		}
	case c.Message().Text != "":
		if err := handlers.HandleAnyText(app.conf, c, state); err != nil {
			return fmt.Errorf("error in application.HandleAnyAction(), action 'OnText': %v", err)
		}
	default:
		return c.Send("unknown action")
	}

	return nil
}
