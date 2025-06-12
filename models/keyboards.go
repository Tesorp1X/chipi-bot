package models

import (
	tele "gopkg.in/telebot.v4"
)

type Button struct {
	BtnTxt string
	Unique string
	Data   string
}

// ReplyMarkup for selecting how owns a check
func CheckOwnershipSelectorInlineKb(btnTxt1, unique1, data1, btnTxt2, unique2, data2 string) *tele.ReplyMarkup {
	var (
		// Universal markup builders.
		selector = &tele.ReplyMarkup{}

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind
		// since it's required for callback routing to work.
		//
		btnLeft  = selector.Data(btnTxt1, unique1, data1)
		btnRight = selector.Data(btnTxt2, unique2, data2)
	)
	selector.Inline(
		selector.Row(btnLeft, btnRight),
	)
	return selector
}

// ReplyMarkup for selecting how owns an item
func ItemOwnershipSelectorInlineKb(btnTxt1, unique1, data1, btnTxt2, unique2, data2, btnTxt3, unique3, data3 string) *tele.ReplyMarkup {
	var (
		// Universal markup builders.
		selector = &tele.ReplyMarkup{}

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind
		// since it's required for callback routing to work.
		//
		btnLeft   = selector.Data(btnTxt1, unique1, data1)
		btnRight  = selector.Data(btnTxt2, unique2, data2)
		btnMiddle = selector.Data(btnTxt3, unique3, data3)
	)
	selector.Inline(
		selector.Row(btnLeft, btnRight),
		selector.Row(btnMiddle),
	)
	return selector
}

// Selector kb factory
func CreateSelectorInlineKb(btnsPerRow int, buttons ...Button) *tele.ReplyMarkup {
	rowsOfButtons := []tele.Btn{}
	for _, b := range buttons {
		rowsOfButtons = append(rowsOfButtons, tele.Btn{
			Text:   b.BtnTxt,
			Unique: b.Unique,
			Data:   b.Data,
		})
	}
	rm := &tele.ReplyMarkup{}
	rm.Inline(rm.Split(btnsPerRow, rowsOfButtons)...)

	return rm
}

func GetScrollKb() *tele.ReplyMarkup {
	return CreateSelectorInlineKb(3,
		Button{
			BtnTxt: "<<",
			Unique: CallbackActionMenuButtonPress.String(),
			Data:   BTN_BACK,
		},
		Button{
			BtnTxt: "edit",
			Unique: CallbackActionMenuButtonPress.String(),
			Data:   BTN_EDIT,
		},
		Button{
			BtnTxt: ">>",
			Unique: CallbackActionMenuButtonPress.String(),
			Data:   BTN_FORWARD,
		})
}
