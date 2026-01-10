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

	targetFilePath, err := downloadFile(c, conf.DownloadPath, d.FileID, d.FileName)
	if err != nil {
		return fmt.Errorf(
			"error in OnDocumentActionHandler(): couldn't download a file (%v)",
			err,
		)
	}

	checkData, err := reader.ExtractCheckData(targetFilePath)
	if err != nil {
		sendErr := c.Send("error: couldn't extract data")

		return fmt.Errorf(
			"error in OnDocumentActionHandler(): could't extract data from file '%s' (%v). send with error: %v",
			targetFilePath,
			err,
			sendErr,
		)
	}

	return c.Send(checkData.OrgName)
}

func downloadFile(c tele.Context, downloadDirPath string, fileId string, fileName string) (string, error) {
	f, err := c.Bot().FileByID(fileId)
	if err != nil {
		sendErr := c.Send("error with 'c.Bot().FileByID': " + err.Error())
		return "", fmt.Errorf(
			"error in downloadFile(): couldn't find a file with id %s (%v), message send with err: %v",
			fileId,
			err,
			sendErr,
		)
	}

	targetFilePath := downloadDirPath + `\` + fileName
	err = c.Bot().Download(&f, targetFilePath)
	if err != nil {
		sendErr := c.Send("error: couldn't download a file: " + err.Error())
		return "", fmt.Errorf(
			"error in downloadFile(): couldn't download a file {id: %s; name: %s; path: %s} (%v), message send with err: %v",
			fileId,
			fileName,
			targetFilePath,
			err,
			sendErr,
		)
	}

	return targetFilePath, nil
}
