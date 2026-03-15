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

var (
	checksTableFields []*field = []*field{
		{
			Name:         "id",
			Type:         IntegerType,
			IsPrimaryKey: true,
		},
		{
			Name: "session_id",
			Type: IntegerType,

			IsForeignKey: true,
			RefTableName: SESSIONS_TABLE_NAME,
			RefFieldName: "id",
		},
		{
			Name: "name",
			Type: TextType,
		},
		{
			Name: "orgname",
			Type: TextType,
		},
		{
			Name: "owner",
			Type: TextType,
		},
		{
			Name: "total",
			Type: RealType,
		},
		{
			Name: "total_pau",
			Type: RealType,
		},
		{
			Name: "total_liz",
			Type: RealType,
		},
		{
			Name: "date_of_purchase",
			Type: TextType,
		},
	}

	itemsTableFields []*field = []*field{
		{
			Name:         "id",
			Type:         IntegerType,
			IsPrimaryKey: true,
		},
		{
			Name: "check_id",
			Type: IntegerType,

			IsForeignKey: true,
			RefTableName: CHECKS_TABLE_NAME,
			RefFieldName: "id",
		},
		{
			Name: "name",
			Type: TextType,
		},
		{
			Name: "owner",
			Type: TextType,
		},
		{
			Name: "price",
			Type: RealType,
		},
		{
			Name: "amount",
			Type: RealType,
		},
		{
			Name: "subtotal",
			Type: RealType,
		},
	}

	sessionsTableFields []*field = []*field{
		{
			Name:         "id",
			Type:         IntegerType,
			IsPrimaryKey: true,
		},
		{
			Name: "opened_at",
			Type: TextType,
		},
		{
			Name: "closed_at",
			Type: TextType,
		},
		{
			Name: "total_spent",
			Type: RealType,
		},
		{
			Name: "total_pau_spent",
			Type: RealType,
		},
		{
			Name: "total_liz_spent",
			Type: RealType,
		},
		{
			Name: "total_topups",
			Type: RealType,
		},
		{
			Name: "total_pau_topups",
			Type: RealType,
		},
		{
			Name: "total_liz_topups",
			Type: RealType,
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
