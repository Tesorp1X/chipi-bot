package static

import (
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

// FSM-states
const (
	StateDefault fsm.State = "default"

	StateStart fsm.State = "start"

	StateWaitForCheckName fsm.State = "wait_for_check_name"

	StateWaitForNewCheckNameUnsaved fsm.State = "wait_for_new_check_name_unsaved"
	StateWaitForNewCheckName        fsm.State = "wait_for_new_check_name"

	StateWaitForCheckOwnerUnsaved fsm.State = "wait_for_check_owner_unsaved"
	StateWaitForCheckOwner        fsm.State = "wait_for_check_owner"

	StateWaitForCheckCreationDateUnsaved fsm.State = "wait_for_creation_date_unsaved"
	StateWaitForCheckCreationDate        fsm.State = "wait_for_creation_date"

	StateWaitingForCheckConfirmationUnsaved fsm.State = "waiting_for_check_confirmation_unsaved"
	StateWaitingForCheckConfirmation        fsm.State = "waiting_for_check_confirmation"

	StateWaitForItemName  fsm.State = "wait_for_item_name"
	StateWaitForItemPrice fsm.State = "wait_for_item_price"
	StateWaitForItemOwner fsm.State = "wait_for_item_owner"
	StateWaitForNewItem   fsm.State = "wait_for_new_item"

	StateShowingChecks fsm.State = "showing_checks"
	StateShowingTotals fsm.State = "showing_totals"

	StateShowingAnItem        fsm.State = "showing_item"
	StateShowingAnItemUnsaved fsm.State = "showing_item_unsaved"

	StateEditingCheck        fsm.State = "editing_check"
	StateEditingCheckUnsaved fsm.State = "editing_check_unsaved"

	StateEditingAnItem        fsm.State = "editing_item"
	StateEditingAnItemUnsaved fsm.State = "editing_item_unsaved"
)

const (
	// For usage in state.Finish method
	DELETE_DATA = true
	KEEP_DATA   = false
)
