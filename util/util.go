package util

import (
	"strconv"
	"strings"

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
			sumL += float64(item.Price)
		case models.OWNER_PAU:
			sumP += float64(item.Price)
		default:
			sumL += float64(item.Price) / 2
			sumP += float64(item.Price) / 2
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
