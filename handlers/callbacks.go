package handlers

import (
	"context"
	"log"

	"github.com/Tesorp1X/chipi-bot/db"
	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/Tesorp1X/chipi-bot/util"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Any callbacks handler. Dispatches callbacks to specifick handlers.
func HandleCallbackAction(c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		log.Fatalf("couldn't recieve users(%d) current state: %v", c.Sender().ID, err)
		return err
	}
	switch {
	case currentState == models.StateShowingChecks &&
		models.CallbackActionMenuButtonPress.DataMatches(c.Callback().Data):

		return ShowChecksMenuButtonCallback(c, state)

	case currentState == models.StateWaitForCheckOwner &&
		models.CallbackActionCheckOwner.DataMatches(c.Callback().Data):

		return CheckOwnerCallback(c, state)

	case currentState == models.StateWaitForItemOwner &&
		models.CallbackActionItemOwner.DataMatches(c.Callback().Data):

		return ItemOwnerCallback(c, state)

	case currentState == models.StateWaitForNewItem &&
		models.CallbackActionHasNewItem.DataMatches(c.Callback().Data):

		return NewItemCallback(c, state)
	case currentState == models.StateShowingTotals &&
		models.CallbackActionMenuButtonPress.DataMatches(c.Callback().Data):

		return ShowTotalsMenuButtonCallback(c, state)

	default:
		//Response to callback
		if err := c.Respond(&tele.CallbackResponse{}); err != nil {
			log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
		}
		return c.Send(models.ErrorInvalidRequest)
	}

}

// Handles buuton-presses('<<' and '>>'), while scrolling through checks in '/show checks'.
func ShowChecksMenuButtonCallback(c tele.Context, state fsm.Context) error {
	// Trying to get session from context.
	var session *checksForSession
	if err := state.Data(context.TODO(), models.CHECKS, &session); err != nil {
		session, err = getChecksForCurrentSession(c)
		if err != nil {
			return c.Send(models.ErrorSometingWentWrong)
		}
		state.Update(context.TODO(), models.CHECKS, session)
	}

	var currentIndex int = 0
	// If currentIndex is not stored in context, then it will be just zero.
	if err := state.Data(context.TODO(), models.CURRENT_INDEX, &currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorSometingWentWrong + " Попробуйте еще раз.",
		})
	}
	if currentIndex < 0 || currentIndex >= len(session.Checks) {
		currentIndex = 0
	}

	buttonPressed := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionMenuButtonPress)
	switch buttonPressed {
	case models.FORWARD:
		currentIndex++
		if currentIndex == len(session.Checks) {
			// to eliminate OutOfBounds Error
			if err := c.Respond(&tele.CallbackResponse{
				Text: "Это последний чек!",
			}); err != nil {
				log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
			}
			return nil
		}
	case models.BACK:
		currentIndex--
		if currentIndex < 0 {
			// to eliminate OutOfBounds Error

			if err := c.Respond(&tele.CallbackResponse{
				Text: "Это первый чек!",
			}); err != nil {

				log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
			}
			return nil
		}
	default:
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorInvalidRequest,
		})
	}
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}

	state.Update(context.TODO(), models.CURRENT_INDEX, currentIndex)

	kb := models.CreateSelectorInlineKb(
		2,
		models.Button{
			BtnTxt: "<<",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BACK,
		},
		models.Button{
			BtnTxt: ">>",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.FORWARD,
		},
	)
	// set state ShowinChecks
	if err := state.SetState(context.TODO(), models.StateShowingChecks); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	return c.EditOrReply(util.GetCheckWithItemsResponse(*session.Checks[currentIndex]), kb)
}

// Handles check ownership responce (from inline keyboard).
func CheckOwnerCallback(c tele.Context, state fsm.Context) error {
	//Response to callback
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}
	//Remove keyboard from callback-message
	c.Bot().EditReplyMarkup(c.Message(), &tele.ReplyMarkup{})
	// state: [StateWaitForCheckOwner]; saving check to db and asking to name first item
	checkOwner := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionCheckOwner)
	if err := state.Update(context.TODO(), models.CHECK_OWNER, checkOwner); err != nil {
		return c.Send(models.ErrorStateDataUpdate)
	}
	// save check to db here
	var checkName string
	if err := state.Data(context.Background(), models.CHECK_NAME, &checkName); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}
	var sessionId int64
	if err := state.Data(context.Background(), models.SESSION_ID, &sessionId); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}
	checkId, errDb := db.AddCheck(&models.Check{Name: checkName, Owner: checkOwner}, sessionId)
	if errDb != nil {
		state.Finish(context.TODO(), true)
		state.SetState(context.TODO(), models.StateDefault)
		return c.Send(models.ErrorSavingInDB)
	}

	if err := state.Update(context.Background(), models.CHECK_ID, checkId); err != nil {
		return c.Send(models.ErrorStateDataUpdate)
	}

	msg := util.GetCheckCreatedResponse(checkOwner)

	if err := state.SetState(context.TODO(), models.StateWaitForItemName); err != nil {
		c.Send(models.ErrorSetState)
		return CancelHandler(c, state)
	}
	return c.Send(msg)
}

// Handles item ownership responce (from inline keyboard).v
func ItemOwnerCallback(c tele.Context, state fsm.Context) error {
	//Response to callback
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}
	//Remove keyboard from callback-message
	c.Bot().EditReplyMarkup(c.Message(), &tele.ReplyMarkup{})
	// state: [StateWaitForItemOwner]; saving item to [state.dataStorage] and asking if ther is one more item
	itemOwner := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionItemOwner)
	var (
		itemName  string
		itemPrice float64
	)
	errA := state.Data(context.Background(), models.ITEM_NAME, &itemName)
	errB := state.Data(context.Background(), models.ITEM_PRICE, &itemPrice)

	if errA != nil || errB != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	var checkId int64
	state.Data(context.Background(), models.CHECK_ID, &checkId)

	newItem := models.Item{CheckId: checkId, Name: itemName, Price: itemPrice, Owner: itemOwner}
	itemsList := []models.Item{}
	state.Data(context.Background(), models.ITEMS_LIST, &itemsList)
	itemsList = append(itemsList, newItem)
	if err := state.Update(context.Background(), models.ITEMS_LIST, itemsList); err != nil {
		c.Send(models.ErrorStateDataUpdate)
	}

	msg := util.GetItemAdded(itemOwner, newItem.Price)

	selector := models.CreateSelectorInlineKb(
		2,
		models.Button{BtnTxt: "Да", Unique: models.CallbackActionHasNewItem.String(), Data: models.HAS_MORE_ITEMS_TRUE},
		models.Button{BtnTxt: "Нет", Unique: models.CallbackActionHasNewItem.String(), Data: models.HAS_MORE_ITEMS_FALSE},
	)

	if err := state.SetState(context.TODO(), models.StateWaitForNewItem); err != nil {
		c.Send(models.ErrorSetState)
	}

	return c.Send(msg, selector)
}

// Handles responce to question if there more items (from inline keyboard).
func NewItemCallback(c tele.Context, state fsm.Context) error {
	//Response to callback
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}
	//Remove keyboard from callback-message
	c.Bot().EditReplyMarkup(c.Message(), &tele.ReplyMarkup{})

	hasNewItems := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionHasNewItem)
	var (
		msg      string = "Хорошо!\n"
		newState fsm.State
	)

	switch hasNewItems {
	case models.HAS_MORE_ITEMS_TRUE:
		msg += "Название товара?👀"
		newState = models.StateWaitForItemName
	case models.HAS_MORE_ITEMS_FALSE:
		// todo add list of items
		msg += "Получился вот такой чек:\n"
		itemsList := []models.Item{}
		if err := state.Data(context.Background(), models.ITEMS_LIST, &itemsList); err == fsm.ErrNotFound {
			state.Finish(context.Background(), true)
			return c.Send(models.ErrorItemsListNotFound)
		}

		if err := db.AddItems(itemsList...); err != nil {
			state.Finish(context.Background(), true)
			return c.Send(models.ErrorSavingInDB)
		}

		msg += util.CreateItemsListResponse(itemsList...)
		state.Finish(context.Background(), true)
	}

	if err := state.SetState(context.TODO(), newState); err != nil {
		return c.Send(models.ErrorSetState)
	}
	return c.Send(msg)
}

// Handles buuton-presses('<<' and '>>'), while scrolling through checks in '/show totals'.
func ShowTotalsMenuButtonCallback(c tele.Context, state fsm.Context) error {
	// retrieve totals from contexxt or db
	var totals []*models.SessionTotal
	if err := state.Data(context.TODO(), models.SESSION_TOTALS, &totals); err != nil {
		totals, err = db.GetAllSessionTotals()
		if err != nil {
			return c.Respond(&tele.CallbackResponse{
				Text: models.ErrorSometingWentWrong + " Попробуйте еще раз.",
			})
		}
	}
	// retrieve currentIndex from contexxt or it will be zero
	var currentIndex int = 0
	if err := state.Data(context.TODO(), models.CURRENT_INDEX, &currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorSometingWentWrong + " Попробуйте еще раз.",
		})
	}
	if currentIndex < 0 || currentIndex >= len(totals) {
		currentIndex = 0
	}

	buttonPressed := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionMenuButtonPress)
	switch buttonPressed {
	case models.FORWARD:
		currentIndex++
		if currentIndex == len(totals) {
			// to eliminate OutOfBounds Error
			if err := c.Respond(&tele.CallbackResponse{
				Text: "Это последний отчет!",
			}); err != nil {
				log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
				return err
			}
			return nil
		}
	case models.BACK:
		currentIndex--
		if currentIndex < 0 {
			// to eliminate OutOfBounds Error

			if err := c.Respond(&tele.CallbackResponse{
				Text: "Это первый отчет!",
			}); err != nil {

				log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
				return err
			}
			return nil
		}
	default:
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorInvalidRequest,
		})
	}

	if err := state.Update(context.TODO(), models.CURRENT_INDEX, currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: models.ErrorStateDataUpdate})
	}

	// if all is good, send empty response. testing kinda...
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}

	kb := models.CreateSelectorInlineKb(
		2,
		models.Button{
			BtnTxt: "<<",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BACK,
		},
		models.Button{
			BtnTxt: ">>",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.FORWARD,
		},
	)

	return c.EditOrReply(util.GetShowTotalsResponse(totals[currentIndex]), kb)
}
