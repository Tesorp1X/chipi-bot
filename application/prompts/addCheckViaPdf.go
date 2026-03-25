package prompts

import (
	"fmt"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Prepares and sends message with check verification and sets state to a 'StateWaitingForCheckConfirmation'.
func SendCheckVerificationMessage(check *static.Check, items []*static.Item, c tele.Context, state fsm.Context) error {
	if err := c.Send(responses.GetVerificationFinalStepResponse(check, items)); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): couldn't send a message (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateWaitingForCheckConfirmation, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): couldn't change a state (%v)",
			err,
		)
	}

	return nil
}
