package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tesorp1X/chipi-bot/config"
	"github.com/Tesorp1X/chipi-bot/static"
	"github.com/doug-martin/goqu/v9"

	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

// Service provides access to DB via methods. Only create new instances via MakeNewDBService.
type DBService struct {
	db *sql.DB
}

// Creates a new instance of DBService with provided path [dsnURI].
// Returns a pointer and a nil, if all goes well.
// If any error occurs, a nil is returned with a wrapped error.
func MakeNewDBService(conf *config.Config) (*DBService, error) {
	db, err := sql.Open("sqlite3", conf.DbPath)
	if err != nil {
		return nil, fmt.Errorf(
			"error in db.MakeNewDBService(): couldn't open a db with uri: '%s' (%v)",
			conf.DbPath,
			err,
		)
	}

	dbs := &DBService{db: db}

	if err := dbs.CreateIfNotExists(); err != nil {
		return nil, fmt.Errorf(
			"error in db.MakeNewDBService(): failed to create tables (%v)",
			err,
		)
	}

	return dbs, nil
}

func (dbs *DBService) Close() error {
	err := dbs.db.Close()
	if err != nil {
		return fmt.Errorf(
			"error in db.Close(): failed to close a db (%v)",
			err,
		)
	}

	return nil
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
		Cols(SESSIONS_OPENED_AT, SESSIONS_IS_OPEN).
		Vals(
			goqu.Vals{time.Now().Format(time.DateTime), "true"},
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
	selectRowSql := fmt.Sprintf(
		"SELECT id FROM %s WHERE %s = ?",
		SESSIONS_TABLE_NAME,
		SESSIONS_IS_OPEN,
	)
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

// Creates a new check record in given transaction. Returns an id of that record,
// if no errors occurred. Otherwise, returns -1 and an error.
func (dbs *DBService) addCheck(tx *sql.Tx, check *static.Check) (int64, error) {
	// save the check
	sessionId, err := dbs.getOrCreateSession(tx)
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.addCheck(): failed to retrieve a session_id (%v)",
			err,
		)
	}

	ds := goqu.Insert(CHECKS_TABLE_NAME).
		Cols(CHECKS_SESSION_ID, CHECKS_NAME, CHECKS_ORGNAME,
			CHECKS_OWNER, CHECKS_TOTAL, CHECKS_TOTAL_PAU,
			CHECKS_TOTAL_LIZ, CHECKS_DATE_OF_PURCHASE).
		Vals(
			goqu.Vals{
				sessionId, check.Name, check.OrgName, check.Owner,
				check.Total, check.TotalPau, check.TotalLiz,
				check.Date.Format(time.DateTime)},
		)
	insertSql, args, _ := ds.ToSQL()

	statement, err := tx.Prepare(insertSql)
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.addCheck(): failed to prepare a statement '%s' (%v)",
			insertSql,
			err,
		)
	}

	res, err := statement.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.addCheck(): failed to execute a statement '%s' (%v)",
			insertSql,
			err,
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf(
			"error in db.addCheck(): to retrieve an id from result (%v)",
			err,
		)
	}

	return id, nil

}

// This method creates a new record in checks and records for every item,
// associated with that check. The method will assign a current session to a check.
func (dbs *DBService) AddNewCheckWithItems(check *static.Check, items []*static.Item) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := dbs.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"error in db.AddNewCheckWithItems(): failed to begin a transaction (%v)",
			err,
		)
	}

	var committed bool
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	checkId, err := dbs.addCheck(tx, check)
	if err != nil {
		return fmt.Errorf(
			"error in db.AddNewCheckWithItems(): failed to add check (%v)",
			err,
		)
	}

	insertItemSql := fmt.Sprintf(
		//"INSERT INTO items (check_id, name, owner, price, amount, subtotal) VALUES (?, ?, ?, ?, ?, ?)"
		"INSERT INTO items (%s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?)",
		ITEMS_CHECK_ID,
		ITEMS_NAME,
		ITEMS_OWNER,
		ITEMS_PRICE,
		ITEMS_AMOUNT,
		ITEMS_SUBTOTAL,
	)
	addItemStatement, err := tx.Prepare(insertItemSql)
	if err != nil {
		return fmt.Errorf(
			"error in db.AddNewCheckWithItems(): failed to prepare a statement with sql: '%s' (%v)",
			insertItemSql,
			err,
		)
	}

	for _, item := range items {
		if _, err := addItemStatement.Exec(
			checkId, item.Name, item.Name,
			item.Price, item.Amount, item.Subtotal); err != nil {
			return fmt.Errorf(
				"error in db.AddNewCheckWithItems(): couldn't insert an item %+v (%v)",
				item,
				err,
			)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	committed = true

	return nil
}
