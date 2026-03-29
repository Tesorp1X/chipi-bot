package prompts

import (
	"fmt"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Send this in case some button or command exists, but not yet implemented.
func SendNotImplementedMessage(c tele.Context, state fsm.Context) error {
	if err := c.Send("Эта функция еще в разработке."); err != nil {
		return fmt.Errorf(
			"error in prompts.SendNotImplementedMessage(): failed to send a 'not implemented' message (%v)",
			err,
		)
	}

	return nil
}
