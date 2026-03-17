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
// Creates all the tables from tablesWithNames array. Returns an error if anything goes wrong.
func (dbs *DBService) CreateIfNotExists() error {
	for _, table := range tablesWithNames {
		if err := dbs.createTable(table.Name, table.Fields...); err != nil {
			return fmt.Errorf(
				"in db.CreateIfNotExists(): couldn't create a '%s' table (%v)",
				table.Name,
				err,
			)
		}
	}

	return nil
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

	if _, err = createTableStatement.Exec(); err != nil {
		return fmt.Errorf(
			"error in db.createTable(): couldn't execute a statement (%v)",
			err,
		)
	}

	if err := createTableStatement.Close(); err != nil {
		return fmt.Errorf(
			"error in db.createTable(): failed to close a createTableStatement (%v)",
			err,
		)
	}

	return nil
}

// Creates a new session in given transaction. Returns an id of that record,
// if no errors occurred. Otherwise, returns -1 and an error.
func (dbs *DBService) createNewSession(tx *sql.Tx) (int64, error) {
	ds := goqu.Insert(SESSIONS_TABLE_NAME).
		Cols("created_at", "is_open").
		Vals(
			goqu.Vals{"time", "true"},
		)
	insertSql, args, _ := ds.ToSQL()

	statement, err := tx.Prepare(insertSql)
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.createNewSession(): failed to prepare a statement '%s' (%v)",
			insertSql,
			err,
		)
	}

	res, err := statement.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.createNewSession(): failed to execute a statement '%s' (%v)",
			insertSql,
			err,
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.createNewSession(): to retrieve an id from result (%v)",
			err,
		)
	}

	return id, nil
}

// Returns an id of a current session in given transaction.
// If there is no active session, then it will be created first.
// If anything goes wrong, '-1' and dn error is returned.
func (dbs *DBService) getOrCreateSession(tx *sql.Tx) (int64, error) {
	selectRowSql := fmt.Sprintf("SELECT id FROM %s WHERE is_open = ?", SESSIONS_TABLE_NAME)
	row := tx.QueryRow(selectRowSql, "true")
	var id int64 = -1

	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		if id, err = dbs.createNewSession(tx); err != nil {
			return -1, fmt.Errorf(
				"error in db.getOrCreateSession(): failed to create a new session (%v)",
				err,
			)
		}

	}

	if err != nil {
		return -1, fmt.Errorf(
			"error in db.getOrCreateSession(): failed to scan a row (%v)",
			err,
		)
	}

	return id, nil
}

	return nil
}
