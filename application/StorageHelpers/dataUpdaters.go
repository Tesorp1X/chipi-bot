package storageHelpers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Updates a CURRENT_INDEX_ITEMS value in a context-storage.
// If error occurs, it will bbe wrapped and returned.
func UpdateCurrentItemsIndex(newCurrentIndex int, c tele.Context, state fsm.Context) error {
	if err := state.Update(context.Background(), static.CURRENT_INDEX_ITEMS, newCurrentIndex); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in UpdateCurrentIndex(): couldn't save current index in state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	return nil
}

// Updates a CHECK value in a context-storage.
// If error occurs, it will bbe wrapped and returned.
func UpdateCheck(check *static.Check, c tele.Context, state fsm.Context) error {
	if err := state.Update(context.Background(), static.CHECK, check); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in UpdateCheck(): couldn't save check in state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	return nil
}
