package callbacks

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"
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
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleKeepCheckNameCallback(conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyCallback(), state 'StateWaitForCheckName', action 'CallbackActionSelector': %v", err)
		}
	case currentState == static.StateShowingAnItem &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleShowingAnItemCallback(conf, c, state); err != nil {
			return fmt.Errorf("error in HandleAnyCallback(), state 'StateShowingAnItem', action 'CallbackActionSelector': %v", err)
		}
	default:
		// if callback query is old, remove inline buttons from that message
		c.Bot().EditReplyMarkup(c.Callback().Message, &tele.ReplyMarkup{})
		return c.Respond(&tele.CallbackResponse{Text: "error todo"})
	}
	return nil
}

func handleKeepCheckNameCallback(conf *config.Config, c tele.Context, state fsm.Context) error {
	c.Respond(&tele.CallbackResponse{})
	// send ok-message
	// show first item and increment currentIndex

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
	if sendErr := c.Send("Хорошо. Название не меняем👌"); sendErr != nil {
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

func handleShowingAnItemCallback(conf *config.Config, c tele.Context, state fsm.Context) error {
	action := static.CallbackActionSelector.GetData(c.Callback().Data)
	switch action {
	case static.CallbackSelectorChange:
		c.Respond(&tele.CallbackResponse{})
		// add new line of text at the bottom of that msg "Что меняем?"
		// change inline kb for that msg to a new one
		// buttons (two in a row): название, цена, кол-во, сумма
		if sendErr := c.EditOrReply(responses.GetEditItemInVerificationResponse(c.Message().Text)); sendErr != nil {
			return fmt.Errorf(
				"error in handleShowingAnItemCallback(): couldn't edit a message (%v)",
				sendErr,
			)
		}
		// state to waitingForMenuAction maybe
		if err := state.SetState(context.Background(), static.StateEditingAnItem); err != nil {
			sendErr := c.Send("error: couldn't change state")
			return fmt.Errorf(
				"error in handleShowingAnItemCallback(): couldn't change a state to StateEditingAnItem (%v). send with error: %v",
				err,
				sendErr,
			)
		}
	case static.CallbackSelectorKeep:
		c.Respond(&tele.CallbackResponse{})
		// ask who's the owner of an item
		// new message with a question and inline kb
		// buttons: liz, both, pau
		if sendErr := c.Send(responses.GetItemOwnershipQuestion()); sendErr != nil {
			return fmt.Errorf(
				"error in handleShowingAnItemCallback(): couldn't send a message (%v)",
				sendErr,
			)
		}

		if err := state.SetState(context.Background(), static.StateWaitForItemOwner); err != nil {
			sendErr := c.Send("error: couldn't change state")
			return fmt.Errorf(
				"error in handleShowingAnItemCallback(): couldn't change a state to StateWaitForItemOwner (%v). send with error: %v",
				err,
				sendErr,
			)
		}

	default:
		return c.Respond(&tele.CallbackResponse{Text: "error todo"})
	}

	return nil
}

func handleItemOwnerCallback(conf *config.Config, c tele.Context, state fsm.Context) error {
	action := static.CallbackActionEditItem.GetData(c.Callback().Data)
	switch action {
	case static.CallbackOwnerLiz, static.CallbackOwnerPau, static.CallbackOwnerBoth:
		c.Respond(&tele.CallbackResponse{})
		// get items and currentIndex
		// update items info
		// get to the next item -> display that
	default:
		return c.Respond(&tele.CallbackResponse{Text: "error todo"})
	}

	return nil
}
