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

// Any callbacks handler. Dispatches callbacks to specific handlers.
func HandleCallbackAction(c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		log.Fatalf("couldn't receive users(%d) current state: %v", c.Sender().ID, err)
		return err
	}
	switch {
	case currentState == models.StateShowingChecks &&
		models.CallbackActionMenuButtonPress.DataMatches(c.Callback().Data):

		return ShowChecksScrollButtonCallback(c, state)

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

	case currentState == models.StateEditingCheck &&
		models.CallbackActionEditMenuButtonPress.DataMatches(c.Callback().Data):

		return EditChecksButtonCallback(c, state)

	case currentState == models.StateWaitForCheckOwner &&
		models.CallbackActionEditMenuButtonPress.DataMatches(c.Callback().Data):

		return NewCheckOwnerCallback(c, state)

	default:
		// if callback query is old, remove inline buttons from that message
		c.Bot().EditReplyMarkup(c.Callback().Message, &tele.ReplyMarkup{})
		return c.Respond(&tele.CallbackResponse{Text: models.ErrorInvalidRequest})
	}

}

// Handles button-presses('<<' and '>>'), while scrolling through checks in '/show checks'.
func ShowChecksScrollButtonCallback(c tele.Context, state fsm.Context) error {
	if util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionMenuButtonPress) == models.BTN_EDIT {
		return ShowChecksEditButtonCallback(c, state)
	}
	// Trying to get session from context.
	var session *checksForSession
	if err := state.Data(context.TODO(), models.CHECKS, &session); err != nil {
		session, err = getChecksForCurrentSession(c)
		if err != nil {
			return c.Send(models.ErrorSomethingWentWrong)
		}
		// length is still zero, then there must be no checks for this session yet.
		if len(session.Checks) == 0 {
			return c.EditOrReply("В текущей сессии пока что нет чеков.", &tele.ReplyMarkup{})
		}
		state.Update(context.TODO(), models.CHECKS, session)
	}

	var currentIndex int = 0
	// If currentIndex is not stored in context, then it will be just zero.
	if err := state.Data(context.TODO(), models.CURRENT_INDEX_CHECKS, &currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorSomethingWentWrong + " Попробуйте еще раз.",
		})
	}
	if currentIndex < 0 || currentIndex >= len(session.Checks) {
		currentIndex = 0
	}

	buttonPressed := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionMenuButtonPress)
	switch buttonPressed {
	case models.BTN_FORWARD:
		currentIndex++
		if currentIndex == len(session.Checks) {
			// moving from the last to the first element
			if err := c.Respond(&tele.CallbackResponse{
				Text: "Это был последний чек. Перехожу к первому.",
			}); err != nil {
				log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
				return err
			}
			currentIndex = 0
		}
	case models.BTN_BACK:
		currentIndex--
		if currentIndex < 0 {
			// moving from the first to the last element
			if err := c.Respond(&tele.CallbackResponse{
				Text: "Это первый чек. Перехожу к последнему.",
			}); err != nil {
				log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
				return err
			}
			currentIndex = len(session.Checks) - 1
		}
	default:
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorInvalidRequest,
		})
	}
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}

	state.Update(context.TODO(), models.CURRENT_INDEX_CHECKS, currentIndex)

	kb := models.GetScrollKb()
	// set state ShowingChecks
	if err := state.SetState(context.TODO(), models.StateShowingChecks); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	return c.EditOrReply(util.GetCheckWithItemsResponse(*session.Checks[currentIndex], currentIndex), kb)
}

// Handles button-presses('edit'), while scrolling through checks in '/show checks'.
func ShowChecksEditButtonCallback(c tele.Context, state fsm.Context) error {
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}

	// set state EditingCheck
	if err := state.SetState(context.TODO(), models.StateEditingCheck); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	msg := c.Callback().Message.Text + "\n\nЧто хотите изменить?👀"

	kb := models.CreateSelectorInlineKb(
		3,
		models.Button{
			BtnTxt: "Владелец",
			Unique: models.CallbackActionEditMenuButtonPress.String(),
			Data:   models.CHECK_OWNER,
		},
		models.Button{
			BtnTxt: "Название",
			Unique: models.CallbackActionEditMenuButtonPress.String(),
			Data:   models.CHECK_NAME,
		},
		models.Button{
			BtnTxt: "Товары",
			Unique: models.CallbackActionEditMenuButtonPress.String(),
			Data:   models.ITEMS_LIST,
		},
		models.Button{
			BtnTxt: "Назад",
			Unique: models.CallbackActionEditMenuButtonPress.String(),
			Data:   models.BTN_BACK,
		},
	)

	return c.EditOrReply(msg, kb, tele.ModeMarkdown)
}

// Handles button-presses('edit'), while scrolling through checks in '/show checks'.
func EditChecksButtonCallback(c tele.Context, state fsm.Context) error {
	buttonPressed := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionEditMenuButtonPress)

	// Trying to get session from context.
	var session *checksForSession
	if err := state.Data(context.TODO(), models.CHECKS, &session); err != nil {
		session, err = getChecksForCurrentSession(c)
		if err != nil {
			return c.Send(models.ErrorSomethingWentWrong)
		}
		// length is still zero, then there must be no checks for this session yet.
		if len(session.Checks) == 0 {
			return c.EditOrReply("В текущей сессии покаа что нет чеков.", &tele.ReplyMarkup{})
		}
		state.Update(context.TODO(), models.CHECKS, session)
	}

	var currentIndex int = 0
	// If currentIndex is not stored in context, then it will be just zero.
	if err := state.Data(context.TODO(), models.CURRENT_INDEX_CHECKS, &currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorTryAgain,
		})
	}

	check := session.Checks[currentIndex]

	if err := state.Update(context.TODO(), models.CHECK, check); err != nil {
		state.Finish(context.Background(), true)
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorStateDataUpdate,
		})
	}

	switch buttonPressed {
	case models.CHECK_OWNER:
		// set state EditingCheck
		if err := state.SetState(context.TODO(), models.StateWaitForCheckOwner); err != nil {
			state.Finish(context.TODO(), true)
			return c.Send(models.ErrorSetState)
		}

		msg := "Кто новый владелец?"
		kb := models.CreateSelectorInlineKb(
			2,
			models.Button{
				BtnTxt: "Лиз💜",
				Unique: models.CallbackActionEditMenuButtonPress.String(),
				Data:   models.OWNER_LIZ,
			},
			models.Button{
				BtnTxt: "Пау💙",
				Unique: models.CallbackActionEditMenuButtonPress.String(),
				Data:   models.OWNER_PAU,
			},
		)
		return c.EditOrReply(msg, kb)

	case models.CHECK_NAME:
		if err := state.SetState(context.TODO(), models.StateWaitForNewCheckName); err != nil {
			return c.Respond(&tele.CallbackResponse{Text: models.ErrorTryAgain})
		}
		msg := "Хорошо, какое новое название?"

		return c.EditOrReply(msg, &tele.ReplyMarkup{})

	case models.ITEMS_LIST:
		return c.Respond(&tele.CallbackResponse{Text: "Фича еще в разработке."})
	case models.BTN_BACK:
		return c.Respond(&tele.CallbackResponse{Text: "Фича еще в разработке."})
	default:
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorInvalidRequest,
		})
	}
}

// Assigns new owner to a check in check-editing scenario
func NewCheckOwnerCallback(c tele.Context, state fsm.Context) error {
	// finding check, that is being edited, in storage
	var check *models.CheckWithItems
	if err := state.Data(context.Background(), models.CHECK, &check); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorSetState)
	}

	if err := state.Update(context.TODO(), models.CHECK, nil); err != nil {
		state.Finish(context.TODO(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	newOwner := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionEditMenuButtonPress)

	switch newOwner {
	case models.OWNER_LIZ, models.OWNER_PAU:
		if err := db.EditCheckOwner(check.Id, newOwner); err != nil {
			return c.Respond(&tele.CallbackResponse{Text: models.ErrorTryAgain})
		}
	default:
		c.Bot().EditReplyMarkup(c.Callback().Message, &tele.ReplyMarkup{})
		return c.Respond(&tele.CallbackResponse{Text: models.ErrorInvalidRequest})
	}

	session, err := getChecksForCurrentSession(c)
	if err != nil {
		state.Finish(context.Background(), true)
		return c.Respond(&tele.CallbackResponse{Text: models.ErrorSomethingWentWrong})
	}

	if err := state.Update(context.Background(), models.CHECKS, session); err != nil {
		state.Finish(context.Background(), true)
		return c.Send(models.ErrorStateDataUpdate)
	}

	var currentIndex int
	if err := state.Data(context.TODO(), models.CURRENT_INDEX_CHECKS, &currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorTryAgain,
		})
	}

	kb := models.GetScrollKb()

	if err := state.SetState(context.TODO(), models.StateShowingChecks); err != nil {
		state.Finish(context.Background(), true)
		return c.Send(models.ErrorSetState)
	}

	msg := "Готово!\n\n" + util.GetCheckWithItemsResponse(*session.Checks[currentIndex], currentIndex)
	return c.EditOrReply(msg, kb)
}

// Handles check ownership response (from inline keyboard).
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
		return c.Send(models.ErrorSomethingWentWrong)
	}
	var sessionId int64
	if err := state.Data(context.Background(), models.SESSION_ID, &sessionId); err != nil {
		return c.Send(models.ErrorSomethingWentWrong)
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

// Handles item ownership response (from inline keyboard).v
func ItemOwnerCallback(c tele.Context, state fsm.Context) error {
	//Response to callback
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}
	//Remove keyboard from callback-message
	c.Bot().EditReplyMarkup(c.Message(), &tele.ReplyMarkup{})
	// state: [StateWaitForItemOwner]; saving item to [state.dataStorage] and asking if there is one more item
	itemOwner := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionItemOwner)
	var (
		itemName  string
		itemPrice float64
	)
	errA := state.Data(context.Background(), models.ITEM_NAME, &itemName)
	errB := state.Data(context.Background(), models.ITEM_PRICE, &itemPrice)

	if errA != nil || errB != nil {
		return c.Send(models.ErrorSomethingWentWrong)
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

// Handles response to question if there more items (from inline keyboard).
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

// Handles button-presses('<<' and '>>'), while scrolling through checks in '/show totals'.
func ShowTotalsMenuButtonCallback(c tele.Context, state fsm.Context) error {
	// retrieve totals from context or db
	var totals []*models.SessionTotal
	if err := state.Data(context.TODO(), models.SESSION_TOTALS, &totals); err != nil {
		totals, err = db.GetAllSessionTotals()
		if err != nil {
			return c.Respond(&tele.CallbackResponse{
				Text: models.ErrorSomethingWentWrong + " Попробуйте еще раз.",
			})
		}
	}
	// retrieve currentIndex from context or it will be zero
	var currentIndex int = 0
	if err := state.Data(context.TODO(), models.CURRENT_INDEX_TOTALS, &currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorSomethingWentWrong + " Попробуйте еще раз.",
		})
	}
	if currentIndex < 0 || currentIndex >= len(totals) {
		currentIndex = 0
	}

	buttonPressed := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionMenuButtonPress)
	switch buttonPressed {
	case models.BTN_FORWARD:
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
	case models.BTN_BACK:
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
	case models.BTN_CHECKS:

		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorNotImplemented,
		})

	default:
		return c.Respond(&tele.CallbackResponse{
			Text: models.ErrorInvalidRequest,
		})
	}

	if err := state.Update(context.TODO(), models.CURRENT_INDEX_TOTALS, currentIndex); err != nil {
		return c.Respond(&tele.CallbackResponse{Text: models.ErrorStateDataUpdate})
	}

	// if all is good, send empty response. testing kinda...
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}

	kb := models.CreateSelectorInlineKb(
		3,
		models.Button{
			BtnTxt: "<<",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BTN_BACK,
		},
		models.Button{
			BtnTxt: "Чеки",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BTN_CHECKS,
		},
		models.Button{
			BtnTxt: ">>",
			Unique: models.CallbackActionMenuButtonPress.String(),
			Data:   models.BTN_FORWARD,
		},
	)

	return c.EditOrReply(util.GetShowTotalsResponse(totals[currentIndex]), kb)
}
