package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ExtractAdminsIDs(adminsStr string) ([]int64, error) {
	adminsStr = strings.ReplaceAll(adminsStr, " ", "")
	adminsStr = strings.ReplaceAll(adminsStr, "[", "")
	adminsStr = strings.ReplaceAll(adminsStr, "]", "")
	admins := strings.Split(adminsStr, ",")
	var res []int64
	for _, s := range admins {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"error in utils.ExtractAdminsIDs(): couldn't parse '%s': %v",
				s,
				err,
			)
		}
		res = append(res, n)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf(
			"error in utils.ExtractAdminsIDs(): got list of length 0 from '%s'",
			adminsStr,
		)
	}

	return res, nil
}

// Returns Callback data part, that comes after '|' symbol.
func ExtractCallbackData(rawData string) string {
	idx := strings.IndexRune(rawData, '|') + 1
	return rawData[idx:]
}

func VerifyName(messageText string) bool {
	// length
	fact := len(messageText) > 0 && len(messageText) < 1000

	// contains letters
	fact = fact && strings.ContainsAny(messageText, "qwertyuiopasdfghjklzxcvbnm")

	return fact
}
