package db

import (
	"fmt"
)

type field struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsNullable   bool
	IsForeignKey bool
	RefTableName string
	RefFieldName string
}

// Returns a ready to use SQL-statement for creating a table with provided name and fields.
// If error occurs, returns an empty string with an error.
func makeSqlStmtCreateTable(tableName string, fields ...*field) (string, error) {
	var foreignKeys string
	sqlStmt := "CREATE TABLE IF NOT EXISTS " + tableName + " ("
	for i, f := range fields {
		if i == 0 {
			if !f.IsPrimaryKey {
				return "", fmt.Errorf(
					"error in db.makeSqlStmtCreateTable(): first field must be a primary key",
				)
			}
			sqlStmt += f.Name + " " + f.Type + " PRIMARY KEY" + ", "
		} else {
			sqlStmt += f.Name + " " + f.Type

			if !f.IsNullable {
				sqlStmt += " NOT NULL"
			}

			if i != len(fields)-1 {
				sqlStmt += ", "
			}

			if f.IsForeignKey {
				switch {
				case f.RefTableName == "":
					return "", fmt.Errorf(
						"error in db.makeSqlStmtCreateTable(): RefTableName must be not empty if IsForeignKey is true",
					)
				case f.RefFieldName == "":
					return "", fmt.Errorf(
						"error in db.makeSqlStmtCreateTable(): RefFieldName must be not empty if IsForeignKey is true",
					)
				default:
					foreignKeys += ", FOREIGN KEY(" + f.Name + ") REFERENCES " + f.RefTableName + " (" + f.RefFieldName + ") "
				}

			}
		}
	}
	sqlStmt += foreignKeys + ")"

	return sqlStmt, nil
}
