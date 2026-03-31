package storageHelpers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	tele "gopkg.in/telebot.v4"
)

// Updates a CURRENT_INDEX_ITEMS value in a context-storage.
// If error occurs, it will bbe wrapped and returned.
func UpdateCurrentItemsIndex(newCurrentIndex int, c tele.Context, state fsm.Context) error {
	if err := state.Update(context.Background(), static.CURRENT_INDEX_ITEMS, newCurrentIndex); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in storageHelpers.UpdateCurrentIndex(): couldn't save current index in state-storage (%v). sent with error: %v",
			err,
			sendErr,
		)
	}

	return nil
}

// Updates a CHECK value in a context-storage.
// If error occurs, it will bbe wrapped and returned.
func UpdateCheck(check *static.Check, c tele.Context, state fsm.Context) error {
	if err := state.Update(context.Background(), static.CHECK, check); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in storageHelpers.UpdateCheck(): couldn't save check in state-storage (%v). sent with error: %v",
			err,
			sendErr,
		)
	}

	return nil
}

// Updates a ITEMS_LIST value in a context-storage.
// If error occurs, it will bbe wrapped and returned.
func UpdateItemsList(items []*static.Item, c tele.Context, state fsm.Context) error {
	if err := state.Update(context.Background(), static.ITEMS_LIST, items); err != nil {
		sendErr := c.Send("error: couldn't save data in context")
		return fmt.Errorf(
			"error in storageHelpers.UpdateItemsList(): couldn't save items in state-storage (%v). sent with error: %v",
			err,
			sendErr,
		)
	}

	return nil
}

const (
	incrementIndex = 1
	decrementIndex = -1
)

func changeCurrentItemsListIndex(diff int, c tele.Context, state fsm.Context) error {
	if diff != incrementIndex && diff != decrementIndex {
		return fmt.Errorf(
			"error in storageHelpers.changeCurrentItemsListIndex(): invalid change value (%d)",
			diff,
		)
	}

	currentIndex, err := GetCurrentIndex(c, state)
	if err != nil {
		return fmt.Errorf(
			"error in storageHelpers.changeCurrentItemsListIndex(): failed to retrieve a current items index (%v)",
			err,
		)
	}

	if err := UpdateCurrentItemsIndex(currentIndex+diff, c, state); err != nil {
		return fmt.Errorf(
			"error in storageHelpers.changeCurrentItemsListIndex(): failed to update a current items index (%v)",
			err,
		)
	}

	return nil
}

func IncrementCurrentItemsListIndex(c tele.Context, state fsm.Context) error {
	if err := changeCurrentItemsListIndex(incrementIndex, c, state); err != nil {
		sendErr := c.Send("error: couldn't update data in context")
		return fmt.Errorf(
			"error in storageHelpers.changeCurrentItemsListIndex(): failed to increment a current items index (%v), sent with an err (%v)",
			err,
			sendErr,
		)
	}

	return nil
}

func DecrementCurrentItemsListIndex(c tele.Context, state fsm.Context) error {
	if err := changeCurrentItemsListIndex(decrementIndex, c, state); err != nil {
		sendErr := c.Send("error: couldn't update data in context")
		return fmt.Errorf(
			"error in storageHelpers.changeCurrentItemsListIndex(): failed to decrement a current items index (%v), sent with an err (%v)",
			err,
			sendErr,
		)
	}

	return nil
}
