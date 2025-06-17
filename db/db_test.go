package db

import (
	"database/sql"
	"testing"

	"github.com/Tesorp1X/chipi-bot/models"
	_ "github.com/mattn/go-sqlite3"
)

func TestMakeSqlStmt(t *testing.T) {

	t.Run("no foreign keys no errors", func(t *testing.T) {
		feilds := []Field{
			{
				Name:         "id",
				Type:         "INTEGER",
				IsPrimeryKey: true,
			},
			{
				Name:       "Name",
				Type:       "STRING",
				IsNullable: false,
			},
			{
				Name:       "Owner",
				Type:       "STRING",
				IsNullable: false,
			},
		}

		stmtGot, err := makeSqlStmtCreateTable("checks", feilds...)
		stmtWant := "CREATE TABLE IF NOT EXISTS checks (id INTEGER PRIMARY KEY, Name STRING NOT NULL, Owner STRING NOT NULL)"

		if err != nil {
			t.Fatalf("didn't expect error, but goy: %v", err)
		}

		if stmtGot != stmtWant {
			t.Fatalf("wanted:\n%s\ngot:\n%s", stmtWant, stmtGot)
		}
	})
	t.Run("one foreign key no errors", func(t *testing.T) {
		feilds := []Field{
			{
				Name:         "id",
				Type:         "INTEGER",
				IsPrimeryKey: true,
			},
			{
				Name:         "check_id",
				Type:         "INTEGER",
				IsNullable:   false,
				IsForeignKey: true,
				RefTableName: "checks",
				RefFieldName: "id",
			},
			{
				Name:       "Name",
				Type:       "STRING",
				IsNullable: false,
			},
			{
				Name:       "Owner",
				Type:       "STRING",
				IsNullable: true,
			},
			{
				Name:       "Price",
				Type:       "INTEGER",
				IsNullable: false,
			},
		}

		stmtGot, err := makeSqlStmtCreateTable("items", feilds...)
		stmtWant := "CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY, check_id INTEGER NOT NULL, Name STRING NOT NULL, Owner STRING, Price INTEGER NOT NULL, FOREIGN KEY(check_id) REFERENCES checks (id) )"

		if err != nil {
			t.Fatalf("didn't expect error, but goy: %v", err)
		}

		if stmtGot != stmtWant {
			t.Fatalf("wanted:\n%s\ngot:\n%s", stmtWant, stmtGot)
		}
	})
}

func makeInMemoryDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	sessionFields := []Field{
		{
			Name:         "id",
			Type:         "INTEGER",
			IsPrimeryKey: true,
		},
		{
			Name: "opened_at",
			Type: "TEXT", // time as text formated as 2006-01-02 15:04:05

		},
		{
			Name:       "closed_at",
			Type:       "TEXT", // time as text formated as 2006-01-02 15:04:05
			IsNullable: true,
		},
		{
			Name: "is_open",
			Type: "TEXT", // bool value as string
		},
	}
	if err = createTable(db, "sessions", sessionFields...); err != nil {
		t.Fatalf("error couldn't create a db: %v", err)
	}

	checkFeilds := []Field{
		{
			Name:         "id",
			Type:         "INTEGER",
			IsPrimeryKey: true,
		},
		{
			Name:       "Name",
			Type:       "TEXT",
			IsNullable: false,
		},
		{
			Name:       "Owner",
			Type:       "TEXT",
			IsNullable: false,
		},
	}
	if err = createTable(db, "checks", checkFeilds...); err != nil {
		t.Fatalf("error couldn't create a db: %v", err)
	}

	itemsFeilds := []Field{
		{
			Name:         "id",
			Type:         "INTEGER",
			IsPrimeryKey: true,
		},
		{
			Name: "check_id",
			Type: "INTEGER",

			IsForeignKey: true,
			RefTableName: "checks",
			RefFieldName: "id",
		},
		{
			Name: "Name",
			Type: "TEXT",
		},
		{
			Name: "Owner",
			Type: "TEXT",
		},
		{
			Name: "Price",
			Type: "REAL",
		},
	}
	if err = createTable(db, "items", itemsFeilds...); err != nil {
		t.Fatalf("error couldn't create a db: %v", err)

	}

	totalsField := []Field{
		{
			Name:         "id",
			Type:         "INTEGER",
			IsPrimeryKey: true,
		},
		{
			Name: "session_id",
			Type: "INTEGER",

			IsForeignKey: true,
			RefTableName: "seesions",
			RefFieldName: "id",
		},
		{
			Name: "total",
			Type: "REAL",
		},
		{
			Name: "recipient",
			Type: "TEXT",
		},
		{
			Name: "amount",
			Type: "REAL",
		},
	}
	if err = createTable(db, "totals", totalsField...); err != nil {
		t.Fatalf("error couldn't create a db: %v", err)

	}

	checksAndSessionsFields := []Field{
		{
			Name:         "id",
			Type:         "INTEGER",
			IsPrimeryKey: true,
		},
		{
			Name: "session_id",
			Type: "INTEGER",

			IsForeignKey: true,
			RefTableName: "seesions",
			RefFieldName: "id",
		},
		{
			Name: "check_id",
			Type: "INTEGER",

			IsForeignKey: true,
			RefTableName: "checks",
			RefFieldName: "id",
		},
	}
	if err = createTable(db, "checks_and_sessions", checksAndSessionsFields...); err != nil {
		t.Fatalf("error couldn't create a db: %v", err)

	}

	return db
}

// Inserts items in db. fails the test if error occures
func populateItemsDB(t *testing.T, db *sql.DB, items []models.Item) {
	t.Helper()
	for _, item := range items {
		_, err := db.Exec("INSERT INTO items (id, check_id, name, owner, price) VALUES (?, ?, ?, ?, ?)",
			item.Id, item.CheckId, item.Name, item.Owner, item.Price)
		if err != nil {
			t.Fatalf("failed to insert item: %v", err)
		}
	}
}
