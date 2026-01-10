package callbacks

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func HandleAnyCallback(conf *config.Config, c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error in HandleAnyCallback(): couldn't receive users(%d) current state: %v",
			c.Sender().ID, err,
		)
	}

	callbackData := c.Callback().Data

	switch {
	case currentState == static.StateWaitForCheckName &&
		static.CallbackActionKeep.DataMatches(callbackData):
		if err := handleKeepCheckNameCallback(conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyCallback(), state 'StateWaitForCheckName', action 'CallbackActionKeep': %v", err)
		}
	default:
		// if callback query is old, remove inline buttons from that message
		c.Bot().EditReplyMarkup(c.Callback().Message, &tele.ReplyMarkup{})
		return c.Respond(&tele.CallbackResponse{Text: "error todo"})
	}
	return nil
}

func handleKeepCheckNameCallback(conf *config.Config, c tele.Context, state fsm.Context) error {
	return nil
}
