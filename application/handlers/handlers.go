package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	storageHelpers "github.com/Tesorp1X/chipi-bot/application/StorageHelpers"
	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/Tesorp1X/chipi-bot/utils/responses"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

func HandleStartCommand(c tele.Context, state fsm.Context) error {
	return c.Send("hello")
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
	case static.StateWaitForNewCheckName:
		if err := handleEditCheckName(c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForNewCheckName': %v", err)
		}
	case static.StateWaitForCheckCreationDate:
		if err := handleEditCheckCreationDate(c, state); err != nil {
			return fmt.Errorf("error in handlers.HandleAnyText(), state 'StateWaitForNewCheckName': %v", err)
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

	if err := storageHelpers.SetState(static.StateWaitForCheckOwner, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't change a state (%v)",
			err,
		)
	}
	// prompt check ownership
	if sendErr := c.Send(responses.GetAskForCheckOwnershipQuestion()); sendErr != nil {
		return fmt.Errorf(
			"error in handlers.handleCheckName(): couldn't send a 'check-ownership'-message (%v)",
			sendErr,
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

	if err := storageHelpers.SetState(static.StateEditingCheck, c, state); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't change a state (%v)",
			err,
		)
	}

	items, err := storageHelpers.GetItemsList(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't retrieve an items list (%v)",
			err,
		)
	}

	responseText, _ := responses.GetVerificationFinalStepResponse(check, items)

	if err := c.Send(responses.GetEditCheckMessage(responseText)); err != nil {
		return fmt.Errorf(
			"error in handlers.handleEditCheckName(): couldn't send a message (%v)",
			err,
		)
	}

	return nil
}

func handleEditCheckCreationDate(c tele.Context, state fsm.Context) error {
	gotDateStr := c.Message().Text
	newDate, errTime := time.Parse(time.DateTime, gotDateStr)
	if errTime != nil {
		errMsg := fmt.Sprintf(
			"error in handlers.handleEditCheckCreationDate(): \nfailed to parse a given date-time '%s' (%v)\n",
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

	items, errItems := storageHelpers.GetItemsList(c, state)

	var checkStr string

	if errItems == nil {
		checkStr, _ = responses.GetVerificationFinalStepResponse(check, items)
	}

	sendErr := c.Send(responses.GetEditCheckMessage(checkStr))
	stateErr := state.SetState(context.Background(), static.StateEditingCheck)

	if sendErr != nil || stateErr != nil {
		errMsg := "error in chandlers.handleEditCheckCreationDate(): "

		if sendErr != nil {
			errMsg += fmt.Sprintf(
				"\nfailed to send a message (%v)\n",
				sendErr,
			)
		}

		if stateErr != nil {
			errMsg += fmt.Sprintf(
				"\nfailed to set the state to 'StateEditingCheck' (%v)\n",
				stateErr,
			)
		}

		if err := c.Send("error: " + errMsg); err != nil {
			errMsg += fmt.Sprintf(
				"sent with error (%v)\n",
				err,
			)
		}

		return errors.New(errMsg)
	}

	return nil
}
