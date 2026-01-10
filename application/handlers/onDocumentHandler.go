package handlers

import (
	"fmt"

	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/utils/reader"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func OnDocumentActionHandler(conf *config.Config, c tele.Context, state fsm.Context) error {
	d := c.Message().Document
	if d == nil {
		return c.Send("error: no file")
	}
	f, err := c.Bot().FileByID(d.FileID)
	if err != nil {
		return c.Send("error with  'c.Bot().FileByID': " + err.Error())
	}

	targetFilePath := conf.DownloadPath + `\` + d.FileName
	err = c.Bot().Download(&f, targetFilePath)
	if err != nil {
		return c.Send("error: couldn't download a file: " + err.Error())
	}

	checkData, err := reader.ExtractCheckData(targetFilePath)
	if err != nil {
		wrappedErr := fmt.Errorf(
			"error in OnDocumentActionHandler(), file '%s': %v",
			targetFilePath,
			err,
		)
		if err := c.Send(wrappedErr.Error()); err != nil {
			return fmt.Errorf(
				wrappedErr.Error()+"\n"+
					"error in OnDocumentActionHandler(), couldn't send a message: %v",
				err,
			)
		}
		return wrappedErr
	}

	return c.Send(checkData.OrgName)
}
