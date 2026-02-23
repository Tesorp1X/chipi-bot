package storageHelpers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Retrieves a new check name from a message text,
// verifies it and saves a new name to a check object.
// If an error occurs, a nil is returned alongside an error itself.
func SetNewCheckNameFromMessage(c tele.Context, state fsm.Context) (*static.Check, error) {
	//get checkData obj from context
	check, err := GetCheck(c, state)
	if err != nil {

		return nil, fmt.Errorf(
			"error in setNewCheckName(): couldn't retrieve check from state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}
	// // todo: verify the message
	// if c.Message().Text == "" {
	// 	// ask again for the name
	// }
	//save new check name in it and put new version in context
	check.Name = c.Message().Text

	if err := state.Update(context.Background(), static.CHECK, check); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return nil, fmt.Errorf(
			"error in setNewCheckName(): couldn't save check in state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}
	return check, nil
}
