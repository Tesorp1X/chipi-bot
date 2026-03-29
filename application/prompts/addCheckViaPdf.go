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
func SendCheckVerificationMessage(c tele.Context, state fsm.Context) error {
	check, err := storageHelpers.GetCheck(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to retrieve a check (%v)",
			err,
		)
	}

	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to retrieve an items-list (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateWaitingForCheckConfirmationUnsaved, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to change a state to a '%s' (%v)",
			static.StateWaitingForCheckConfirmationUnsaved,
			err,
		)
	}

	if err := c.EditOrSend(responses.GetVerificationFinalStepResponse(check, items)); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckVerificationMessage(): failed to send a message (%v)",
			err,
		)
	}

	return nil
}

// Prepares and sends message with check edit text and buttons and sets state to a 'StateEditingCheck'.
func SendEditCheckMessage(c tele.Context, state fsm.Context) error {
	check, err := storageHelpers.GetCheck(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendEditCheckMessage(): failed to retrieve a check (%v)",
			err,
		)
	}

	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendEditCheckMessage(): failed to retrieve an items-list (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateEditingCheckUnsaved, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendEditCheckMessage(): failed to change a state to a '%s' (%v)",
			static.StateEditingCheckUnsaved,
			err,
		)
	}

	verificationText, _ := responses.GetVerificationFinalStepResponse(check, items)
	if err := c.EditOrSend(responses.GetEditCheckMessage(verificationText)); err != nil {
		return fmt.Errorf(
			"error in prompts.SendEditCheckMessage(): failed to send a message (%v)",
			err,
		)
	}

	return nil
}

const (
	// In case of calling from AddCheck sequence, before final stage.
	FromAddCheck = iota
	// In case of calling from EditCheck in final stage of AddCheck scenarios.
	FromEditCheckFinal
	// In case of calling from EditCheck scenarios.
	FromEditCheck
)

// Sends a check ownership message with or without go-back button, depending on a cameFrom argument.
func SendCheckOwnershipMessage(cameFrom int, c tele.Context, state fsm.Context) error {
	var withGoBackButton bool
	switch cameFrom {
	case FromAddCheck:
		withGoBackButton = false
	case FromEditCheckFinal:
		state.Update(context.Background(), static.IS_FROM_FINAL_STAGE, true)
		withGoBackButton = true
	default:
		return fmt.Errorf(
			"error in prompts.SendCheckOwnerMessage(): invalid cameFrom value (%d)",
			cameFrom,
		)
	}

	if err := storageHelpers.SetState(static.StateWaitForCheckOwner, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckOwnerMessage(): failed change state to a '%s' (%v)",
			static.StateWaitForCheckOwner,
			err,
		)
	}

	if err := c.EditOrSend(responses.GetAskForCheckOwnershipQuestion(withGoBackButton)); err != nil {
		return fmt.Errorf(
			"error in prompts.SendCheckOwnerMessage(): failed to send check ownership message (%v)",
			err,
		)
	}

	return nil
}

// Sends a 'new check-name question' message.
// Depending on a cameFrom argument, a different behavior is going to be applied:
// - if FromAddCheck, then state will be set to StateWaitForCheckName and response used from GenerateNameVerificationResponse;
// - if FromEditCheckFinal, then state will be set to StateWaitForNewCheckName and response used from GetAskForNewCheckNameResponse.
func SendChangeCheckNameMessage(cameFrom int, c tele.Context, state fsm.Context) error {
	check, err := storageHelpers.GetCheck(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendNewCheckNameQuestionMessage(): failed to retrieve a check (%v)",
			err,
		)
	}

	var newState fsm.State

	switch cameFrom {
	case FromAddCheck:
		newState = static.StateWaitForCheckName
	case FromEditCheckFinal:
		newState = static.StateWaitForNewCheckNameUnsaved
	default:
		return fmt.Errorf(
			"error in prompts.SendNewCheckNameQuestionMessage(): invalid cameFrom value (%d)",
			cameFrom,
		)
	}

	text, kb := responses.GetAskForNewCheckNameResponse(check.Name)

	if err := storageHelpers.SetState(newState, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendNewCheckNameQuestionMessage(): failed change state to a '%s' (%v)",
			newState,
			err,
		)
	}

	if err := c.EditOrSend(text, kb); err != nil {
		return fmt.Errorf(
			"error in prompts.SendNewCheckNameQuestionMessage(): failed to send an edit-check message (%v)",
			err,
		)
	}

	return nil
}

// Sends a 'check creation date question' message.
func SendChangeCreationDateMessage(c tele.Context, state fsm.Context) error {
	if err := storageHelpers.SetState(static.StateWaitForCheckCreationDateUnsaved, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendNewCheckNameQuestionMessage(): failed change state to a '%s' (%v)",
			static.StateWaitForCheckCreationDateUnsaved,
			err,
		)
	}

	if err := c.EditOrSend(responses.GetAskForNewCheckCreationDateQuestion()); err != nil {
		return fmt.Errorf(
			"error in prompts.SendNewCheckNameQuestionMessage(): failed to send an edit-check message (%v)",
			err,
		)
	}

	return nil
}

// Sends a message with an item. Depending on a cameFrom argument,
// a different behavior is going to be applied:
// - if FromAddCheck, then state is set to StateShowingAnItem and response is from GetItemVerificationResponse;
// - if FromEditCheckFinal, then state is set to StateShowingAnItem and response is from GetShowItemForEditResponse;
func SendShowItemsMessage(cameFrom int, c tele.Context, state fsm.Context) error {
	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendShowItemsMessage(): failed to retrieve an items-list (%v)",
			err,
		)
	}

	currentIndex, err := storageHelpers.GetCurrentIndex(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in prompts.SendShowItemsMessage(): failed to retrieve a current index (%v)",
			err,
		)
	}

	var newState fsm.State
	var text string
	var kb *tele.ReplyMarkup

	switch cameFrom {
	case FromAddCheck:
		newState = static.StateShowingAnItem
		text, kb = responses.GetItemVerificationResponse(items[currentIndex], currentIndex, len(items))
	case FromEditCheckFinal:
		newState = static.StateShowingAnItem
		text, kb = responses.GetShowItemForEditResponse(items[currentIndex], currentIndex, len(items))
	default:
		return fmt.Errorf(
			"error in prompts.SendShowItemsMessage(): invalid cameFrom value (%d)",
			cameFrom,
		)
	}

	if err := storageHelpers.SetState(newState, c, state); err != nil {
		return fmt.Errorf(
			"error in prompts.SendShowItemsMessage(): failed change state to a '%s' (%v)",
			newState,
			err,
		)
	}

	if err := c.EditOrSend(text, kb); err != nil {
		return fmt.Errorf(
			"error in prompts.SendShowItemsMessage(): failed to send a 'show item' message (%v)",
			err,
		)
	}

	return nil
}
