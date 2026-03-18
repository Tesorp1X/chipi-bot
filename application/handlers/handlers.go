package handlers

import (
	"context"
	"fmt"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"

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
			"error in handlers.HandleAnyCallback(): couldn't receive users(%d) current state: %v",
			c.Sender().ID, err,
		)
	}

	switch currentState {
	case static.StateWaitForCheckName:
		if err := handleCheckName(conf, c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForCheckName': %v", err)
		}
	case static.StateWaitForNewCheckName:
		if err := handleEditCheckName(conf, c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForNewCheckName': %v", err)
		}
	}
	return nil
}

func handleCheckName(conf *config.Config, c tele.Context, state fsm.Context) error {
	_, err := storageHelpers.SetNewCheckNameFromMessage(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't set new check name (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateWaitForCheckOwner, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't change a state (%v)",
			err,
		)
	}
	// prompt check ownership
	if sendErr := c.Send(responses.GetAskForCheckOwnershipQuestion()); sendErr != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't send a 'check-ownership'-message (%v)",
			sendErr,
		)
	}

	return nil
}

func handleEditCheckName(conf *config.Config, c tele.Context, state fsm.Context) error {
	check, err := storageHelpers.SetNewCheckNameFromMessage(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't set new check name (%v)",
			err,
		)
	}

	if err := c.Send(responses.GetNewCheckNameIsSavedResponse(check.Name)); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't send a message (%v)",
			err,
		)
	}

	if err := storageHelpers.SetState(static.StateEditingCheck, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't change a state (%v)",
			err,
		)
	}

	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't retrieve an items list (%v)",
			err,
		)
	}

	responseText, _ := responses.GetVerificationFinalStepResponse(check, items)

	if err := c.Send(responses.GetEditCheckMessage(responseText)); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't send a message (%v)",
			err,
		)
	}

	return nil
}
