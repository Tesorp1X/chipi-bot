package callbacks

import (
	"context"
	"errors"
	"fmt"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/application/prompts"
	"github.com/Tesorp1X/chipi-bot/db"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils"
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
	case utils.ExtractCallbackData(callbackData) == static.CallbackMenuGoBack:
		if err := handleGoBackButtonCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state '%s', 'CallbackMenuGoBack': %v",
				currentState.GoString(),
				err,
			)
		}
	case currentState == static.StateWaitForCheckName &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleKeepCheckNameCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForCheckName', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateShowingAnItem &&
		static.CallbackActionSelector.DataMatches(callbackData):
		if err := handleShowingItemInVerificationCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateShowingAnItem', action 'CallbackActionSelector': %v",
				err,
			)
		}
	case currentState == static.StateShowingAnItemUnsaved:
		if err := handleItemsScrollCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateShowingAnItemUnsaved': %v",
				err,
			)
		}
	case currentState == static.StateEditingAnItem:
		if err := handleEditItemInVerificationCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateEditingAnItem', action 'CallbackActionEditItem': %v",
				err,
			)
		}
	case currentState == static.StateWaitForItemOwner &&
		static.CallbackActionEditItem.DataMatches(callbackData):
		if err := handleItemOwnerCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForItemOwner', action 'CallbackActionEditItem': %v",
				err,
			)
		}
	case currentState == static.StateWaitForCheckConfirmationUnsaved:
		if err := handleFinalVerificationStage(dbs, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitingForCheckConfirmationUnsaved': %v",
				err,
			)
		}
	case currentState == static.StateEditingCheckUnsaved &&
		static.CallbackActionEditUnsavedCheck.DataMatches(callbackData):
		if err := handleEditFinalizedCheck(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateEditingCheck', action 'CallbackActionEditUnsavedCheck': %v",
				err,
			)
		}
	case currentState == static.StateWaitForCheckOwnerUnsaved:
		if err := handleCheckOwnerFromEditCheckCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForCheckOwnerUnsaved': %v",
				err,
			)
		}
	case currentState == static.StateWaitForCheckOwner:
		if err := handleCheckOwnerCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForCheckOwner': %v",
				err,
			)
		}
	case currentState == static.StateEditingAnItemUnsaved:
		if err := handleEditUnsavedItemCallback(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.HandleAnyCallback(), state 'StateWaitForCheckOwner': %v",
				err,
			)
		}
	default:
		// if callback query is old, remove inline buttons from that message
		c.Bot().EditReplyMarkup(c.Callback().Message, &tele.ReplyMarkup{})
		return c.Respond(&tele.CallbackResponse{Text: "error: unsupported action", ShowAlert: true})
	}

	return nil
}

func handleKeepCheckNameCallback(c tele.Context, state fsm.Context) error {
	// Replying ok!
	if err := c.EditOrSend("Хорошо. Название не меняем👌"); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckNameCallback(): couldn't send an 'ok'-message (%v)",
			err,
		)
	}

	if err := prompts.SendCheckOwnershipMessage(prompts.FromAddCheck, c, state); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckNameCallback(): failed to send a check ownership message (%v)",
			err,
		)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckNameCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleCheckOwnerFromEditCheckCallback(c tele.Context, state fsm.Context) error {
	errMsg := "error in callbacks.handleCheckOwnerFromEditCheckCallback():\n"

	if err := storageHelpers.DeleteKeyFromStorage(static.IS_FROM_FINAL_STAGE, c, state); err != nil {
		errMsg += fmt.Sprintf(
			"failed to delete a key from the storage (%v)\n",
			err,
		)
	}

	if err := prompts.SendEditUnsavedCheckMessage(c, state); err != nil {
		errMsg += fmt.Sprintf(
			"prompt failed (%v)\n",
			err,
		)
	}

	if errMsg != "error in callbacks.handleCheckOwnerFromEditCheckCallback():\n" {
		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			errMsg += fmt.Sprintf(
				"failed to respond with errMsg to a user (%v)",
				err,
			)
		}

		return errors.New(errMsg)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleCheckOwnerFromEditCheckCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleCheckOwnerCallback(c tele.Context, state fsm.Context) error {
	// try to set a new owner
	_, err := storageHelpers.SetNewCheckOwnerFromCallback(c, state)
	if err != nil {
		respErr := c.Respond(&tele.CallbackResponse{Text: "Error: failed to update check owner"})
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't set a check owner (%v), responded with error (%v)",
			err,
			respErr,
		)
	}

	var currentIndex int

	if err := storageHelpers.UpdateCurrentItemsIndex(currentIndex, c, state); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): couldn't save current index in state-storage (%v)",
			err,
		)
	}

	if err := prompts.SendShowItemsMessage(prompts.FromAddCheck, c, state); err != nil {
		errMsg := "error in callbacks.handleKeepCheckOwnerCallback():\n"
		errMsg += fmt.Sprintf(
			"prompt failed (%v)",
			err,
		)

		if err := c.Respond(&tele.CallbackResponse{Text: "error!"}); err != nil {
			errMsg += fmt.Sprintf(
				"failed to respond to a callback query (%v)",
				err,
			)
		}

		return errors.New(errMsg)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleKeepCheckOwnerCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

// item in verification stage
func handleShowingItemInVerificationCallback(c tele.Context, state fsm.Context) error {
	action := static.CallbackActionSelector.GetData(c.Callback().Data)
	switch action {
	case static.CallbackSelectorChange:
		if err := prompts.SendShowEditItemOptions(prompts.FromAddCheck, c, state); err != nil {
			errMsg := "error in callbacks.handleShowingAnItemCallback():\n"
			errMsg += fmt.Sprintf(
				"prompt failed (%v)",
				err,
			)

			if err := c.Respond(&tele.CallbackResponse{Text: "error!"}); err != nil {
				errMsg += fmt.Sprintf(
					"failed to respond to a callback query (%v)",
					err,
				)
			}

			return errors.New(errMsg)
		}

	case static.CallbackSelectorKeep:
		if sendErr := c.EditOrSend(responses.GetItemOwnershipQuestion()); sendErr != nil {
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
		if err := c.Respond(&tele.CallbackResponse{Text: "error: unknown action"}); err != nil {
			return fmt.Errorf(
				"failed to respond to a callback query (%v)",
				err,
			)
		}
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleShowingAnItemCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleEditItemInVerificationCallback(c tele.Context, state fsm.Context) error {
	whatToChange := utils.ExtractCallbackData(c.Callback().Data)
	var promptErr error
	switch whatToChange {
	case static.CallbackEditItemName:
		promptErr = prompts.SendChangeItemNameMessage(prompts.FromAddCheck, c, state)
	case static.CallbackEditItemPrice:
		promptErr = prompts.SendNotImplementedMessage(c, state)
	case static.CallbackEditItemAmount:
		promptErr = prompts.SendNotImplementedMessage(c, state)
	case static.CallbackEditItemOwner:
		promptErr = prompts.SendNotImplementedMessage(c, state)
	}

	if promptErr != nil {
		errMsg := "error in callbacks.handleEditItemInVerificationCallback():\n"
		errMsg += fmt.Sprintf(
			"prompt failed (%v)",
			promptErr,
		)

		if err := c.Respond(&tele.CallbackResponse{Text: "error!"}); err != nil {
			errMsg += fmt.Sprintf(
				"failed to respond to a callback query (%v)",
				err,
			)
		}

		return errors.New(errMsg)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleEditItemInVerificationCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleItemOwnerCallback(c tele.Context, state fsm.Context) error {
	itemOwner := static.CallbackActionEditItem.GetData(c.Callback().Data)
	switch itemOwner {
	case static.CallbackOwnerLiz, static.CallbackOwnerPau, static.CallbackOwnerBoth:
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

			if err := prompts.SendShowItemsMessage(prompts.FromAddCheck, c, state); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): prompt failed (%v)",
					err,
				)
			}

			if err := c.Respond(&tele.CallbackResponse{}); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): failed to respond to a callback query (%v)",
					err,
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

			if err := prompts.SendCheckVerificationMessage(c, state); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): failed to send a verification message (%v)",
					err,
				)
			}

			if err := c.Respond(&tele.CallbackResponse{}); err != nil {
				return fmt.Errorf(
					"error in callbacks.handleItemOwnerCallback(): failed to respond to a callback query (%v)",
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
		if err := prompts.SendShowItemsMessage(prompts.FromAddCheck, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemOwnerCallback(): prompt failed (%v)",
				err,
			)
		}

	default:
		return c.Respond(&tele.CallbackResponse{Text: "error: invalid response: " + itemOwner})
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleItemOwnerCallback(): failed to respond to a callback query (%v)",
			err,
		)
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

		if err := storageHelpers.FinishState(static.DELETE_DATA, c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): couldn't finish state (%v).",
				err,
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
		if err := prompts.SendEditUnsavedCheckMessage(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleFinalVerificationStage(): prompt failed (%v)",
				err,
			)
		}
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleFinalVerificationStage(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleEditFinalizedCheck(c tele.Context, state fsm.Context) error {
	whatToChange := static.CallbackActionEditUnsavedCheck.GetData(c.Callback().Data)

	var promptErr error
	var action string
	var storageError error

	switch whatToChange {
	case static.CallbackEditCheckName:
		promptErr = prompts.SendChangeCheckNameMessage(prompts.FromEditCheckFinal, c, state)
		action = static.CallbackEditCheckName

	case static.CallbackEditCheckOwner:
		promptErr = prompts.SendCheckOwnershipMessage(prompts.FromEditCheckFinal, c, state)
		action = static.CallbackEditCheckOwner

	case static.CallbackEditCheckCreationDate:
		promptErr = prompts.SendChangeCreationDateMessage(c, state)
		action = static.CallbackEditCheckCreationDate

	case static.CallbackEditCheckItems:
		action = static.CallbackEditCheckName
		currentIndex := 0
		storageError = storageHelpers.UpdateCurrentItemsIndex(currentIndex, c, state)
		if storageError != nil {
			break
		}

		promptErr = prompts.SendShowItemsMessage(prompts.FromEditCheckFinal, c, state)
	}

	if promptErr != nil || storageError != nil {
		errMsg := "error in callbacks.handleEditFinalizedCheck():\n"
		if promptErr != nil {
			errMsg += fmt.Sprintf(
				"in action '%s' prompt failed (%v)\n",
				action,
				promptErr,
			)
		}

		if storageError != nil {
			errMsg += fmt.Sprintf(
				"in action '%s' storageHelper error (%v)\n",
				action,
				storageError,
			)
		}

		// if err := c.Send("error: " + errMsg); err != nil {
		// 	errMsg += fmt.Sprintf(
		// 		"sent with error (%v)\n",
		// 		err,
		// 	)
		// }

		// todo: better error message
		if err := c.Respond(&tele.CallbackResponse{Text: "error!"}); err != nil {
			errMsg += fmt.Sprintf(
				"failed to respond to a callback query (%v)",
				err,
			)
		}

		return errors.New(errMsg)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleEditFinalizedCheck(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleGoBackButtonCallback(c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		errResp := c.Respond(&tele.CallbackResponse{Text: "error: couldn't get your state"})
		return fmt.Errorf(
			"error in callbacks.handleGoBackButtonCallback(): failed to retrieve a current state (%v), responded with error (%v)",
			err,
			errResp,
		)
	}

	errMsg := "error in callbacks.handleGoBackButtonCallback():\n"
	//userErrMsg := "error: "

	switch currentState {
	case static.StateEditingCheckUnsaved:
		// go to verification stage
		if err := prompts.SendCheckVerificationMessage(c, state); err != nil {
			errMsg += fmt.Sprintf(
				"failed to send a 'check-verification' message (%v)\n",
				err,
			)
		}

	case static.StateWaitForNewCheckNameUnsaved,
		static.StateWaitForCheckOwnerUnsaved,
		static.StateWaitForCheckCreationDateUnsaved,
		static.StateShowingAnItemUnsaved:
		// go to edit unsaved check menu
		if err := prompts.SendEditUnsavedCheckMessage(c, state); err != nil {
			errMsg += fmt.Sprintf(
				"failed to send a 'edit unsaved check' message (%v)\n",
				err,
			)
		}
	case static.StateEditingAnItemUnsaved:
		// go back to items carousel
		if err := prompts.SendShowItemsMessage(prompts.FromEditCheckFinal, c, state); err != nil {
			errMsg += fmt.Sprintf(
				"failed to send a 'show items carousel' message (%v)\n",
				err,
			)
		}
	case static.StateEditingAnItem:
		// go back to items carousel
		if err := prompts.SendShowItemsMessage(prompts.FromAddCheck, c, state); err != nil {
			errMsg += fmt.Sprintf(
				"failed to send a 'show items carousel' message (%v)\n",
				err,
			)
		}
	case static.StateWaitForItemName:
		// go back to show item's edit options (unverified)
		if err := prompts.SendShowEditItemOptions(prompts.FromAddCheck, c, state); err != nil {
			errMsg += fmt.Sprintf(
				"failed to send a 'show items edit options' message (%v)\n",
				err,
			)
		}
	case static.StateWaitForNewItemNameUnsaved:
		// go back to show item's edit options (unsaved)
		if err := prompts.SendShowEditItemOptions(prompts.FromAddCheck, c, state); err != nil {
			errMsg += fmt.Sprintf(
				"failed to send a 'show items edit options' message (%v)\n",
				err,
			)
		}
	default:
		respErr := c.Respond(&tele.CallbackResponse{Text: "error: couldn't take you back"})
		return fmt.Errorf(
			"error in callbacks.handleGoBackButtonCallback(): invalid state (%s) for 'go back' to be called from, responded with err (%v)",
			currentState.GoString(),
			respErr,
		)
	}

	if errMsg != "error in callbacks.handleGoBackButtonCallback():\n" {
		return errors.New(errMsg)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleGoBackButtonCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}

func handleItemsScrollCallback(c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error in callbacks.handleItemsScrollCallback(): failed to retrieve current state (%v)",
			err,
		)
	}

	var cameFrom int
	switch currentState {
	case static.StateShowingAnItemUnsaved:
		cameFrom = prompts.FromEditCheckFinal
	}

	buttonPressed := utils.ExtractCallbackData(c.Callback().Data)
	var promptErr error
	switch buttonPressed {
	case static.CallbackMenuGoBackward, static.CallbackMenuGoForward:
		storageFunc := storageHelpers.IncrementCurrentItemsListIndex
		if buttonPressed == static.CallbackMenuGoBackward {
			storageFunc = storageHelpers.DecrementCurrentItemsListIndex
		}

		if err := storageFunc(c, state); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemsScrollCallback(): failed to change current items index (%v)",
				err,
			)
		}

		promptErr = prompts.SendShowItemsMessage(cameFrom, c, state)

	case static.CallbackSelectorChange:
		promptErr = prompts.SendShowEditItemOptions(cameFrom, c, state)

	default:
		respErr := c.Respond(&tele.CallbackResponse{Text: "error: unknown action"})
		return fmt.Errorf(
			"error in callbacks.handleItemsScrollCallback(): invalid callback data '%s', responded with err (%v)",
			buttonPressed,
			respErr,
		)
	}

	if promptErr != nil {
		errMsg := fmt.Sprintf(
			"error in callbacks.handleItemsScrollCallback(): prompt failed (%v)\n",
			promptErr,
		)
		if err := c.Respond(&tele.CallbackResponse{Text: "error!"}); err != nil {
			errMsg += fmt.Sprintf(
				"failed to respond (%v)",
				err,
			)
		}

		return errors.New(errMsg)
	}

	if !utils.IsCbQueryResponded(c) {
		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			return fmt.Errorf(
				"error in callbacks.handleItemsScrollCallback(): failed to respond to a callback query (%v)",
				err,
			)
		}
	}

	return nil
}

func handleEditUnsavedItemCallback(c tele.Context, state fsm.Context) error {
	whatToChange := utils.ExtractCallbackData(c.Callback().Data)

	var promptErr error
	var action string

	switch whatToChange {
	case static.CallbackEditItemName:
		promptErr = prompts.SendChangeItemNameMessage(prompts.FromEditCheckFinal, c, state)
		action = static.CallbackEditItemName
	case static.CallbackEditItemPrice:
		promptErr = prompts.SendNotImplementedMessage(c, state)
		action = static.CallbackEditItemPrice
	case static.CallbackEditItemAmount:
		promptErr = prompts.SendNotImplementedMessage(c, state)
		action = static.CallbackEditItemAmount
	case static.CallbackEditItemOwner:
		promptErr = prompts.SendNotImplementedMessage(c, state)
		action = static.CallbackEditItemOwner
	}

	if promptErr != nil {
		errMsg := "error in callbacks.handleEditUnsavedItemCallback():\n"

		errMsg += fmt.Sprintf(
			"in action '%s' prompt failed (%v)\n",
			action,
			promptErr,
		)

		if err := c.Send("error: " + errMsg); err != nil {
			errMsg += fmt.Sprintf(
				"sent with error (%v)\n",
				err,
			)
		}

		// todo: better error message
		if err := c.Respond(&tele.CallbackResponse{Text: "error!"}); err != nil {
			errMsg += fmt.Sprintf(
				"failed to respond to a callback query (%v)",
				err,
			)
		}

		return errors.New(errMsg)
	}

	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		return fmt.Errorf(
			"error in callbacks.handleEditUnsavedItemCallback(): failed to respond to a callback query (%v)",
			err,
		)
	}

	return nil
}
