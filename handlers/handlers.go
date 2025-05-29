package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Tesorp1X/chipi-bot/db"
	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/Tesorp1X/chipi-bot/util"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func CancelHandler(c tele.Context, state fsm.Context) error {
	_ = state.Finish(context.TODO(), c.Data() != "")
	return c.Send("Canceled!")
}

func HelloHandler(c tele.Context, state fsm.Context) error {
	return c.Send("Hello!")
}

func NewCheckHandler(c tele.Context, state fsm.Context) error {
	if err := state.SetState(context.TODO(), models.StateWaitForCheckName); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}
	val, ok := c.Get(models.SESSION_ID).(int64)
	if !ok {
		return c.Send(models.ErrorSometingWentWrong)
	}
	if err := state.Update(context.Background(), models.SESSION_ID, val); err != nil {
		return c.Send(models.ErrorStateDataUpdate)
	}
	return c.Send("Хорошо, как назовем новый чек?👀")
}

func CheckNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}
	state.Update(context.TODO(), models.CHECK_NAME, msgText)
	if err := state.SetState(context.TODO(), models.StateWaitForCheckOwner); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}
	selector := models.CheckOwnershipSelectorInlineKb(
		"Liz :3", models.CallbackActionCheckOwner.String(), models.OWNER_LIZ,
		"Пау <3", models.CallbackActionCheckOwner.String(), models.OWNER_PAU,
	)
	return c.Send("Хорошо. Кто заплатил?🤑", selector)
}

func ItemNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorNameMustBeTxtMsg)
	}
	state.Update(context.TODO(), models.ITEM_NAME, msgText)

	if err := state.SetState(context.TODO(), models.StateWaitForItemPrice); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	return c.Send("Сколько это столо?")
}

func ItemPriceResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	var (
		price float64
		err   error
	)

	if price, err = strconv.ParseFloat(msgText, 64); err != nil {
		return c.Send(models.ErrorItemPriceMustBeANumberMsg)
	}
	state.Update(context.TODO(), models.ITEM_PRICE, price)

	if err := state.SetState(context.TODO(), models.StateWaitForItemOwner); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	selector := models.ItemOwnershipSelectorInlineKb(
		"Liz :3", models.CallbackActionItemOwner.String(), models.OWNER_LIZ,
		"Пау <3", models.CallbackActionItemOwner.String(), models.OWNER_PAU,
		"Оба", models.CallbackActionItemOwner.String(), models.OWNER_BOTH,
	)
	return c.Send("Хорошо. Чей это товар?😺", selector)
}

// /current -- shows how much both payed and who owns money to whom and how much.
func ShowCurrentTotalCommand(c tele.Context, state fsm.Context) error {
	sessionId, ok := c.Get(models.SESSION_ID).(int64)
	if !ok {
		var err error
		sessionId, err = db.GetSessionId()
		if err != nil {
			return c.Send(models.ErrorSometingWentWrong)
		}
	}

	checks, err := db.GetAllChecksWithItemsForSesssionId(sessionId)
	if err != nil {
		return c.Send(err)
	}

	sessionTotal := util.CalculateSessionTotal(sessionId, checks)

	msg := fmt.Sprintf("Вот промежуточный итог за этот период:\nВсего заплачено: %.2f руб\n", sessionTotal.Total)
	if sessionTotal.Recipient == models.OWNER_LIZ {
		msg += fmt.Sprintf("Пау должен Лиз %.2f руб.", sessionTotal.Amount)
	} else {
		msg += fmt.Sprintf("Лиз должна Пау %.2f руб.", sessionTotal.Amount)
	}

	return c.Send(msg)
}

// // /finish -- finishes current session and makes a record in totals table.
// func FinishSession(c tele.Context, state fsm.Context) error {

// }
