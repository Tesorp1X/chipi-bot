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
	BTN_CHECKS  = "checks"
)

// storage keys
const (
	ITEMS_LIST = "items_list"

	CHECK_ID             = "check_id"
	CHECK                = "check"
	CHECKS               = "checks"
	CURRENT_INDEX_CHECKS = "current_index_checks"

	SESSION_ID           = "session_id"
	CURRENT_INDEX_TOTALS = "current_index_totals"
	SESSION_TOTALS       = "session_totals"
)

// error msgs
const (
	ErrorSomethingWentWrong = "Произошла ошибка. Am souwy😔"
	ErrorStateDataUpdate    = "Произошла ошибка при сохранении данных состояния."
	ErrorSetState           = "Произошла ошибка при смене состояния."
	ErrorNameMustBeTxtMsg   = "Название обязательно должно быть текстовым сообщением."

	ErrorItemPriceMustBeANumberMsg     = "Цена должна быть целым числом без пробелов."
	ErrorAmountOfItemsMustBeANumberMsg = "Количество должно быть целым числом без пробелов."

	ErrorItemsListNotFound = "Ошибка, лист с покупками не найден."
	ErrorSavingInDB        = "Произошла ошибка при сохранении данных в базу данных."
	ErrorInvalidRequest    = "Invalid request"
	ErrorTryAgain          = "Произошла ошибка, попробуйте еще раз."

	ErrorNotImplemented = "Этот раздел пока в разработке 🛠️"
)

// help msgs
const (
	AmountOfItemsHelpMsg = "Если хотите указать несколько товаров, то делайте это в формате: <i>&lt;кол-во&gt;*&lt;цена&gt;</i>"
)
