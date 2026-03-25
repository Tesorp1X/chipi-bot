package prompts

import (
	"context"
	"fmt"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Prepares and sends message with check-verification text and buttons and sets state to a 'StateWaitingForCheckConfirmation'.
func SendCheckVerificationMessage(check *static.Check, items []*static.Item, c tele.Context, state fsm.Context) error {
	if err := c.Send(responses.GetVerificationFinalStepResponse(check, items)); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to send a message (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateWaitingForCheckConfirmation, c, state); err != nil {
		currentState, _ := state.State(context.Background())
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to change a state to a '%s' (%v)",
			currentState,
			err,
		)
	}

	return nil
}

// Prepares and sends message with check edit text and buttons and sets state to a 'StateEditingCheck'.
func SendEditCheckMessage(check *static.Check, items []*static.Item, c tele.Context, state fsm.Context) error {
	verificationText, _ := responses.GetVerificationFinalStepResponse(check, items)
	if err := c.Send(responses.GetEditCheckMessage(verificationText)); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to send a message (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateEditingCheck, c, state); err != nil {
		currentState, _ := state.State(context.Background())
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to change a state to a '%s' (%v)",
			currentState,
			err,
		)
	}

	return nil
}
