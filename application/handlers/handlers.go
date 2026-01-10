package handlers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func HandleStartCommand(c tele.Context, state fsm.Context) error {
	return c.Send("hello")
}

func HandleAnyText(conf *config.Config, c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error in HandleAnyCallback(): couldn't receive users(%d) current state: %v",
			c.Sender().ID, err,
		)
	}

	switch currentState {
	case static.StateWaitForCheckName:
		if err := handleCheckName(conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyText(), state 'StateWaitForCheckName': %v", err)
		}
	}
	return nil
}

func handleCheckName(conf *config.Config, c tele.Context, state fsm.Context) error {
	//get checkData obj from context
	//save new check name in it and put new version in context
	//get first item and start verification of items
	return nil
}
