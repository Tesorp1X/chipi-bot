package util

/*
	Html-codes:
	< &lt;
	> &gt;
	— &#8212;

*/
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
		msg += "<i>" + strconv.Itoa(no) + ") " + item.Name + " " + strconv.FormatFloat(item.Price, 'f', 2, 64) + " руб</i>\n"

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

	msg += "\nЛиз заплатила: <b>" + strconv.FormatFloat(sumL, 'f', 2, 64) + " руб</b>\n"
	msg += "Пау заплатил: <b>" + strconv.FormatFloat(sumP, 'f', 2, 64) + " руб</b>\n\n"
	msg += "Итого: <b>" + strconv.FormatFloat(total, 'f', 2, 64) + " бублей.</b>"

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
		switch item.Owner {
		case ct.OwnerId:
			ct.OwnerTotal += item.Price
		case models.OWNER_BOTH:
			ct.OwnerTotal += item.Price / 2
			ct.DebtorTotal += item.Price / 2
		default:
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

	st := &models.SessionTotal{
		SessionId: sessionId,
		Total:     lizTotal.All + pauTotal.All,
		TotalLiz:  lizTotal.Owner + pauTotal.Debtor,
		TotalPau:  pauTotal.Owner + lizTotal.Debtor,
	}

	if lizTotal.Debtor > pauTotal.Debtor {
		st.Recipient = models.OWNER_LIZ
		st.Amount = lizTotal.Debtor - pauTotal.Debtor
	} else {
		st.Recipient = models.OWNER_PAU
		st.Amount = pauTotal.Debtor - lizTotal.Debtor
	}

	return st
}

// Returns a response based on [models.SessionTotal].
// Make [isPreliminary] true, if you wish to get preliminary results message.
func GetTotalResponse(sessionTotal *models.SessionTotal, isPreliminary bool) string {
	var msg string
	if isPreliminary {
		msg = fmt.Sprintf("<b>Вот промежуточный итог за этот период:</b>\n\n<b><u>Всего заплачено: %.2f руб</u></b>\n", sessionTotal.Total)
	} else {
		msg = fmt.Sprintf("<b>Вот итог за этот период:</b>\n\n<b><u>Всего заплачено: %.2f руб</u></b>\n", sessionTotal.Total)
	}

	if sessionTotal.Recipient == models.OWNER_LIZ {
		msg += fmt.Sprintf("<i>Пау должен Лиз <b>%.2f руб.</b></i>", sessionTotal.Amount)
	} else {
		msg += fmt.Sprintf("<i>Лиз должна Пау <b>%.2f руб.</b></i>", sessionTotal.Amount)
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
	msg := "<b>Чек:</b> <i>" + check.GetCheckName() + "</i>\n"

	switch check.GetCheckOwner() {
	case models.OWNER_LIZ:
		msg += "<b>Заплачено:</b> <i>Лиз :3 </i>\n\n"
	case models.OWNER_PAU:
		msg += "<b>Заплачено:</b> <i>Пау &lt;3 </i>\n\n"
	}

	lizItems := "<b><u>Товары Лиз</u></b>💜:\n"
	pauItems := "<b><u>Товары Пау</u></b>💙:\n"
	mutualItems := "<b><u>Общие товары</u></b>💜💙:\n"
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
			lizItems += "<i>" + strconv.Itoa(lizCount) + ") " + item.Name + " — <b>" + strconv.FormatFloat(item.Price, 'f', 2, 64) + "</b></i>\n"
			lizSum += item.Price
			lizCount++
		case models.OWNER_PAU:
			pauItems += "<i>" + strconv.Itoa(pauCount) + ") " + item.Name + " — <b>" + strconv.FormatFloat(item.Price, 'f', 2, 64) + "</b></i>\n"
			pauSum += item.Price
			pauCount++
		case models.OWNER_BOTH:
			mutualItems += "<i>" + strconv.Itoa(mutualCount) + ") " + item.Name + " — <b>" + strconv.FormatFloat(item.Price, 'f', 2, 64) + "</b></i>\n"
			mutualSum += item.Price
			mutualCount++
		}

		total += item.Price
	}

	msg += lizItems
	if lizCount > 1 {
		msg += "\n" + "<i>Товаров Лиз на:</i> <b>" + strconv.FormatFloat(lizSum, 'f', 2, 64) + "</b> <i>руб</i>\n\n"
	} else {
		msg += "<i>товаров нет</i>" + "\n\n"
	}

	msg += pauItems
	if pauCount > 1 {
		msg += "\n" + "<i>Товаров Пау на:</i> <b>" + strconv.FormatFloat(pauSum, 'f', 2, 64) + "</b> <i>руб</i>\n\n"
	} else {
		msg += "<i>товаров нет</i>" + "\n\n"
	}

	msg += mutualItems
	if mutualCount > 1 {
		msg += "\n" + "<i>Общих товаров на:</i> <b>" + strconv.FormatFloat(mutualSum, 'f', 2, 64) + "</b> <i>руб</i>\n\n"
	} else {
		msg += "<i>товаров нет</i>" + "\n\n"
	}

	msg += "<i>Итого:</i> <b>" + strconv.FormatFloat(total, 'f', 2, 64) + "</b> <i>бублей.</i>"

	return msg
}

// Response for '/show totals' command. Show one at a time.
func GetShowTotalsResponse(sessionTotal *models.SessionTotal) string {
	msg := "<b>Результат сессии №" + strconv.FormatInt(sessionTotal.GetSessionId(), 10) + ":</b>\n\n"

	msg += "<i><b>Дата начала:</b> " + sessionTotal.GetOpenedAtTime().Format(time.DateTime) + "</i>\n"
	msg += "<i><b>Дата окончания:</b> " + sessionTotal.GetClosedAtTime().Format(time.DateTime) + "</i>\n\n"

	if sessionTotal.Recipient == models.OWNER_LIZ {
		msg += fmt.Sprintf("Пау перевел Лиз <b>%.2f руб.</b>\n\n", sessionTotal.Amount)
	} else {
		msg += fmt.Sprintf("Лиз перевела Пау <b>%.2f руб.</b>\n\n", sessionTotal.Amount)
	}

	msg += fmt.Sprintf("<b><i>Лиз купила на: %.2f руб</i></b>\n", sessionTotal.TotalLiz)
	msg += fmt.Sprintf("<b><i>Пау купил на: %.2f руб</i></b>\n\n", sessionTotal.TotalPau)
	msg += fmt.Sprintf("<b><u>Всего заплачено: %.2f руб</u></b>", sessionTotal.Total)

	return msg
}
/*
[ParsePrice] looks for valid item price response patterns
and returns converted value and an error. If a valid pattern was found,
then it will be converted according to context (see examples).
If there were no matches, 0.0 with an error will be returned.
Valid patterns are:
- Single value (45 or 45.05 or 45,05);
- Price with multiplier (2*45 or 2*45.05 or 2*45,05 or 2 * 45 or 2 * 45.05 or 2 * 45,05).
In second case price is on the right and will be multiplied by quantity on left and returned afterwards.
All other cases will result an error.
*/
func ParsePrice(itemPriceStr string) (float64, error) {
	convertedPrice := 0.0
	if strings.Contains(itemPriceStr, "*") {
		// price with multiplier case
		tokens := strings.Split(itemPriceStr, "*")
		if len(tokens) != 2 {
			return convertedPrice, models.ErrItemPriceWrongFormat
		}

		if !verifyPrice(tokens[1]) {
			return convertedPrice, models.ErrItemPriceNotSingleIntOrFloat
		}

		if !verifyPriceMultiplier(tokens[0]) {
			return convertedPrice, models.ErrItemPriceMultiplierNotSingleInt
		}
		// todo: clear trailing spaces
		multiplier, _ := strconv.Atoi(tokens[0])
		price, _ := strconv.ParseFloat(tokens[1], 64)

		convertedPrice = float64(multiplier) * price
	}
	return convertedPrice, nil
}
