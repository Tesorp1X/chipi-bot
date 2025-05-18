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
)

// storage keys
const (
	ITEMS_LIST = "items_list"
)

// error msgs
const (
	ErrorSometingWentWrong     = "Произошла ошибка. Am souwy😔"
	ErrorStateDataUpdate       = "Произошла ошибка при сохранении данных состояния."
	ErrorSetState              = "Произошла ошибка при смене состояния."
	ErrorNameMustBeTxtMsg      = "Название обязательно должно быть текстовым сообщением."
	ErrorItemPriceMustBeIntMsg = "Цена должна быть целым числом без пробелов и других знаков (только цифры)."
	ErrorItemsListNotFound     = "Ошибка, лист с покупками не найден."
)
