package storageHelpers

import (
	"context"
	"fmt"

	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"

	tele "gopkg.in/telebot.v4"
)

// Retrieves ITEMS_LIST from given context.
// If error occurs, an empty slice is returned alongside an error itself.
func GetItemsList(c tele.Context, state fsm.Context) ([]*static.Item, error) {
	var items []*static.Item
	if err := state.Data(context.Background(), static.ITEMS_LIST, &items); err != nil {
		sendErr := c.Send("error: couldn't retrieve data from context")
		return nil, fmt.Errorf(
			"error in handleCheckName(): couldn't retrieve items from state-storage (%v). send with error: %v",
			err,
			sendErr,
		)
	}

	if len(items) == 0 {
		sendErr := c.Send("error: retrieved items-list iss empty")
		return nil, fmt.Errorf(
			"error in handleCheckName(): retrieved items-list is empty. send with error: %v",
			sendErr,
		)
	}

	return items, nil
}
