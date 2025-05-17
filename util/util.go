package util

import (
	"strings"

	"github.com/Tesorp1X/chipi-bot/models"
)

// Extracts data from [Callback.Data] by removing prefix '\f + CallbackAction + |'
func ExtractDataFromCallback(data string, action models.CallbackAction) string {
	return strings.TrimPrefix(data, "\f"+action.String()+"|")
}
