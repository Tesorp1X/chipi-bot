package handlers

import (
	"fmt"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/application/prompts"
	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/reader"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func OnDocumentActionHandler(conf *config.Config, c tele.Context, state fsm.Context) error {
	d := c.Message().Document
	if d == nil {
		return c.Send("error: no file")
	}

	// verify that filename doesn't exist already
	// ask to

	// todo: log an error, but don't interrupt
	c.Send("Скачиваю...📡💾")

	targetFilePath, err := downloadFile(c, conf.DownloadPath, d.FileID, d.FileName)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.OnDocumentActionHandler(): couldn't download a file (%v)",
			err,
		)
	}

	// todo: log an error, but don't interrupt
	c.Send("Анализирую...🤔🔬🔍")

	checkData, err := reader.ExtractCheckData(targetFilePath)
	if err != nil {
		sendErr := c.Send("error: couldn't extract data")

		return fmt.Errorf(
			"error in handlers.OnDocumentActionHandler(): could't extract data from file '%s' (%v). send with error: %v",
			targetFilePath,
			err,
			sendErr,
		)
	}

	// convert checkData to check obj and save it in context
	check := static.CreateCheckFromCheckData(checkData)
	if err := storageHelpers.UpdateCheck(check, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.OnDocumentActionHandler(): couldn't save check in state-storage (%v)",
			err,
		)
	}

	// save items in context
	if err := storageHelpers.UpdateItemsList(static.CreateItemsFromCheckData(checkData), c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.OnDocumentActionHandler(): couldn't save items in state-storage (%v)",
			err,
		)
	}

	if err := prompts.SendNewCheckNameQuestionMessage(prompts.FromAddCheck, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.OnDocumentActionHandler(): couldn't send a 'new check name' prompt (%v)",
			err,
		)
	}

	return nil
}

func downloadFile(c tele.Context, downloadDirPath string, fileId string, fileName string) (string, error) {
	f, err := c.Bot().FileByID(fileId)
	if err != nil {
		sendErr := c.Send("error with 'c.Bot().FileByID': " + err.Error())
		return "", fmt.Errorf(
			"error in handlers.downloadFile(): couldn't find a file with id %s (%v), message sent with err: %v",
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
			"error in handlers.downloadFile(): couldn't download a file {id: %s; name: %s; path: %s} (%v), message sent with err: %v",
			fileId,
			fileName,
			targetFilePath,
			err,
			sendErr,
		)
	}

	return targetFilePath, nil
}
