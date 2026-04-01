package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
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

const (
	EnglishAlphabet = "abcdefghijklmnopqrstuvwxyz"
	RussianAlphabet = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
)

func VerifyName(messageText string) bool {
	// length
	fact := len(messageText) > 0 && len(messageText) < 1000

	// contains letters
	fact = fact && strings.ContainsAny(
		strings.ToLower(messageText),
		EnglishAlphabet+RussianAlphabet,
	)
	// todo: perhaps add return error to specify a problem with a string
	return fact
}

const isResponded = "is_callback_query_responded"

func IsCbQueryResponded(c tele.Context) bool {
	// false or absence of value in ctx storage means not responded
	if val := c.Get(isResponded); val != nil && val == true {
		return true
	}

	return false
}

func MarkCbQueryAsResponded(c tele.Context) error {
	if IsCbQueryResponded(c) {
		return errors.New("error in utils.MarkCbQueryAsResponded(): a callback query is already responded")
	}

	c.Set(isResponded, true)

	return nil
}
