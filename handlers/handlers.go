package handlers

import (
	"context"

	"github.com/Tesorp1X/chipi-bot/models"

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
	return c.Send("Хорошо, как назовем новый чек?👀")
}

func CheckNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorCheckNameMustBeTxtMsg)
	}
	state.Update(context.TODO(), models.CHECK_NAME, msgText)
	if err := state.SetState(context.TODO(), models.StateWaitForCheckOwner); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}
	selector := models.SelectorInlineKb(
		"Liz :3", models.CallbackActionCheckOwner.String(), models.OWNER_LIZ,
		"Пау <3", models.CallbackActionCheckOwner.String(), models.OWNER_PAU,
	)
	return c.Send("Хорошо. Кто заплатил?🤑", selector)
}

func ItemNameResponseHandler(c tele.Context, state fsm.Context) error {
	msgText := c.Text()
	if len(msgText) == 0 {
		return c.Send(models.ErrorCheckNameMustBeTxtMsg)
	}
	state.Update(context.TODO(), models.ITEM_NAME, msgText)

	if err := state.SetState(context.TODO(), models.StateWaitForItemPrice); err != nil {
		return c.Send(models.ErrorSometingWentWrong)
	}

	return c.Send("Сколько это столо?")
}
