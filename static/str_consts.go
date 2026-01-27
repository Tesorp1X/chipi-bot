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
)

//Callback data types
const (
	CallbackSelectorKeep   = "keep"
	CallbackSelectorChange = "change"

	CallbackEditItemName     = "name"
	CallbackEditItemPrice    = "price"
	CallbackEditItemAmount   = "amount"
	CallbackEditItemSubtotal = "subtotal"
	CallbackEditItemGoBack   = "go_back"

	CallbackOwnerPau  = "owner_pau"
	CallbackOwnerLiz  = "owner_liz"
	CallbackOwnerBoth = "owner_both"
)
