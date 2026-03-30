package responses

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tesorp1X/chipi-bot/static"
	tele "gopkg.in/telebot.v4"
)

type Button struct {
	BtnTxt string
	Unique string
	Data   string
}

type RowOfButtons struct {
	BtnsPerRow int
	Btns       []Button
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

func makeTeleRowFromButtons(amountOfBtns int, buttons ...Button) tele.Row {
	var row tele.Row
	for i, b := range buttons {
		if i >= amountOfBtns {
			// the rest of btns are cut of
			break
		}
		row = append(row,
			tele.Btn{
				Text:   b.BtnTxt,
				Unique: b.Unique,
				Data:   b.Data,
			},
		)
	}

	return row
}

// Allows to create Inline Keyboards with different amount of buttons per row.
// Returns a ready to use keyboard.
func createCustomRowsInlineKb(rows ...RowOfButtons) *tele.ReplyMarkup {
	kb := &tele.ReplyMarkup{}
	var teleRows []tele.Row
	for _, row := range rows {
		teleRows = append(teleRows, makeTeleRowFromButtons(row.BtnsPerRow, row.Btns...))
	}

	kb.Inline(teleRows...)
	return kb
}

func GetVerificationFirstMessage(checkName string) (string, *tele.ReplyMarkup) {
	response := `<b>Сканирование завершено! Уточним данные.</b>` + "\n\n"
	response += `Начнем с названия. Я предлагаю назвать <b>` + checkName + "</b>." + "\n\n"
	response += `Оставляем? Если хочешь поменять, то отправь новое название.`

	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Оставляем ✅",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorKeep,
		},
	)

	return response, kb
}

func sPrintItemInfo(item *static.Item) string {
	return fmt.Sprintf(
		"<i>Название:</i> %s\n<i>Цена:</i> %.2f\n<i>Кол-во:</i> %.3f\n<i>Сумма:</i> %.2f\n",
		item.Name,
		item.Price,
		item.Amount,
		item.Subtotal,
	)
}

func GetItemVerificationResponse(item *static.Item, currentIndex, outOf int) (string, *tele.ReplyMarkup) {
	response := fmt.Sprintf(
		"<b>Проверяем товары: %d из %d</b>\n\n",
		currentIndex+1,
		outOf,
	)

	response += sPrintItemInfo(item)

	kb := createSelectorInlineKb(
		2,
		Button{
			BtnTxt: "All good ✅",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorKeep,
		},
		Button{
			BtnTxt: "Изменить ✏️",
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
			BtnTxt: "Вернуться ⬅️",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackMenuGoBack,
		},
	)

	return text, kb
}

func GetItemOwnershipQuestion() (string, *tele.ReplyMarkup) {
	text := "<b>Отлично! А чей это товар?👀</b>"
	kb := createSelectorInlineKb(
		2,
		Button{
			BtnTxt: "Liz 💜",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackOwnerLiz,
		},
		Button{
			BtnTxt: "Pau 💙",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackOwnerPau,
		},
		Button{
			BtnTxt: "Общий 💜💙",
			Unique: static.CallbackActionEditItem.String(),
			Data:   static.CallbackOwnerBoth,
		},
	)

	return text, kb
}

func GetVerificationFinalStepResponse(check *static.Check, items []*static.Item) (string, *tele.ReplyMarkup) {
	text := "<b>Почти всё! Проверь чек и можно сохранять.</b>\n\n"

	text += fmt.Sprintf("<b>Название:</b> %s; <b>Дата:</b> %s\n", check.Name, check.Date.Format(time.DateTime))
	text += fmt.Sprintf("<b>Товары:\n</b>%s\n", sPrintItemsBasedOnOwnership(items))

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
			BtnTxt: "All Good ✅",
			Unique: static.CallbackActionSelector.String(),
			Data:   static.CallbackSelectorKeep,
		},
		Button{
			BtnTxt: "Изменить ✏️",
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

	var pau strings.Builder
	pau.WriteString("<b><i><u>Товары Пау:</u></i></b>\n")
	var liz strings.Builder
	liz.WriteString("<b><i><u>Товары Лиз:</u></i></b>\n")
	var both strings.Builder
	both.WriteString("<b><i><u>Общие товары:</u></i></b>\n")

	var lizIdx, pauIdx, bothIdx int

	for _, item := range items {
		switch item.Owner {
		case static.CallbackOwnerLiz:
			lizIdx++
			liz.WriteString(itemToStr(lizIdx, item))
		case static.CallbackOwnerPau:
			pauIdx++
			pau.WriteString(itemToStr(pauIdx, item))
		case static.CallbackOwnerBoth:
			bothIdx++
			both.WriteString(itemToStr(bothIdx, item))
		}
	}

	addNoItemsMsgIfEmpty := func(idx int, sb *strings.Builder) {
		if idx == 0 {
			sb.WriteString("<i>Товаров нет</i>\n")
		}
	}

	addNoItemsMsgIfEmpty(lizIdx, &liz)
	addNoItemsMsgIfEmpty(pauIdx, &pau)
	addNoItemsMsgIfEmpty(bothIdx, &both)

	text := pau.String() + "\n" + liz.String() + "\n" + both.String() + "\n"

	return text
}

func GetCheckSavedMessage(checkName string) (string, *tele.ReplyMarkup) {
	text := "<b>Чек <i>" + checkName + "</i> сохранен!</b>"
	kb := &tele.ReplyMarkup{}

	return text, kb
}

func GetShowCheckEditOptionsMessage(prevMsg string) (string, *tele.ReplyMarkup) {
	text := prevMsg + "\n\n" + "<b>Что меняем?👀</b>"
	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Название ✏️",
			Unique: static.CallbackActionEditUnsavedCheck.String(),
			Data:   static.CallbackEditCheckName,
		},
		Button{
			BtnTxt: "Дату 📆",
			Unique: static.CallbackActionEditUnsavedCheck.String(),
			Data:   static.CallbackEditCheckCreationDate,
		},
		Button{
			BtnTxt: "Кто заплатил 🧑‍🤝‍🧑",
			Unique: static.CallbackActionEditUnsavedCheck.String(),
			Data:   static.CallbackEditCheckOwner,
		},
		Button{
			BtnTxt: "Изменить товары 📝",
			Unique: static.CallbackActionEditUnsavedCheck.String(),
			Data:   static.CallbackEditCheckItems,
		},
		Button{
			BtnTxt: "Назад ⬅️",
			Unique: static.CallbackActionEditUnsavedCheck.String(),
			Data:   static.CallbackMenuGoBack,
		},
	)

	return text, kb
}

const (
	RETRY       = true
	NOT_A_RETRY = false
)

func GetAskForNewCheckNameResponse(currentCheckName string, action static.CallbackAction, isRetry bool) (string, *tele.ReplyMarkup) {
	var text string

	if isRetry {
		text += "<b><u>Ошибка</u></b>\n" +
			"<i>Название не может быть пустым, должно содержать буквы и не превышать 1000 символов.</i>\n" +
			"Попробуй еще раз.\n\n"
	}

	text += fmt.Sprintf(
		"<b><u>Текущее название этого чека:</u> %s</b>\n\nНапиши новое название. Если передумал менять название, просто нажми на кнопку <u>\"Назад ⬅️\"</u>",
		currentCheckName,
	)

	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Назад ⬅️",
			Unique: action.String(),
			Data:   static.CallbackMenuGoBack,
		},
	)

	return text, kb
}

func GetNewCheckNameIsSavedResponse(checkName string) (string, *tele.ReplyMarkup) {
	text := fmt.Sprintf(
		"<b>Название чека изменено на <i>%s</i>!</b>",
		checkName,
	)

	return text, nil
}

func GetAskForCheckOwnershipQuestion(action static.CallbackAction) (string, *tele.ReplyMarkup) {
	text := "<b>С чьей карты был оплачен чек?💳💵</b>"

	rows := []RowOfButtons{
		{ // First row
			BtnsPerRow: 2,
			Btns: []Button{
				{
					BtnTxt: "Liz 💜",
					Unique: action.String(),
					Data:   static.CallbackOwnerLiz,
				},
				{
					BtnTxt: "Pau 🩵",
					Unique: action.String(),
					Data:   static.CallbackOwnerPau,
				},
			},
		},
		{ // Second row
			BtnsPerRow: 1,
			Btns: []Button{
				{
					BtnTxt: "С общей 💜🩵",
					Unique: action.String(),
					Data:   static.CallbackOwnerBoth,
				},
			},
		},
	}

	if action == static.CallbackActionEditUnsavedCheck {
		rows = append(rows,
			RowOfButtons{
				BtnsPerRow: 1,
				Btns: []Button{
					{
						BtnTxt: "Назад ⬅️",
						Unique: static.CallbackActionEditUnsavedCheck.String(),
						Data:   static.CallbackMenuGoBack,
					},
				},
			},
		)
	}

	return text, createCustomRowsInlineKb(rows...)
}

func GetAskForNewCheckCreationDateQuestion() (string, *tele.ReplyMarkup) {
	text := "<b><u>Изменение даты и времени</u></b>\n\n"
	text += "Укажите новую дату и время в формате: <i>ГГГГ-ММ-ДД ЧЧ:ММ</i>"

	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Назад ⬅️",
			Unique: static.CallbackActionEditUnsavedCheck.String(),
			Data:   static.CallbackMenuGoBack,
		},
	)

	return text, kb
}

func GetShowItemForEditResponse(item *static.Item, currentIndex, outOf int) (string, *tele.ReplyMarkup) {
	text := fmt.Sprintf(
		"<b>Товар: %d из %d</b>\n\n",
		currentIndex+1,
		outOf,
	)
	text += sPrintItemInfo(item)

	rows := []RowOfButtons{
		{ // First row
			BtnsPerRow: 3,
			Btns: []Button{
				{
					BtnTxt: "⬅️",
					Unique: static.CallbackActionEditUnsavedItem.String(),
					Data:   static.CallbackMenuGoBackward,
				},
				{
					BtnTxt: "Изменить ✏️",
					Unique: static.CallbackActionEditUnsavedItem.String(),
					Data:   static.CallbackOwnerPau,
				},
				{
					BtnTxt: "➡️",
					Unique: static.CallbackActionEditUnsavedItem.String(),
					Data:   static.CallbackMenuGoForward,
				},
			},
		},
		{ // Second row
			BtnsPerRow: 1,
			Btns: []Button{
				{
					BtnTxt: "Назад к чеку ⬅️",
					Unique: static.CallbackActionEditUnsavedItem.String(),
					Data:   static.CallbackMenuGoBack,
				},
			},
		},
	}

	return text, createCustomRowsInlineKb(rows...)
}

func GetShowItemEditOptions(item *static.Item, action static.CallbackAction) (string, *tele.ReplyMarkup) {
	text := "<b>Изменение товара</b>\n\n"
	text += sPrintItemInfo(item) + "\n\n<b>Что меняем?</b>👀"
	kb := createSelectorInlineKb(
		1,
		Button{
			BtnTxt: "Название",
			Unique: action.String(),
			Data:   static.CallbackEditItemName,
		},
		Button{
			BtnTxt: "Цена",
			Unique: action.String(),
			Data:   static.CallbackEditItemPrice,
		},
		Button{
			BtnTxt: "Количество",
			Unique: action.String(),
			Data:   static.CallbackEditItemAmount,
		},
		Button{
			BtnTxt: "Чей товар",
			Unique: action.String(),
			Data:   static.CallbackEditItemOwner,
		},
		Button{
			BtnTxt: "Вернуться ⬅️",
			Unique: action.String(),
			Data:   static.CallbackMenuGoBack,
		},
	)

	return text, kb
}
