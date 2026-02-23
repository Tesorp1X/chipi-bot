package handlers

import (
	"context"
	"fmt"

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
			"error in HandleAnyCallback(): couldn't receive users(%d) current state: %v",
			c.Sender().ID, err,
		)
	}

	switch currentState {
	case static.StateWaitForCheckName:
		if err := handleCheckName(conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyText(), state 'StateWaitForCheckName': %v", err)
		}
	case static.StateWaitForNewCheckName:
		if err := handleEditCheckName(conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyText(), state 'StateWaitForNewCheckName': %v", err)
		}
	}
	return nil
}

func handleCheckName(conf *config.Config, c tele.Context, state fsm.Context) error {
	//get checkData obj from context
	var check *static.Check
	if err := state.Data(context.Background(), static.CHECK, &check); err != nil {
		sendErr := c.Send("error: couldn't retrieve data from context")
		return fmt.Errorf(
			"error in handleCheckName(): couldn't retrieve check from state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}
	//save new check name in it and put new version in context
	check.Name = c.Message().Text
	if err := state.Update(context.Background(), static.CHECK, check); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in handleCheckName(): couldn't save check in state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	//get first item and start verification of items
	var items []*static.Item
	if err := state.Data(context.Background(), static.ITEMS_LIST, &items); err != nil {
		sendErr := c.Send("error: couldn't retrieve data from context")
		return fmt.Errorf(
			"error in handleCheckName(): couldn't retrieve items from state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	if len(items) == 0 {
		sendErr := c.Send("error: retrieved items-list iss empty")
		return fmt.Errorf(
			"error in handleCheckName(): retrieved items-list is empty. send with error: %v",
			sendErr,
		)
	}

	// First message
	// Replying ok!
	if sendErr := c.Send("Отлично, название изменено на <b>" + check.Name + "</b>!"); sendErr != nil {
		return fmt.Errorf(
			"error in handleCheckName(): couldn't send an 'ok'-message (%v)",
			sendErr,
		)
	}

	var currentIndex int

	if err := state.SetState(context.Background(), static.StateShowingAnItem); err != nil {
		sendErr := c.Send("error: couldn't change state")
		return fmt.Errorf(
			"error in handleCheckName(): couldn't change a state to StateShowingAnItems (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	responseTxt, kb := responses.GetItemVerificationResponse(
		items[currentIndex],
		currentIndex, len(items),
	)

	currentIndex++ // now points at the next item
	if err := state.Update(context.Background(), static.CURRENT_INDEX_ITEMS, currentIndex); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in handleCheckName(): couldn't save current index in state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	// Second message
	if sendErr := c.Send(responseTxt, kb); sendErr != nil {
		return fmt.Errorf(
			"error in handleCheckName(): couldn't send a 'item verification'-message (%v)",
			sendErr,
		)
	}

	return nil
}

func handleEditCheckName(conf *config.Config, c tele.Context, state fsm.Context) error {
	check, err := storageHelpers.SetNewCheckName(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handleEditCheckName(): couldn't set new check name (%v)",
			err,
		)
	}

	if err := c.Send(responses.GetNewCheckNameIsSavedResponse(check.Name)); err != nil {
		return fmt.Errorf(
			"error in handleEditCheckName(): couldn't send a message (%v)",
			err,
		)
	}

	if err := state.SetState(context.Background(), static.StateEditingCheck); err != nil {
		sendErr := c.Send("error: couldn't change state")
		return fmt.Errorf(
			"error in handleEditCheckName(): couldn't change a state to StateEditingCheck (%v)\n sent with error (%v)",
			err,
			sendErr,
		)
	}

	//responseText, kb := responses.GetVerificationFinalStepResponse(check, )

	if err := c.Send(responses.GetEditCheckMessage("")); err != nil {
		return fmt.Errorf(
			"error in handleEditCheckName(): couldn't send a message (%v)",
			err,
		)
	}

	return nil
}
