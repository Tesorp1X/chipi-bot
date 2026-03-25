package application

import (
	"fmt"

	"github.com/Tesorp1X/chipi-bot/application/handlers"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	tele "gopkg.in/telebot.v4"
)

var allCommandsAndActions []commandWithStates = []commandWithStates{
	{
		command: "/start",
		states:  []fsm.State{fsm.AnyState},
	},
	{
		command: "/cancel",
		states:  []fsm.State{fsm.AnyState},
	},
}

func (app *Application) HandleAnyCommand(c tele.Context, state fsm.Context) error {
	cmd := c.Message().Text
	switch cmd {
	case "/start":
		if err := handlers.HandleStartCommand(c, state); err != nil {
			return fmt.Errorf("error in application.HandleAnyCommand(), command '%s': %v", cmd, err)
		}
	case "/cancel":
		if err := handlers.HandleCancelCommand(c, state); err != nil {
			return fmt.Errorf("error in application.HandleAnyCommand(), command '%s': %v", cmd, err)
		}
	default:
		return c.Send("unknown command")
	}

	return nil
}
