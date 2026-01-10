package application

import (
	"fmt"

	"github.com/Tesorp1X/chipi-bot/application/handlers"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	tele "gopkg.in/telebot.v4"
)

var allAllowedActionsAndStates []actionsWithStates = []actionsWithStates{
	{
		action: tele.OnDocument,

		states: []fsm.State{fsm.AnyState},
	},
}

// default function
func (app *Application) HandleAnyAction(c tele.Context, state fsm.Context) error {
	switch {
	case c.Message().Document != nil:
		if err := handlers.OnDocumentActionHandler(app.conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyAction(), action 'OnDocument': %v", err)
		}
	default:
		return c.Send("unknown action")
	}

	return nil
}
