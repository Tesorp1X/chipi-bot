package responses

import (
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"
	tele "gopkg.in/telebot.v4"
)

type Button struct {
	BtnTxt string
	Unique string
	Data   string
}

// Selector kb factory
func createSelectorInlineKb(btnsPerRow int, buttons ...Button) *tele.ReplyMarkup {
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

func GenerateNameVerificationResponse(checkName string) (string, *tele.ReplyMarkup) {
	response := `<b>Сканирование завершено! Уточним данные.</b>` + "\n\n"
	response += `Начнем с названия. Я предлагаю назвать <b>` + checkName + "</b>." + "\n\n"
	response += `Оставляем? Если хочешь поменять, то отправь новое название.`

	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Оставляем",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorKeep,
		},
	)

	return response, kb
}

func GetItemVerificationResponse(item *static.Item, currentIndex, outOf int) (string, *tele.ReplyMarkup) {
	response := fmt.Sprintf(
		"<b>Проверяем товары: %d из %d</b>\n\n",
		currentIndex+1,
		outOf,
	)

	response += fmt.Sprintf(
		"<i>Название:</i> %s\n<i>Цена:</i> %.2f\n<i>Кол-во:</i> %.3f\n<i>Сумма:</i> %.2f\n",
		item.Name,
		item.Price,
		item.Amount,
		item.Subtotal,
	)

	kb := createSelectorInlineKb(
		2,
		Button{
			BtnTxt: "All good✅",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorKeep,
		},
		Button{
			BtnTxt: "Изменить",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorChange,
		},
	)

	return response, kb
}

func GetEditItemInVerificationResponse(msgText string) (string, *tele.ReplyMarkup) {
	text := msgText + "\n\nЧто меняем?👀"
	kb := createSelectorInlineKb(
		2,
		Button{
			BtnTxt: "Название",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackEditItemName,
		},
		Button{
			BtnTxt: "Цена",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackEditItemPrice,
		},
		Button{
			BtnTxt: "Кол-во",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackEditItemAmount,
		},
		Button{
			BtnTxt: "Сумма",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackEditItemSubtotal,
		},
		Button{
			BtnTxt: "Вернуться⬅️",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackEditItemGoBack,
		},
	)

	return text, kb
}
