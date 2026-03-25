package storageHelpers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Retrieves a new check name from a message text,
// verifies it and saves a new name to a check object.
// If an error occurs, a nil is returned alongside an error itself.
func SetNewCheckNameFromMessage(c tele.Context, state fsm.Context) (*static.Check, error) {
	//get checkData obj from context
	check, err := GetCheck(c, state)
	if err != nil {

		return nil, fmt.Errorf(
			"error in storageHelpers.SetNewCheckName(): couldn't retrieve check from context (%v).",
			err,
		)
	}
	// // todo: verify the message
	// if c.Message().Text == "" {
	// 	// ask again for the name
	// }
	//save new check name in it and put new version in context
	check.Name = c.Message().Text

	if err := state.Update(context.Background(), static.CHECK, check); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return nil, fmt.Errorf(
			"error in storageHelpers.SetNewCheckName(): couldn't save check in state-storage (%v). sent with error: %v",
			err,
			sendErr,
		)
	}
	return check, nil
}

// Retrieves new check owner from a callback-data and updates a Check.Owner object
// stored in a context-storage. Returns a pointer to a Check object with a nil,
// if no errors occurred. If an error occurs, while working with a context-storage,
// a nil is returned, alongside an error.
func SetNewCheckOwnerFromCallback(c tele.Context, state fsm.Context) (*static.Check, error) {
	owner := static.CallbackActionEditCheck.GetData(c.Callback().Data)

	switch owner {
	case static.CallbackOwnerLiz, static.CallbackOwnerPau, static.CallbackOwnerBoth:
		c.Respond(&tele.CallbackResponse{})
		check, err := GetCheck(c, state)
		if err != nil {
			return nil, fmt.Errorf(
				"error in storageHelpers.SetNewCheckOwnerFromCallback(): couldn't retrieve a check (%v)",
				err,
			)
		}

		check.Owner = owner

		if err := state.Update(context.Background(), static.CHECK, check); err != nil {
			sendErr := c.Send("error: couldn't save data in context")
			return nil, fmt.Errorf(
				"error in storageHelpers.SetNewCheckOwnerFromCallback(): couldn't save check in state-storage (%v). sent with error: %v",
				err,
				sendErr,
			)
		}

		return check, nil

	default:
		sendErr := c.Respond(&tele.CallbackResponse{Text: "error: invalid response: " + owner})
		return nil, fmt.Errorf(
			"error in storageHelpers.SetNewCheckOwnerFromCallback(): invalid response (%s). sent with error: %v",
			owner,
			sendErr,
		)
	}
}

func SetState(newState fsm.State, c tele.Context, state fsm.Context) error {
	if err := state.SetState(context.Background(), newState); err != nil {
		sendErr := c.Send("error: couldn't change a state")
		return fmt.Errorf(
			"error in storageHelpers.SetState(): couldn't change a state to %s (%v). sent with error: %v",
			newState,
			err,
			sendErr,
		)
	}

	return nil
}

// Sets a state of fsm.Context to default and removes data, if removeData is true.
// Returns any error, that happened during work.
func FinishState(removeData bool, c tele.Context, state fsm.Context) error {
	currentState, _ := state.State(context.Background())
	if stateErr := state.Finish(context.Background(), removeData); stateErr != nil {
		sendErr := c.Send("error: failed to finish your state and delete all context data")
		return fmt.Errorf(
			"error in storageHelpers.FinishState(): failed to finish '%s' state (%v). sent with err (%v)",
			currentState,
			stateErr,
			sendErr,
		)
	}

	msg := "error happened: all data erased"
	if sendErr := c.Send(msg); sendErr != nil {
		return fmt.Errorf(
			"error in storageHelpers.FinishState(): failed to to send a message '%s' (%v).",
			msg,
			sendErr,
		)
	}

	return nil
}
