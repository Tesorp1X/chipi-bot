package static

//storage keys
const (
	ITEMS_LIST          = "items_list"
	CURRENT_INDEX_ITEMS = "current_index_items"

	CHECK_ID             = "check_id"
	CHECK                = "check"
	CHECKS               = "checks"
	CURRENT_INDEX_CHECKS = "current_index_checks"

	SESSION_ID           = "session_id"
	CURRENT_INDEX_TOTALS = "current_index_totals"
	SESSION_TOTALS       = "session_totals"

	IS_FROM_FINAL_STAGE = "is_from_final_stage"
)

//Callback data types
const (
	CallbackSelectorKeep   = "keep"
	CallbackSelectorChange = "change"

	CallbackEditItemName     = "item_name"
	CallbackEditItemPrice    = "item_price"
	CallbackEditItemAmount   = "item_amount"
	CallbackEditItemSubtotal = "item_subtotal"

	CallbackEditCheckName         = "check_name"
	CallbackEditCheckCreationDate = "check_date"
	CallbackEditCheckOwner        = "check_owner"
	CallbackEditCheckItems        = "check_items"

	CallbackMenuGoForward  = "go_forward"
	CallbackMenuGoBackward = "go_backward"
	CallbackMenuGoBack     = "go_back"

	CallbackOwnerPau  = "owner_pau"
	CallbackOwnerLiz  = "owner_liz"
	CallbackOwnerBoth = "owner_both"
)
