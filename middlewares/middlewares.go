package middlewares

import (
	"log"

	"github.com/Tesorp1X/chipi-bot/db"
	"github.com/Tesorp1X/chipi-bot/models"
	tele "gopkg.in/telebot.v4"
)

func AutoSessionAssigner(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		sessionId, err := db.GetSessionId()
		if err != nil {
			log.Fatal(err)
		}
		c.Set(models.SESSION_ID, sessionId)
		return next(c) // continue execution chain
	}
}
