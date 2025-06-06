package models

import (
	"strings"

	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

const (
	StateDefault fsm.State = "default"

	StateStart fsm.State = "start"

	StateWaitForCheckName  fsm.State = "wait_for_check_name"
	StateWaitForCheckOwner fsm.State = "wait_for_check_owner"

	StateWaitForItemName  fsm.State = "wait_for_item_name"
	StateWaitForItemPrice fsm.State = "wait_for_item_price"
	StateWaitForItemOwner fsm.State = "wait_for_item_owner"
	StateWaitForNewItem   fsm.State = "wait_for_new_item"

	StateShowingChecks fsm.State = "showing_checks"
	StateShowingTotals fsm.State = "showing_totals"
)

type CallbackAction string

func (a CallbackAction) String() string {
	return string(a)
}

func (a CallbackAction) DataMatches(data string) bool {
	cringePrefix := "\f" + a.String()
	return data == cringePrefix || strings.HasPrefix(data, cringePrefix+"|")
}

const (
	CallbackActionCheckOwner CallbackAction = "check_owner"

	CallbackActionItemOwner  CallbackAction = "item_owner"
	CallbackActionHasNewItem CallbackAction = "has_new_item"

	CallbackActionMenuButtonPress CallbackAction = "menu_button_press"
)
