package db

// SQLite data-types
const (
	IntegerType = "INTEGER"
	TextType    = "TEXT"
	RealType    = "REAL"
)

//Tables names
const (
	CHECKS_TABLE_NAME   = "checks"
	ITEMS_TABLE_NAME    = "items"
	SESSIONS_TABLE_NAME = "sessions"
)

// Column names
const (
	// Checks

	CHECKS_ID               = "id"
	CHECKS_SESSION_ID       = "session_id"
	CHECKS_NAME             = "name"
	CHECKS_ORGNAME          = "orgname"
	CHECKS_OWNER            = "owner"
	CHECKS_TOTAL            = "total"
	CHECKS_TOTAL_PAU        = "total_pau"
	CHECKS_TOTAL_LIZ        = "total_liz"
	CHECKS_DATE_OF_PURCHASE = "date_of_purchase"

	// Items

	ITEMS_ID       = "id"
	ITEMS_CHECK_ID = "check_id"
	ITEMS_NAME     = "name"
	ITEMS_OWNER    = "owner"
	ITEMS_PRICE    = "price"
	ITEMS_AMOUNT   = "amount"
	ITEMS_SUBTOTAL = "subtotal"

	// Sessions

	//SESSIONS_
	SESSIONS_ID              = "id"
	SESSIONS_OPENED_AT       = "opened_at"
	SESSIONS_CLOSED_AT       = "closed_at"
	SESSIONS_IS_OPEN         = "is_open"
	SESSIONS_TOTAL_SPENT     = "total_spent"
	SESSIONS_PAU_TOTAL_SPENT = "total_pau_spent"
	SESSIONS_LIZ_TOTAL_SPENT = "total_liz_spent"
	SESSIONS_TOPUPS          = "total_topups"
	SESSIONS_PAU_TOPUPS      = "total_pau_topups"
	SESSIONS_LIZ_TOPUPS      = "total_liz_topups"
)

var (
	checksTableFields []*field = []*field{
		{
			Name:         CHECKS_ID,
			Type:         IntegerType,
			IsPrimaryKey: true,
		},
		{
			Name: CHECKS_SESSION_ID,
			Type: IntegerType,

			IsForeignKey: true,
			RefTableName: SESSIONS_TABLE_NAME,
			RefFieldName: SESSIONS_ID,
		},
		{
			Name: CHECKS_NAME,
			Type: TextType,
		},
		{
			Name: CHECKS_ORGNAME,
			Type: TextType,
		},
		{
			Name: CHECKS_OWNER,
			Type: TextType,
		},
		{
			Name: CHECKS_TOTAL,
			Type: RealType,
		},
		{
			Name: CHECKS_TOTAL_PAU,
			Type: RealType,
		},
		{
			Name: CHECKS_TOTAL_LIZ,
			Type: RealType,
		},
		{
			Name: CHECKS_DATE_OF_PURCHASE,
			Type: TextType,
		},
	}

	itemsTableFields []*field = []*field{
		{
			Name:         ITEMS_ID,
			Type:         IntegerType,
			IsPrimaryKey: true,
		},
		{
			Name: ITEMS_CHECK_ID,
			Type: IntegerType,

			IsForeignKey: true,
			RefTableName: CHECKS_TABLE_NAME,
			RefFieldName: CHECKS_ID,
		},
		{
			Name: ITEMS_NAME,
			Type: TextType,
		},
		{
			Name: ITEMS_OWNER,
			Type: TextType,
		},
		{
			Name: ITEMS_PRICE,
			Type: RealType,
		},
		{
			Name: ITEMS_AMOUNT,
			Type: RealType,
		},
		{
			Name: ITEMS_SUBTOTAL,
			Type: RealType,
		},
	}

	sessionsTableFields []*field = []*field{
		{
			Name:         SESSIONS_ID,
			Type:         IntegerType,
			IsPrimaryKey: true,
		},
		{
			Name: SESSIONS_OPENED_AT,
			Type: TextType,
		},
		{
			Name: SESSIONS_CLOSED_AT,
			Type: TextType,

			IsNullable: true,
		},
		{
			Name: SESSIONS_IS_OPEN,
			Type: TextType,
		},
		{
			Name: SESSIONS_TOTAL_SPENT,
			Type: RealType,

			IsNullable: true,
		},
		{
			Name: SESSIONS_PAU_TOTAL_SPENT,
			Type: RealType,

			IsNullable: true,
		},
		{
			Name: SESSIONS_LIZ_TOTAL_SPENT,
			Type: RealType,

			IsNullable: true,
		},
		{
			Name: SESSIONS_TOPUPS,
			Type: RealType,

			IsNullable: true,
		},
		{
			Name: SESSIONS_PAU_TOPUPS,
			Type: RealType,

			IsNullable: true,
		},
		{
			Name: SESSIONS_LIZ_TOPUPS,
			Type: RealType,

			IsNullable: true,
		},
	}
)

type tableNameAndFields struct {
	Name   string
	Fields []*field
}

type tables []*tableNameAndFields

// Array of all tables
var tablesWithNames = tables{
	{
		Name:   SESSIONS_TABLE_NAME,
		Fields: sessionsTableFields,
	},
	{
		Name:   CHECKS_TABLE_NAME,
		Fields: checksTableFields,
	},
	{
		Name:   ITEMS_TABLE_NAME,
		Fields: itemsTableFields,
	},
}
