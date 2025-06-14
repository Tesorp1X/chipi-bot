package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Tesorp1X/chipi-bot/models"
)

// Extracts data from [Callback.Data] by removing prefix '\f + CallbackAction + |'
func ExtractDataFromCallback(data string, action models.CallbackAction) string {
	return strings.TrimPrefix(data, "\f"+action.String()+"|")
}

func CreateItemsListResponse(itemsList ...models.Item) string {
	var (
		msg   string
		no    int
		total float64
		sumL  float64
		sumP  float64
	)

	for i, item := range itemsList {
		no = i + 1
		msg += strconv.Itoa(no) + ") " + item.Name + " " + strconv.FormatFloat(item.Price, 'f', 2, 64) + " руб\n"

		total += item.Price
		switch item.Owner {
		case models.OWNER_LIZ:
			sumL += item.Price
		case models.OWNER_PAU:
			sumP += item.Price
		default:
			sumL += item.Price / 2
			sumP += item.Price / 2
		}
	}

	msg += "Лиз заплатила: " + strconv.FormatFloat(sumL, 'f', 2, 64) + " руб\n"
	msg += "Пау заплатил: " + strconv.FormatFloat(sumP, 'f', 2, 64) + " руб\n"
	msg += "Итого: " + strconv.FormatFloat(total, 'f', 2, 64) + " бублей."

	return msg
}

func ExtractAdminsIDs(adminsStr string) []int64 {
	adminsStr = strings.ReplaceAll(adminsStr, " ", "")
	adminsStr = strings.ReplaceAll(adminsStr, "[", "")
	adminsStr = strings.ReplaceAll(adminsStr, "]", "")
	admins := strings.Split(adminsStr, ",")
	var res []int64
	for _, s := range admins {
		n, _ := strconv.ParseInt(s, 10, 64)
		res = append(res, n)
	}

	return res
}

func CalculateCheckTotal(check *models.CheckWithItems) *models.CheckTotal {
	ct := &models.CheckTotal{Id: check.Id, OwnerId: check.GetCheckOwner()}
	for _, item := range check.GetItems() {
		if item.Owner == ct.OwnerId {
			ct.OwnerTotal += item.Price
		} else if item.Owner == models.OWNER_BOTH {
			ct.OwnerTotal += item.Price / 2
			ct.DebtorTotal += item.Price / 2
		} else {
			ct.DebtorTotal += item.Price
		}
		ct.Total += item.Price
	}
	return ct
}

type total struct {
	Owner  float64
	Debtor float64
	All    float64
}

func CalculateSessionTotal(sessionId int64, checks []*models.CheckWithItems) *models.SessionTotal {

	lizTotal := total{}
	pauTotal := total{}

	for _, check := range checks {
		t := CalculateCheckTotal(check)
		if check.GetCheckOwner() == models.OWNER_LIZ {
			lizTotal.Owner += t.OwnerTotal
			lizTotal.Debtor += t.DebtorTotal
			lizTotal.All += t.Total
		} else {
			pauTotal.Owner += t.OwnerTotal
			pauTotal.Debtor += t.DebtorTotal
			pauTotal.All += t.Total
		}
	}

	st := &models.SessionTotal{SessionId: sessionId, Total: lizTotal.All + pauTotal.All}
	if lizTotal.Debtor > pauTotal.Debtor {
		st.Recipient = models.OWNER_LIZ
		st.Amount = lizTotal.Debtor - pauTotal.Debtor
	} else {
		st.Recipient = models.OWNER_PAU
		st.Amount = pauTotal.Debtor - lizTotal.Debtor
	}

	return st
}

// Returns a responce based on [models.SessionTotal].
// Make [isPreliminary] true, if you wish to get preliminary results message.
func GetTotalResponse(sessionTotal *models.SessionTotal, isPreliminary bool) string {
	var msg string
	if isPreliminary {
		msg = fmt.Sprintf("Вот промежуточный итог за этот период:\nВсего заплачено: %.2f руб\n", sessionTotal.Total)
	} else {
		msg = fmt.Sprintf("Вот итог за этот период:\nВсего заплачено: %.2f руб\n", sessionTotal.Total)
	}

	if sessionTotal.Recipient == models.OWNER_LIZ {
		msg += fmt.Sprintf("Пау должен Лиз %.2f руб.", sessionTotal.Amount)
	} else {
		msg += fmt.Sprintf("Лиз должна Пау %.2f руб.", sessionTotal.Amount)
	}

	return msg
}

func GetCheckCreatedResponse(checkOwner string) string {
	msg := "Чек создан!😇\n"
	switch checkOwner {
	case models.OWNER_LIZ:
		msg += "Заплатила Лиз💜\n"
	case models.OWNER_PAU:
		msg += "Заплатил Пау💙\n"
	}
	msg += "Теперь давай добавим покупочки😋\n\n"
	msg += "Название товара?👀"

	return msg
}

func GetItemAdded(itemOwner string, itemPrice float64) string {
	msg := "Товар добавлен.\n"
	switch itemOwner {
	case models.OWNER_LIZ:
		msg += "Заплатила Лиз💜\n"
	case models.OWNER_PAU:
		msg += "Заплатил Пау💙\n"
	case models.OWNER_BOTH:
		msg += "Товар общий💜💙\n"
	}
	msg += "Цена: " + strconv.FormatFloat(itemPrice, 'f', 2, 64) + "\n\n"
	msg += "Еще товары?"

	return msg
}

func GetCheckWithItemsResponse(check models.CheckWithItems) string {
	msg := "Чек: " + check.GetCheckName() + " "
	msg += "Заплачено: " + check.GetCheckOwner() + "\n\n"

	lizItems := "Товары Лиз💜:\n"
	pauItems := "Товары Пау💙:\n"
	mutualItems := "Общие товары💜💙:\n"
	lizCount := 1
	pauCount := 1
	mutualCount := 1

	var (
		total     float64
		lizSum    float64
		pauSum    float64
		mutualSum float64
	)

	for _, item := range check.GetItems() {
		switch item.Owner {
		case models.OWNER_LIZ:
			lizItems += strconv.Itoa(lizCount) + ") " + item.Name + " -- " + strconv.FormatFloat(item.Price, 'f', 2, 64) + "\n"
			lizSum += item.Price
			lizCount++
		case models.OWNER_PAU:
			pauItems += strconv.Itoa(pauCount) + ") " + item.Name + " -- " + strconv.FormatFloat(item.Price, 'f', 2, 64) + "\n"
			pauSum += item.Price
			pauCount++
		case models.OWNER_BOTH:
			mutualItems += strconv.Itoa(mutualCount) + ") " + item.Name + " -- " + strconv.FormatFloat(item.Price, 'f', 2, 64) + "\n"
			mutualSum += item.Price
			mutualCount++
		}

		total += item.Price
	}

	msg += lizItems
	if lizCount > 1 {
		msg += "\n" + "Товаров Лиз на: " + strconv.FormatFloat(lizSum, 'f', 2, 64) + " руб\n\n"
	} else {
		msg += "товаров нет" + "\n\n"
	}

	msg += pauItems
	if pauCount > 1 {
		msg += "\n" + "Товаров Пау на: " + strconv.FormatFloat(pauSum, 'f', 2, 64) + " руб\n\n"
	} else {
		msg += "товаров нет" + "\n\n"
	}

	msg += mutualItems
	if mutualCount > 1 {
		msg += "\n" + "Общих товаров на: " + strconv.FormatFloat(mutualSum, 'f', 2, 64) + " руб\n\n"
	} else {
		msg += "товаров нет" + "\n\n"
	}

	msg += "Итого: " + strconv.FormatFloat(total, 'f', 2, 64) + " бублей."

	return msg
}

// Responce for '/show totals' command. Show one at a time.
func GetShowTotalsResponse(total *models.SessionTotal) string {
	msg := "Результат сессии №" + strconv.FormatInt(total.GetSessionId(), 10) + ":\n\n"

	msg += "Дата начала: " + total.GetOpenedAtTime().Format(time.DateTime) + "\n"
	msg += "Дата окончания: " + total.GetClosedAtTime().Format(time.DateTime) + "\n\n"

	if total.Recipient == models.OWNER_LIZ {
		msg += fmt.Sprintf("Пау перевел Лиз %.2f руб.", total.Amount)
	} else {
		msg += fmt.Sprintf("Лиз перевела Пау %.2f руб.", total.Amount)
	}

	return msg
}
