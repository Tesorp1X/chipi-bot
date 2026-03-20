package callbacks

import (
	"context"
	"errors"
	"fmt"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/db"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func HandleAnyCallback(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error in callbacks.HandleAnyCallback(): couldn't receive users(%d) current state: %v",
			c.Sender().ID, err,
		)
	}

	callbackData := c.Callback().Data

	switch {
	case currentState == static.StateWaitForCheckName &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleKeepCheckNameCallback(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForCheckName', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateShowingAnItem &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleShowingAnItemCallback(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateShowingAnItem', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateWaitForItemOwner &&
		static.CallbackActionEditItem.DataMatches(callbackData):
		if err := handleItemOwnerCallback(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateShowingAnItem', action 'CallbackActionEditItem': %v",
				err,
			)
		}
	case currentState == static.StateWaitingForCheckConfirmation &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleFinalVerificationStage(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitingForCheckConfirmation', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateEditingCheck &&
		static.CallbackActionEditCheck.DataMatches(callbackData):
		if err := handleEditFinalizedCheck(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateEditingCheck', action 'CallbackActionEditCheck': %v",
				err,
			)
		}
	case currentState == static.StateWaitForCheckOwner &&
		static.CallbackActionEditCheck.DataMatches(callbackData):
		if err := handleCheckOwnerCallback(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForCheckOwner', action 'CallbackActionEditCheck': %v",
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

func handleKeepCheckNameCallback(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	c.Respond(&tele.CallbackResponse{})
	// Replying ok!
	if sendErr := c.EditOrSend("Хорошо. Название не меняем👌"); sendErr != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckNameCallback(): couldn't send an 'ok'-message (%v)",
			sendErr,
		)
	}

	if err := storageHelpers.SetState(static.StateWaitForCheckOwner, c, state); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckNameCallback(): couldn't change a state (%v)",
			err,
		)
	}
	// prompt check ownership
	if sendErr := c.Send(responses.GetAskForCheckOwnershipQuestion()); sendErr != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckNameCallback(): couldn't send a 'check-ownership'-message (%v)",
			sendErr,
		)
	}

	return nil
}

func handleCheckOwnerFromEditCheckCallback(c tele.Context, state fsm.Context) error {
	errStorage := state.Update(context.Background(), static.IS_FROM_FINAL_STAGE, nil)

	check, errCheck := storageHelpers.GetCheck(c, state)
	items, errItems := storageHelpers.GetItemsList(c, state)

	var checkStr string

	if errCheck == nil && errItems == nil {
		checkStr, _ = responses.GetVerificationFinalStepResponse(check, items)
	}

	sendErr := c.Send(responses.GetEditCheckMessage(checkStr))
	stateErr := state.SetState(context.Background(), static.StateEditingCheck)

	if sendErr != nil || stateErr != nil || errStorage != nil {
		errMsg := "error in callbacks.handleCheckOwnerFromEditCheckCallback(): "
		if errStorage != nil {
			errMsg += fmt.Sprintf(
				"\ncouldn't delete value with key '%s' from the storage (%v)\n",
				static.IS_FROM_FINAL_STAGE,
				sendErr,
			)
		}

		if sendErr != nil {
			errMsg += fmt.Sprintf(
				"\nfailed to send a message (%v)\n",
				sendErr,
			)
		}

		if stateErr != nil {
			errMsg += fmt.Sprintf(
				"\nfailed to set the state to 'StateEditingCheck' (%v)\n",
				stateErr,
			)
		}

		if err := c.Send("error: " + errMsg); err != nil {
			errMsg += fmt.Sprintf(
				"sent with error (%v)\n",
				err,
			)
		}

		return errors.New(errMsg)
	}

	return nil
}

func handleCheckOwnerCallback(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	// try to set a new owner
	_, err := storageHelpers.SetNewCheckOwnerFromCallback(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't set a check owner (%v)",
			err,
		)
	}

	// check if it's from EditCheck
	var isFromFinalStage bool
	if err := state.Data(
		context.Background(),
		static.IS_FROM_FINAL_STAGE,
		&isFromFinalStage,
	); err == nil && isFromFinalStage {
		return handleCheckOwnerFromEditCheckCallback(c, state)
	}

	// prompt item-verification process
	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't retrieve items (%v)",
			err,
		)
	}

	var currentIndex int

	if err := storageHelpers.SetState(static.StateShowingAnItem, c, state); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't change a state (%v)",
			err,
		)
	}

	responseTxt, kb := responses.GetItemVerificationResponse(
		items[currentIndex],
		currentIndex, len(items),
	)

	if err := storageHelpers.UpdateCurrentItemsIndex(currentIndex, c, state); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't save current index in state-storage (%v)",
			err,
		)
	}

	if sendErr := c.EditOrSend(responseTxt, kb); sendErr != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't send a 'item verification'-message (%v)",
			sendErr,
		)
	}

	return nil
}

func handleShowingAnItemCallback(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	action := static.CallbackActionSelector.GetData(c.Callback().Data)
	switch action {
	case static.CallbackSelectorChange:
		c.Respond(&tele.CallbackResponse{})
		// add new line of text at the bottom of that msg "Что меняем?"
		// change inline kb for that msg to a new one
		// buttons (two in a row): название, цена, кол-во, сумма
		if sendErr := c.EditOrReply(responses.GetEditItemInVerificationResponse(c.Message().Text)); sendErr != nil {
			return fmt.Errorf(
				"error in callbacks.handleShowingAnItemCallback(): couldn't edit a message (%v)",
				sendErr,
			)
		}

		if err := storageHelpers.SetState(static.StateEditingAnItem, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleShowingAnItemCallback(): couldn't change a state (%v).",
				err,
			)
		}

	case static.CallbackSelectorKeep:
		c.Respond(&tele.CallbackResponse{})
		// ask who's the owner of an item
		// new message with a question and inline kb
		// buttons: liz, both, pau
		if sendErr := c.Send(responses.GetItemOwnershipQuestion()); sendErr != nil {
			return fmt.Errorf(
				"error in callbacks.handleShowingAnItemCallback(): couldn't send a message (%v)",
				sendErr,
			)
		}

		if err := storageHelpers.SetState(static.StateWaitForItemOwner, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleShowingAnItemCallback(): couldn't change a state (%v)",
				err,
			)
		}
	default:
		return c.Respond(&tele.CallbackResponse{Text: "error todo"})
	}

	return nil
}

func handleItemOwnerCallback(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	itemOwner := static.CallbackActionEditItem.GetData(c.Callback().Data)
	switch itemOwner {
	case static.CallbackOwnerLiz, static.CallbackOwnerPau, static.CallbackOwnerBoth:
		c.Respond(&tele.CallbackResponse{})
		// get items and currentIndex

		items, err := storageHelpers.GetItemsList(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemOwnerCallback(): couldn't retrieve items (%v)",
				err,
			)
		}

		currentIndex, err := storageHelpers.GetCurrentIndex(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in callbacks.handleCheckName(): couldn't retrieve current index from context (%v)",
				err,
			)
		}

		// check if currentIndex is ok
		if currentIndex < 0 || currentIndex >= len(items) {
			// todo: ehh what to do?..
			// maybe "sorry, there was an error, let's start over"
			// and set currentIndex to 0.
			currentIndex = 0
			if err := storageHelpers.UpdateCurrentItemsIndex(currentIndex, c, state); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): couldn't update current_index_items in context (%v)",
					err,
				)
			}

			// setting up verification process from ground up
			sendErrMsgErr := c.Send("error: items list index is out of bounds. let's start over")
			sendMsgErr := c.Send(responses.GetItemVerificationResponse(items[currentIndex], currentIndex, len(items)))
			stateTransitionErr := storageHelpers.SetState(static.StateShowingAnItem, c, state)
			if sendErrMsgErr != nil || sendMsgErr != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): problem with currentIndex being out of bounds.\nerrorMsg sent with error (%v)\nnew verification message sent with error (%v)\nstate transitioned with an error (%v)",
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
					"error in callbacks.handleItemOwnerCallback(): couldn't retrieve a check (%v)",
					err,
				)
			}

			if err := check.CalculateTotals(items); err != nil {
				sendErr := c.Send("error: while calculating totals")
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): couldn't calculate totals (%v).\nsent with error (%v)",
					err,
					sendErr,
				)
			}

			if err := storageHelpers.UpdateCheck(check, c, state); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): couldn't update check info in context (%v)",
					err,
				)
			}

			if err := c.Send(responses.GetVerificationFinalStepResponse(check, items)); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): couldn't send a message (%v)",
					err,
				)
			}

			if err := storageHelpers.SetState(static.StateWaitingForCheckConfirmation, c, state); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): couldn't change a state (%v)",
					err,
				)
			}

			return nil
		}

		// put updated info back into context
		if err := storageHelpers.UpdateItemsList(items, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemOwnerCallback(): couldn't update items_list in context (%v)",
				err,
			)
		}

		if err := storageHelpers.UpdateCurrentItemsIndex(currentIndex, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemOwnerCallback(): couldn't update current_index_items in context (%v)",
				err,
			)
		}

		// get to the next item -> display that
		err = c.Send(responses.GetItemVerificationResponse(items[currentIndex], currentIndex, len(items)))
		if err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemOwnerCallback(): couldn't send a message with a new item (%v).",
				err,
			)
		}

		// set state to StateShowingAnItem
		if err := storageHelpers.SetState(static.StateShowingAnItem, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemOwnerCallback(): couldn't change a state (%v)",
				err,
			)
		}

	default:
		return c.Respond(&tele.CallbackResponse{Text: "error: invalid response: " + itemOwner})
	}

	return nil
}

func handleFinalVerificationStage(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	action := static.CallbackActionSelector.GetData(c.Callback().Data)

	switch action {
	case static.CallbackSelectorKeep:
		// retrieve check from context
		check, err := storageHelpers.GetCheck(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't set new check name (%v)",
				err,
			)
		}

		items, err := storageHelpers.GetItemsList(c, state)
		if err != nil {
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't retrieve items (%v)",
				err,
			)
		}

		if err := dbs.AddNewCheckWithItems(check, items); err != nil {
			sendErr := c.Send("error: couldn't save your check. try again")
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't save check-data to db (%v). sent with error (%v)",
				err,
				sendErr,
			)
		}

		if err := state.Finish(context.Background(), true); err != nil {
			sendErr := c.Send("error: couldn't finish your state")
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't finish state (%v). sent with error (%v)",
				err,
				sendErr,
			)
		}
		// send an ok msg
		if err := c.EditOrReply(responses.GetCheckSavedMessage(check.Name)); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't send an ok-message (%v)",
				err,
			)
		}
	case static.CallbackSelectorChange:
		if err := storageHelpers.SetState(static.StateEditingCheck, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't change a state (%v)",
				err,
			)
		}

		// add text "Что меняем?"
		// change kb
		// buttons: check name, check owner, items
		return c.EditOrReply(responses.GetEditCheckMessage(c.Message().Text))
	}

	return nil
}

func handleEditFinalizedCheck(dbs *db.DBService, c tele.Context, state fsm.Context) error {
	// figure out an action: what do we change
	whatToChange := static.CallbackActionEditCheck.GetData(c.Callback().Data)
	// retrieve check and items from context
	check, err := storageHelpers.GetCheck(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in callbacks.handleEditFinalizedCheck(): couldn't retrieve a check (%v)",
			err,
		)
	}

	var sendErr, stateErr error
	var action string

	switch whatToChange {
	case static.CallbackEditCheckName:
		sendErr = c.EditOrSend(responses.GetAskForNewCheckNameResponse(check.Name))
		stateErr = storageHelpers.SetState(static.StateWaitForNewCheckName, c, state)
		action = static.CallbackEditCheckName

	case static.CallbackEditCheckOwner:
		state.Update(context.Background(), static.IS_FROM_FINAL_STAGE, true)
		sendErr = c.Send(responses.GetAskForCheckOwnershipQuestion())
		stateErr = storageHelpers.SetState(static.StateWaitForCheckOwner, c, state)
		action = static.CallbackEditCheckOwner

	case static.CallbackEditCheckCreationDate:
		sendErr = c.Send(responses.GetAskForNewCheckCreationDateQuestion())
		stateErr = storageHelpers.SetState(static.StateWaitForCheckCreationDate, c, state)
		action = static.CallbackEditCheckCreationDate

	case static.CallbackEditCheckItems:
		sendErr = nil
		action = static.CallbackEditCheckName

	case static.CallbackSelectorGoBack:
		sendErr = nil
		action = static.CallbackEditCheckName
	}

	if sendErr != nil || stateErr != nil {
		errMsg := "error in callbacks.handleEditFinalizedCheck(): "
		if sendErr != nil {
			errMsg += fmt.Sprintf(
				"\nin action '%s' couldn't send a message (%v)\n",
				action,
				sendErr,
			)
		}

		if stateErr != nil {
			errMsg += fmt.Sprintf(
				"in action '%s' couldn't set a state (%v)\n",
				action,
				stateErr,
			)
		}

		if err := c.Send("error: " + errMsg); err != nil {
			errMsg += fmt.Sprintf(
				"sent with error (%v)\n",
				err,
			)
		}

		return errors.New(errMsg)
	}

	return nil
}
