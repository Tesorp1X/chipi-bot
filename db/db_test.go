package db

import (
	"testing"
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
