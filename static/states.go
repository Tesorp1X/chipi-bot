package static

import (
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

// FSM-states
const (
	StateDefault fsm.State = "default"

	StateStart fsm.State = "start"

	StateWaitForCheckName    fsm.State = "wait_for_check_name"
	StateWaitForNewCheckName fsm.State = "wait_for_new_check_name"
	StateWaitForCheckOwner   fsm.State = "wait_for_check_owner"

	StateWaitForItemName  fsm.State = "wait_for_item_name"
	StateWaitForItemPrice fsm.State = "wait_for_item_price"
	StateWaitForItemOwner fsm.State = "wait_for_item_owner"
	StateWaitForNewItem   fsm.State = "wait_for_new_item"

	StateShowingChecks fsm.State = "showing_checks"
	StateShowingTotals fsm.State = "showing_totals"
	StateShowingAnItem fsm.State = "showing_item"

	StateEditingCheck fsm.State = "editing_check"
)
