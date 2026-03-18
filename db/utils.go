package db

import (
	"fmt"
	"strings"
)

type column struct {
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
func makeSqlStmtCreateTable(tableName string, columns ...*column) (string, error) {
	var foreignKeys strings.Builder
	var sqlStmt strings.Builder
	sqlStmt.WriteString("CREATE TABLE IF NOT EXISTS " + tableName + " (")
	for i, col := range columns {
		if i == 0 {
			if !col.IsPrimaryKey {
				return "", fmt.Errorf(
					"error in db.makeSqlStmtCreateTable(): first field must be a primary key",
				)
			}
			sqlStmt.WriteString(col.Name + " " + col.Type + " PRIMARY KEY" + ", ")
		} else {
			sqlStmt.WriteString(col.Name + " " + col.Type)

			if !col.IsNullable {
				sqlStmt.WriteString(" NOT NULL")
			}

			if i != len(columns)-1 {
				sqlStmt.WriteString(", ")
			}

			if col.IsForeignKey {
				switch {
				case col.RefTableName == "":
					return "", fmt.Errorf(
						"error in db.makeSqlStmtCreateTable(): RefTableName must be not empty if IsForeignKey is true",
					)
				case col.RefFieldName == "":
					return "", fmt.Errorf(
						"error in db.makeSqlStmtCreateTable(): RefFieldName must be not empty if IsForeignKey is true",
					)
				default:
					foreignKeys.WriteString(", FOREIGN KEY(" + col.Name + ") REFERENCES " + col.RefTableName + " (" + col.RefFieldName + ") ")
				}

			}
		}
	}
	sqlStmt.WriteString(foreignKeys.String() + ")")

	return sqlStmt.String(), nil
}
