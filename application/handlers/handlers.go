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
	"github.com/Tesorp1X/chipi-bot/utils"

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
	case static.StateWaitForCheckName, static.StateWaitForNewCheckNameUnsaved:
		if err := handleCheckName(c, state); err != nil {
			return fmt.Errorf(
				"error in handlers.HandleAnyText(), state '%s': %v",
				currentState,
				err,
			)
		}
	case static.StateWaitForCheckCreationDateUnsaved:
		if err := handleEditCheckCreationDate(c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForCheckCreationDate': %v", err)
		}
	}

	return nil
}

func handleCheckName(c tele.Context, state fsm.Context) error {
	if !utils.VerifyName(c.Message().Text) {
		// retry prompt
		if err := prompts.SendRetryCheckNameMessage(c, state); err != nil {
			return fmt.Errorf(
				"error in handlers.handleCheckName(): prompt failed (%v)",
				err,
			)
		}

		return nil
	}

	// name is okay
	_, err := storageHelpers.SetNewCheckNameFromMessage(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't set new check name (%v)",
			err,
		)
	}

	currentState, err := state.State(context.Background())
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): failed to retrieve state (%v)",
			err,
		)
	}

	var promptErr error
	switch currentState {
	case static.StateWaitForCheckName:
		// move on with the verification process: ask for owner
		promptErr = prompts.SendCheckOwnershipMessage(prompts.FromAddCheck, c, state)
	case static.StateWaitForNewCheckNameUnsaved:
		// come back to check edit menu
		promptErr = prompts.SendEditUnsavedCheckMessage(c, state)
	default:
		return fmt.Errorf(
			"error in handlers.handleCheckName(): invalid state for check name change (%s)",
			currentState,
		)
	}

	if promptErr != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): prompt failed (%v)",
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

	if err := prompts.SendEditUnsavedCheckMessage(c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckCreationDate(): failed to send a check-edit message (%v)",
			err,
		)
	}

	return nil
}
