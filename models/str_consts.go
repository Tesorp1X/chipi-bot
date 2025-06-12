package models

const (
	CHECK_OWNER = "check_owner"
	CHECK_NAME  = "check_name"

	OWNER_LIZ  = "liz"
	OWNER_PAU  = "pau"
	OWNER_BOTH = "both"

	ITEM_NAME  = "item_name"
	ITEM_PRICE = "item_price"
	ITEM_OWNER = "item_owner"

	HAS_MORE_ITEMS_TRUE  = "true"
	HAS_MORE_ITEMS_FALSE = "false"

	BTN_FORWARD = "forward"
	BTN_BACK    = "back"
	BTN_EDIT    = "edit"
)

// storage keys
const (
	ITEMS_LIST     = "items_list"
	CHECK_ID       = "check_id"
	SESSION_ID     = "session_id"
	CHECK          = "check"
	CHECKS         = "checks"
	CURRENT_INDEX  = "current_index"
	SESSION_TOTALS = "session_totals"
)

// error msgs
const (
	ErrorSometingWentWrong         = "Произошла ошибка. Am souwy😔"
	ErrorStateDataUpdate           = "Произошла ошибка при сохранении данных состояния."
	ErrorSetState                  = "Произошла ошибка при смене состояния."
	ErrorNameMustBeTxtMsg          = "Название обязательно должно быть текстовым сообщением."
	ErrorItemPriceMustBeANumberMsg = "Цена должна быть числом без пробелов."
	ErrorItemsListNotFound         = "Ошибка, лист с покупками не найден."
	ErrorSavingInDB                = "Произошла ошибка при сохранении данных в базу данных."
	ErrorInvalidRequest            = "Invalid request"
	ErrorTryAgain                  = "Произошла ошибка, попробуйте еще раз."
)
