package callbacks

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
			return fmt.Errorf(
				"error in HandleAnyCallback(), state 'StateWaitForCheckName', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateShowingAnItem &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleShowingAnItemCallback(conf, c, state); err != nil {
			return fmt.Errorf(
				"error in HandleAnyCallback(), state 'StateShowingAnItem', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateWaitForItemOwner &&
		static.CallbackActionEditItem.DataMatches(callbackData):
		if err := handleItemOwnerCallback(conf, c, state); err != nil {
			return fmt.Errorf(
				"error in HandleAnyCallback(), state 'StateShowingAnItem', action 'CallbackActionEditItem': %v",
				err,
			)
		}
	case currentState == static.StateWaitingForCheckConfirmation &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleFinalVerificationStage(conf, c, state); err != nil {
			return fmt.Errorf(
				"error in HandleAnyCallback(), state 'StateWaitingForCheckConfirmation', action 'CallbackActionSelector': %v",
				err,
			)
		}

	case currentState == static.StateEditingCheck &&
		static.CallbackActionEditCheck.DataMatches(callbackData):
		if err := handleEditFinalizedCheck(conf, c, state); err != nil {
			return fmt.Errorf(
				"error in HandleAnyCallback(), state 'StateEditingCheck', action 'CallbackActionSelector': %v",
				err,
			)
		}
	default:
		// if callback query is old, remove inline buttons from that message
		c.Bot().EditReplyMarkup(c.Callback().Message, &tele.ReplyMarkup{})
		return c.Respond(&tele.CallbackResponse{Text: "error todo", ShowAlert: true})
	}
	return nil
}

func handleKeepCheckNameCallback(conf *config.Config, c tele.Context, state fsm.Context) error {
	c.Respond(&tele.CallbackResponse{})
	// send ok-message
	// show first item and increment currentIndex

	//get first item and start verification of items
	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handleItemOwnerCallback(): couldn't retrieve items (%v)",
			err,
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
	itemOwner := static.CallbackActionEditItem.GetData(c.Callback().Data)
	switch itemOwner {
	case static.CallbackOwnerLiz, static.CallbackOwnerPau, static.CallbackOwnerBoth:
		c.Respond(&tele.CallbackResponse{})
		// get items and currentIndex

		items, err := storageHelpers.GetItemsList(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in handleItemOwnerCallback(): couldn't retrieve items (%v)",
				err,
			)
		}

		currentIndex, err := storageHelpers.GetCurrentIndex(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in handleCheckName(): couldn't retrieve current index from context (%v)",
				err,
			)
		}

		if len(items) == 0 {
			sendErr := c.Send("error: items list is empty")
			stateErr := state.SetState(context.Background(), static.StateDefault)
			return fmt.Errorf(
				"error in handleItemOwnerCallback(): slice of items is empty. send with error: %v, state transition err: %v",
				sendErr,
				stateErr,
			)
		}

		// check if currentIndex is ok
		if currentIndex < 0 || currentIndex >= len(items) {
			// todo: ehh what to do?..
			// maybe "sorry, there was an error, let's start over"
			// and set currentIndex to 0.
			currentIndex = 0
			if err := state.Update(context.Background(), static.CURRENT_INDEX_ITEMS, currentIndex); err != nil {
				sendErr := c.Send("error: couldn't update items current index")
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): couldn't update current_index_items in context (%v). send with error: %v",
					err,
					sendErr,
				)
			}

			// setting up verification process from ground up
			sendErrMsgErr := c.Send("error: items list index is out of bounds. let's start over")
			sendMsgErr := c.Send(responses.GetItemVerificationResponse(items[currentIndex], currentIndex, len(items)))
			stateTransitionErr := state.SetState(context.Background(), static.StateShowingAnItem)
			if sendErrMsgErr != nil || sendMsgErr != nil {
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): problem with currentIndex being out of bounds.\nerrorMsg sent with error (%v)\nnew verification message sent with error (%v)\nstate transitioned with an error (%v)",
					sendErrMsgErr,
					sendMsgErr,
					stateTransitionErr,
				)
			}
			return nil
		}

		// update items info
		items[currentIndex].Owner = itemOwner
		currentIndex++
		//check if it was the last item
		if currentIndex == len(items) {
			// todo
			// make a show check func idk
			// send a message with all new info and calculated totals
			// has inline buttons: all good (saves it all to db) and edit (lets edit name and items)
			check, err := storageHelpers.GetCheck(c, state)
			if err != nil {
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): couldn't retrieve a check (%v)",
					err,
				)
			}

			if err := check.CalculateTotals(items); err != nil {
				sendErr := c.Send("error: while calculating totals")
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): couldn't calculate totals (%v).\nsent with error (%v)",
					err,
					sendErr,
				)
			}

			if err := state.Update(context.Background(), static.CHECK, check); err != nil {
				sendErr := c.Send("error: couldn't update check info in context")
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): couldn't update check info in context (%v).\nsent with error (%v)",
					err,
					sendErr,
				)
			}

			if err := c.Send(responses.GetVerificationFinalStepResponse(check, items)); err != nil {
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): couldn't send a message (%v)",
					err,
				)
			}

			if err := state.SetState(context.Background(), static.StateWaitingForCheckConfirmation); err != nil {
				sendErr := c.Send("error: couldn't change the state")
				return fmt.Errorf(
					"error in handleItemOwnerCallback(): couldn't change the state to StateWaitingForCheckConfirmation (%v).\nsent with error (%v)",
					err,
					sendErr,
				)
			}

			return nil
		}

		// put updated info back into context
		if err := state.Update(context.Background(), static.ITEMS_LIST, items); err != nil {
			sendErr := c.Send("error: couldn't update items_list")
			return fmt.Errorf(
				"error in handleItemOwnerCallback(): couldn't update items_list in context (%v). send with error: %v",
				err,
				sendErr,
			)
		}

		if err := state.Update(context.Background(), static.CURRENT_INDEX_ITEMS, currentIndex); err != nil {
			sendErr := c.Send("error: couldn't update current items's index")
			return fmt.Errorf(
				"error in handleItemOwnerCallback(): couldn't update currentIndex in context (%v). send with error: %v",
				err,
				sendErr,
			)
		}

		// get to the next item -> display that
		err = c.Send(responses.GetItemVerificationResponse(items[currentIndex], currentIndex, len(items)))
		if err != nil {
			return fmt.Errorf(
				"error in handleItemOwnerCallback(): couldn't send a message with a new item (%v).",
				err,
			)
		}

		// set state to StateShowingAnItem
		if err := state.SetState(context.Background(), static.StateShowingAnItem); err != nil {
			sendErr := c.Send("error: couldn't set state to StateShowingAnItem")
			return fmt.Errorf(
				"error in handleItemOwnerCallback(): couldn't set state to StateShowingAnItem (%v). send with error: %v",
				err,
				sendErr,
			)
		}

	default:
		return c.Respond(&tele.CallbackResponse{Text: "error todo"})
	}

	return nil
}

func handleFinalVerificationStage(conf *config.Config, c tele.Context, state fsm.Context) error {
	action := static.CallbackActionSelector.GetData(c.Callback().Data)

	switch action {
	case static.CallbackSelectorKeep:
		// retrieve check and items from context
		check, err := storageHelpers.SetNewCheckNameFromMessage(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in handleFinalVerificationStage(): couldn't set new check name (%v)",
				err,
			)
		}

		// getting items
		// items, err := storageHelpers.GetItemsList(c, state)
		// if err != nil {
		// 	return fmt.Errorf(
		// 		"error in handleFinalVerificationStage(): couldn't retrieve items (%v)",
		// 		err,
		// 	)
		// }

		// get session id
		// assign session id to a check and save it to a db
		// assign check id to every item and save it to a db
		// transition state to Default
		if err := state.Finish(context.Background(), true); err != nil {
			sendErr := c.Send("error: couldn't finish your state")
			return fmt.Errorf(
				"error in handleFinalVerificationStage(): couldn't finish state (%v). sent with error (%v)",
				err,
				sendErr,
			)
		}
		// send an ok msg
		if err := c.EditOrReply(responses.GetCheckSavedMessage(check.Name)); err != nil {
			return fmt.Errorf(
				"error in handleFinalVerificationStage(): couldn't send an ok-message (%v)",
				err,
			)
		}

	case static.CallbackSelectorChange:
		if err := state.SetState(context.Background(), static.StateEditingCheck); err != nil {
			sendErr := c.Send("error: couldn't set state to StateEditingCheck")
			return fmt.Errorf(
				"error in handleFinalVerificationStage(): couldn't set state to StateEditingCheck (%v). send with error: %v",
				err,
				sendErr,
			)
		}
		// add text "Что меняем?"
		// change kb
		// buttons: check name, check owner, items
		return c.EditOrReply(responses.GetEditCheckMessage(c.Message().Text))
	}

	return nil
}

func handleEditFinalizedCheck(conf *config.Config, c tele.Context, state fsm.Context) error {
	// figure out an action: what do we change
	whatToChange := static.CallbackActionEditCheck.GetData(c.Callback().Data)
	// retrieve check and items from context
	check, err := storageHelpers.GetCheck(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handleEditFinalizedCheck(): couldn't retrieve a check (%v)",
			err,
		)
	}

	var sendErr, stateErr error
	var action string

	switch whatToChange {
	case static.CallbackEditCheckName:
		sendErr = c.EditOrSend(responses.GetAskForNewCheckNameResponse(check.Name))
		stateErr = state.SetState(context.Background(), static.StateWaitForNewCheckName)
		action = static.CallbackEditCheckName

	case static.CallbackEditCheckOwner:
		sendErr = nil
		action = static.CallbackEditCheckOwner

	case static.CallbackEditCheckCreationDate:
		sendErr = nil
		action = static.CallbackEditCheckCreationDate

	case static.CallbackEditCheckItems:
		sendErr = nil
		action = static.CallbackEditCheckName

	case static.CallbackSelectorGoBack:
		sendErr = nil
		action = static.CallbackEditCheckName
	}

	if sendErr != nil {
		sendErrorMsgErr := c.Send("error: couldn't retrieve check info from context")
		return fmt.Errorf(
			"error in handleEditFinalizedCheck(): in action '%s' couldn't send a message (%v).\nsent with error (%v)",
			action,
			sendErr,
			sendErrorMsgErr,
		)
	}

	if stateErr != nil {
		sendErrorMsgErr := c.Send("error: couldn't change a state")
		return fmt.Errorf(
			"error in handleEditFinalizedCheck(): in action '%s' couldn't set a state (%v).\nsent with error (%v)",
			action,
			stateErr,
			sendErrorMsgErr,
		)
	}

	return nil
}
