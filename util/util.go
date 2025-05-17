package util

import (
	"strings"

	"github.com/Tesorp1X/chipi-bot/models"
)

func ExtractDataFromCallback(data string, action models.CallbackAction) string {
	return strings.TrimPrefix(data, "\f"+action.String()+"|")
}
