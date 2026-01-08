package reader

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type CheckItem struct {
	// Пример
	// 'Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый] 1.0 18 699,00 18 699,00'
	rawText string
	// 'Пылесос-робот Dreame D20 [Лидар, 5200 мА*ч, 13000Па, уборка: влаж, сух, 350 мл, 0.7 л, смарт упр, картография, белый]'
	Name string
	// 18699.00
	Price float64
	// 1.0
	Amount float64
	// 18699.00
	SubTotal float64
}

func (a *CheckItem) IsEqual(b *CheckItem) bool {
	return a.Name == b.Name && a.Amount == b.Amount &&
		a.Price == b.Price && a.SubTotal == b.SubTotal
}

func NewCheckItem(rawItemText string) (CheckItem, error) {
	subTotal, err := extractPrice(rawItemText, SUBTOTAL)
	if err != nil {
		return CheckItem{}, fmt.Errorf("error in NewCheckIte: extractPrice('%s', SUBTOTAL) finished with an error: %v.", rawItemText, err)
	}

	price, err := extractPrice(rawItemText, PRICE)
	if err != nil {
		return CheckItem{}, fmt.Errorf("error in NewCheckIte: extractPrice('%s', PRICE) finished with an error: %v.", rawItemText, err)
	}

	name, err := extractItemName(rawItemText)
	if err != nil {
		return CheckItem{}, fmt.Errorf("error in NewCheckIte: extractItemName('%s') finished with an error: %v.", rawItemText, err)
	}

	return CheckItem{
		rawText:  rawItemText,
		Price:    price,
		SubTotal: subTotal,
		Name:     name,
		Amount:   round(subTotal/price, 3),
	}, nil
}

// analogue to python's round
func round(val float64, precision int) float64 {
	return math.Round(val*(math.Pow10(precision))) / math.Pow10(precision)
}

const (
	SUBTOTAL = iota
	PRICE
)

func extractPrice(rawItemText string, pos int) (float64, error) {
	var idx int
	var priceStr string
	switch pos {
	case SUBTOTAL:
		idx = strings.LastIndex(rawItemText[:len(rawItemText)-3], ",") + 4
		priceStr = rawItemText[idx:]
	case PRICE:
		idx = strings.LastIndex(rawItemText[:len(rawItemText)-3], ",") + 4
		s := rawItemText[:idx]
		dotIdx := strings.LastIndex(s, ".")
		spaceIdx := strings.Index(s[dotIdx:], " ") + 1
		priceStr = s[dotIdx:][spaceIdx:]
	default:
		return -1, fmt.Errorf("error: invalid 'pos' option: %v", pos)
	}

	price, err := strconv.ParseFloat(normalizePriceStr(priceStr), 32)

	return round(price, 2), err
}

func normalizePriceStr(rawPriceStr string) string {

	return strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(rawPriceStr, ",", "."),
			" ", ""),
	)
}

func extractItemName(rawItemText string) (string, error) {
	dotIdx := strings.LastIndex(rawItemText, ".")
	if dotIdx == -1 {
		return "", fmt.Errorf("error: in extractItemName. couldn't find a '.' symbol in string: '%s'", rawItemText)
	}
	rawItemText = rawItemText[:dotIdx]
	spaceIdx := strings.LastIndex(rawItemText, " ")
	if spaceIdx == -1 {
		return "", fmt.Errorf("error: in extractItemName. couldn't find a ' ' symbol in string: '%s'", rawItemText)
	}

	return strings.TrimSpace(rawItemText[:spaceIdx]), nil
}
