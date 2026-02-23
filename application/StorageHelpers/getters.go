package storageHelpers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	tele "gopkg.in/telebot.v4"
)

// Retrieves ITEMS_LIST from a given context.
// If error occurs, an empty slice is returned alongside an error itself.
func GetItemsList(c tele.Context, state fsm.Context) ([]*static.Item, error) {
	var items []*static.Item
	if err := state.Data(context.Background(), static.ITEMS_LIST, &items); err != nil {
		sendErr := c.Send("error: couldn't retrieve data from context")
		return nil, fmt.Errorf(
			"error in handleCheckName(): couldn't retrieve items from state-storage (%v). sent with error: %v",
			err,
			sendErr,
		)
	}

	if len(items) == 0 {
		sendErr := c.Send("error: retrieved items-list iss empty")
		return nil, fmt.Errorf(
			"error in handleCheckName(): retrieved items-list is empty. sent with error: %v",
			sendErr,
		)
	}

	return items, nil
}

// Retrieves a CHECK object from a given context.
// If error occurs, a nil is returned alongside an error itself.
func GetCheck(c tele.Context, state fsm.Context) (*static.Check, error) {
	var check *static.Check
	if err := state.Data(context.Background(), static.CHECK, &check); err != nil {
		sendErr := c.Send("error: couldn't retrieve check data from context")
		return nil, fmt.Errorf(
			"error in GetCheck(): couldn't retrieve check from state-storage (%v). sent with error: %v",
			err,
			sendErr,
		)
	}

	return check, nil
}

// Retrieves CURRENT_INDEX_ITEMS from a given context.
// If error occurs, a -1 is returned alongside an error itself.
func GetCurrentIndex(c tele.Context, state fsm.Context) (int, error) {
	var currentIndex int
	if err := state.Data(context.Background(), static.CURRENT_INDEX_ITEMS, &currentIndex); err != nil {
		sendErr := c.Send("error: couldn't retrieve current items's index")
		return -1, fmt.Errorf(
			"error in GetCurrentIndex(): couldn't retrieve currentIndex from context (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	return currentIndex, nil
}
