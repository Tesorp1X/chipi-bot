package models

import (
	tele "gopkg.in/telebot.v4"
)

func SelectorInlineKb(btnTxt1, unique1, data1, btnTxt2, unique2, data2 string) *tele.ReplyMarkup {
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
