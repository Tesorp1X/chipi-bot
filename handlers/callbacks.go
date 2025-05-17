package handlers

import (
	"context"
	"log"
	"strconv"

	"github.com/Tesorp1X/chipi-bot/models"
	"github.com/Tesorp1X/chipi-bot/util"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func HandleCallbackAction(c tele.Context, state fsm.Context) error {
	//Response to callback
	if err := c.Respond(&tele.CallbackResponse{}); err != nil {
		log.Fatalf("couldn't respond to callback %v: %v", c.Callback(), err)
	}
	//Remove keyboard from callback-message
	c.Bot().EditReplyMarkup(c.Message(), &tele.ReplyMarkup{})
	currentState, err := state.State(context.Background())
	if err != nil {
		log.Fatalf("couldn't recieve users(%d) current state: %v", c.Sender().ID, err)
		return err
	}
	switch {
	case currentState == models.StateWaitForCheckOwner && models.CallbackActionCheckOwner.DataMatches(c.Callback().Data):
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
		msg := "Чек создан!😇\n"
		switch checkOwner {
		case models.OWNER_LIZ:
			msg += "Заплатила Лиз💜\n"
		case models.OWNER_PAU:
			msg += "Заплатил Пау💙\n"
		}
		msg += "Теперь давай добавим покупочки😋\n\n"
		msg += "Название товара?👀"

		if err := state.SetState(context.TODO(), models.StateWaitForItemName); err != nil {
			c.Send(models.ErrorSetState)
			return CancelHandler(c, state)
		}
		c.Send(msg)
	case currentState == models.StateWaitForItemOwner && models.CallbackActionItemOwner.DataMatches(c.Callback().Data):
		// state: [StateWaitForItemOwner]; saving item to [state.dataStorage] and asking if ther is one more item
		itemOwner := util.ExtractDataFromCallback(c.Callback().Data, models.CallbackActionItemOwner)
		var (
			itemName  string
			itemPrice int
		)
		errA := state.Data(context.Background(), models.ITEM_NAME, &itemName)
		errB := state.Data(context.Background(), models.ITEM_PRICE, &itemPrice)

		if errA != nil || errB != nil {
			return c.Send(models.ErrorSometingWentWrong)
		}
		// todo: make []Item in storage and save item there.
		msg := "Товар добавлен.\n"
		switch itemOwner {
		case models.OWNER_LIZ:
			msg += "Заплатила Лиз💜\n"
		case models.OWNER_PAU:
			msg += "Заплатил Пау💙\n"
		case models.OWNER_BOTH:
			msg += "Товар общий💜💙\n"
		}
		msg += "Цена: " + strconv.Itoa(itemPrice) + "\n\n"
		msg += "Еще товары?"

		selector := models.CreateSelectorInlineKb(2,
			models.Button{BtnTxt: "Да", Unique: models.HAS_MORE_ITEMS, Data: models.HAS_MORE_ITEMS_TRUE},
			models.Button{BtnTxt: "Нет", Unique: models.HAS_MORE_ITEMS, Data: models.HAS_MORE_ITEMS_FALSE},
		)

		if err := state.SetState(context.TODO(), models.StateWaitForNewItem); err != nil {
			c.Send(models.ErrorSetState)
		}

		return c.Send(msg, selector)

	}

	return nil
}
