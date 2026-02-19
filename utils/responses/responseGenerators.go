package responses

import (
	"fmt"
	"time"

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
			Data:   static.CallbackSelectorGoBack,
		},
	)

	return text, kb
}

func GetItemOwnershipQuestion() (string, *tele.ReplyMarkup) {
	text := "Отлично! А чей это товар?👀"
	kb := createSelectorInlineKb(
		2,
		Button{
			BtnTxt: "Liz💜",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackOwnerLiz,
		},
		Button{
			BtnTxt: "Pau💙",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackOwnerPau,
		},
		Button{
			BtnTxt: "Both💜💙",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackOwnerBoth,
		},
	)

	return text, kb
}

func GetVerificationFinalStepResponse(check *static.Check, items []*static.Item) (string, *tele.ReplyMarkup) {
	text := "<b>Почти всё! Проверь чек и можно сохранять.</b>\n\n"

	text += fmt.Sprintf("<b>Название:</b> %s; <b>Дата:</b> %s\n", check.Name, check.Date.Format(time.DateTime))
	text += fmt.Sprintf("<b>Товары: %s</b>\n", sPrintItemsBasedOnOwnership(items))

	text += fmt.Sprintf("<b><i><u>Итого с Пау:</u> %.2f</i></b>\n", check.TotalPau)
	text += fmt.Sprintf("<b><i><u>Итого с Лиз:</u> %.2f</i></b>\n", check.TotalLiz)
	text += fmt.Sprintf("<b><i><u>Итого:</u> %.2f</i></b>\n\n", check.Total)

	switch check.Owner {
	case static.CallbackOwnerLiz:
		text += "Заплатила Лиз💜"
	case static.CallbackOwnerPau:
		text += "Заплатил Пау💙"
	case static.CallbackOwnerBoth:
		text += "Заплатили с совместной карты💜💙"
	}

	kb := createSelectorInlineKb(
		2,
		Button{
			BtnTxt: "All Good✅",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorKeep,
		},
		Button{
			BtnTxt: "Edit✏️",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorChange,
		},
	)

	return text, kb
}

func sPrintItemsBasedOnOwnership(items []*static.Item) string {
	itemToStr := func(idx int, item *static.Item) string {
		return fmt.Sprintf(
			"<i>%d) %s %.2f %.3f %.2f </i>\n",
			idx,
			item.Name,
			item.Price,
			item.Amount,
			item.Subtotal,
		)
	}

	pau := "Товары Пау:\n"
	liz := "Товары Лиз:\n"
	both := "Общие товары:\n"

	var lizIdx, pauIdx, bothIdx int

	for _, item := range items {
		switch item.Owner {
		case static.CallbackOwnerLiz:
			lizIdx++
			liz += itemToStr(lizIdx, item)
		case static.CallbackOwnerPau:
			pauIdx++
			pau += itemToStr(pauIdx, item)
		case static.CallbackOwnerBoth:
			bothIdx++
			both += itemToStr(bothIdx, item)
		}
	}

	text := pau + "\n" + liz + "\n" + both + "\n"

	return text
}

func GetCheckSavedMessage(checkName string) (string, *tele.ReplyMarkup) {
	text := "<b>Чек <i>" + checkName + "</i> сохранен!</b>"
	kb := &tele.ReplyMarkup{}

	return text, kb
}

func GetEditCheckMessage(prevMsg string) (string, *tele.ReplyMarkup) {
	text := prevMsg + "\n\n" + "<b>Что меняем?👀</b>"
	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Название ✏️",
			Unique: static.CallbackActionEditCheck.String(),
			Data:   static.CallbackEditCheckName,
		},
		Button{
			BtnTxt: "Дату 📆",
			Unique: static.CallbackActionEditCheck.String(),
			Data:   static.CallbackEditCheckCreationDate,
		},
		Button{
			BtnTxt: "Кто заплатил 🧑‍🤝‍🧑",
			Unique: static.CallbackActionEditCheck.String(),
			Data:   static.CallbackEditCheckOwner,
		},
		Button{
			BtnTxt: "Изменить товары 📝",
			Unique: static.CallbackActionEditCheck.String(),
			Data:   static.CallbackEditCheckItems,
		},
		Button{
			BtnTxt: "Назад ⬅️",
			Unique: static.CallbackActionEditCheck.String(),
			Data:   static.CallbackSelectorGoBack,
		},
	)

	return text, kb
}

func GetAskForNewCheckNameResponse(currentCheckName string) (string, *tele.ReplyMarkup) {
	text := fmt.Sprintf(
		"Текущее название этого чека: %s\n\nНапиши новое название. Если передумал менять название, просто нажми на кнопку \"Назад ⬅️\"",
		currentCheckName,
	)

	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Назад ⬅️",
			Unique: static.CallbackActionEditCheck.String(),
			Data:   static.CallbackSelectorGoBack,
		},
	)

	return text, kb
}
