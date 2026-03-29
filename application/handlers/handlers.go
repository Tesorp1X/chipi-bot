package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/application/prompts"
	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func HandleStartCommand(c tele.Context, state fsm.Context) error {
	return c.Send("hello")
}

func HandleCancelCommand(c tele.Context, state fsm.Context) error {
	if err := storageHelpers.FinishState(static.DELETE_DATA, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.HandleCancelCommand(): failed to finish state (%v)",
			err,
		)
	}

	c.Send("Canceled! All data removed.")

	return nil
}

func HandleAnyText(conf *config.Config, c tele.Context, state fsm.Context) error {
	currentState, err := state.State(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error in handlers.HandleAnyCallback(): couldn't receive users(%d) current state: %v",
			c.Sender().ID, err,
		)
	}

	switch currentState {
	case static.StateWaitForCheckName:
		if err := handleCheckName(c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForCheckName': %v", err)
		}
	case static.StateWaitForNewCheckNameUnsaved:
		if err := handleEditCheckName(c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForNewCheckNameUnsaved': %v", err)
		}
	case static.StateWaitForCheckCreationDate:
		if err := handleEditCheckCreationDate(c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForCheckCreationDate': %v", err)
		}
	}

	return nil
}

func handleCheckName(c tele.Context, state fsm.Context) error {
	_, err := storageHelpers.SetNewCheckNameFromMessage(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't set new check name (%v)",
			err,
		)
	}

	if err := prompts.SendCheckOwnershipMessage(prompts.FromAddCheck, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): failed to send a check ownership message (%v)",
			err,
		)
	}

	return nil
}

func handleEditCheckName(c tele.Context, state fsm.Context) error {
	check, err := storageHelpers.SetNewCheckNameFromMessage(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't set new check name (%v)",
			err,
		)
	}

	if err := c.Send(responses.GetNewCheckNameIsSavedResponse(check.Name)); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't send a message (%v)",
			err,
		)
	}

	if err := prompts.SendEditCheckMessage(c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): failed to send a check-edit message (%v)",
			err,
		)
	}

	return nil
}

func handleEditCheckCreationDate(c tele.Context, state fsm.Context) error {
	gotDateStr := c.Message().Text + ":00"
	newDate, errTime := time.Parse(time.DateTime, gotDateStr)
	if errTime != nil {
		errMsg := fmt.Sprintf(
			"error in handlers.handleEditCheckCreationDate():\nfailed to parse a given date-time '%s' (%v)\n",
			gotDateStr,
			errTime,
		)
		// todo: specify error maybe
		sendErr := c.Send("error: wrong date-time format: " + errTime.Error())
		if sendErr != nil {
			errMsg += fmt.Sprintf("couldn't send a message (%v)\n", sendErr)
		}

		return errors.New(errMsg)
	}

	// date is fine
	check, errCheck := storageHelpers.GetCheck(c, state)
	if errCheck != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckCreationDate(): failed to retrieve a check (%v)",
			errCheck,
		)
	}

	check.Date = &newDate
	if err := storageHelpers.UpdateCheck(check, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckCreationDate(): failed to update a check (%v)",
			err,
		)
	}

	if err := prompts.SendEditCheckMessage(c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckCreationDate(): failed to send a check-edit message (%v)",
			err,
		)
	}

	return nil
}
