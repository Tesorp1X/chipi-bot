package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Service provides access to DB via methods. Only create new instances via MakeNewDBService.
type DBService struct {
	db *sql.DB
}

// Creates a new instance of DBService with provided path [dsnURI].
// Returns a pointer and a nil, if all goes well.
// If any error occurs, a nil is returned with a wrapped error.
func MakeNewDBService(dsnURI string) (*DBService, error) {
	db, err := sql.Open("sqlite3", dsnURI)
	if err != nil {
		return nil, fmt.Errorf(
			"error in db.MakeNewDBService(): couldn't open a db with uri: '%s' (%v)",
			dsnURI,
			err,
		)
	}

	return &DBService{db: db}, nil
}

// Creates a table with provided name and fields.
// If error occurs, returns a wrapped error.
func (dbs *DBService) createTable(name string, fields ...*field) error {
	sqlStmt, err := makeSqlStmtCreateTable(name, fields...)

	if err != nil {
		return fmt.Errorf(
			"error in db.createTable(): couldn't prepare a sql-statement (%v)",
			err,
		)
	}

	createTableStatement, err := dbs.db.Prepare(sqlStmt)
	if err != nil {

		return fmt.Errorf(
			"error in db.createTable(): couldn't prepare a db query (%v)",
			err,
		)
	}

	_, err = createTableStatement.Exec()
	if err != nil {
		return fmt.Errorf(
			"error in db.createTable(): couldn't execute a statement (%v)",
			err,
		)
	}

	return nil
}
