package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Tesorp1X/chipi-bot/models"
	_ "github.com/mattn/go-sqlite3"
)

type Field struct {
	Name         string
	Type         string
	IsPrimeryKey bool
	IsNullable   bool
	IsForeignKey bool
	RefTableName string
	RefFieldName string
}

func makeSqlStmtCreateTable(name string, fields ...Field) (string, error) {
	var foreignKeys string
	sqlStmt := "CREATE TABLE IF NOT EXISTS " + name + " ("
	for i, f := range fields {
		if i == 0 {
			if !f.IsPrimeryKey {
				return "", errors.New("first field must be a primery key")
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
					return "", errors.New("RefTableName must be not empty if IsForeignKey is true")
				case f.RefFieldName == "":
					return "", errors.New("RefFieldName must be not empty if IsForeignKey is true")
				default:
					foreignKeys += ", FOREIGN KEY(" + f.Name + ") REFERENCES " + f.RefTableName + " (" + f.RefFieldName + ") "
				}

			}
		}
	}
	sqlStmt += foreignKeys + ")"

	return sqlStmt, nil
}

func createTable(db *sql.DB, name string, fields ...Field) error {
	sqlStmt, err := makeSqlStmtCreateTable(name, fields...)

	if err != nil {
		return err
	}

	createTableStatement, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("error while preparing query db: %v", err)
		return err
	}

	_, err = createTableStatement.Exec()
	if err != nil {
		log.Printf("error with executing the statement: %v", err)
		return err
	}

	return nil
}

func InitDB() (*sql.DB, error) {
	dsnURI := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", dsnURI)
	if err != nil {
		log.Printf("error while opening db: %v", err)
		return nil, err
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
		log.Printf("error couldn't create a db: %v", err)
		return nil, err
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
		log.Printf("error couldn't create a db: %v", err)
		return nil, err
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
		log.Printf("error couldn't create a db: %v", err)
		return nil, err
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
		log.Printf("error couldn't create a db: %v", err)
		return nil, err
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
		log.Printf("error couldn't create a db: %v", err)
		return nil, err
	}

	return db, nil
}

// Creates a record in checks and sessions db.
func bindSessionAndCheck(db *sql.DB, sessionId, checkId int64) error {
	sql := `INSERT INTO checks_and_sessions (session_id, check_id) VALUES (?, ?)`
	statement, err := db.Prepare(sql)
	if err != nil {
		log.Printf("error while preparing query db: %v", err)
		return err
	}
	_, err = statement.Exec(sessionId, checkId)
	if err != nil {
		log.Printf("error with executing the statement: %v", err)
		return err
	}
	return nil
}

// adds a check in db and returns id of that chec if no error occured.
func AddCheck(c *models.Check, sessionId int64) (int64, error) {
	db, err := InitDB()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO checks (Name, Owner) VALUES (?, ?)")
	if err != nil {
		log.Printf("error while preparing query db: %v", err)
		return -1, err
	}

	res, err := statement.Exec(c.Name, c.Owner)
	if err != nil {
		log.Printf("error with executing the statement: %v", err)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}
	if err := bindSessionAndCheck(db, sessionId, id); err != nil {
		log.Printf("error binding check with session: %v", err)
		return id, err
	}
	return id, nil
}

// adds items to db and returns whatever error happened
func AddItems(items ...models.Item) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	statement, err := db.Prepare("INSERT INTO items (check_id, Name, Owner, Price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	for _, item := range items {
		if _, err := statement.Exec(item.CheckId, item.Name, item.Owner, item.Price); err != nil {
			log.Printf("error adding item: %v", err)
		}
	}

	return nil
}

func addNewSession(db *sql.DB) (int64, error) {
	statement, err := db.Prepare("INSERT INTO sessions (opened_at, is_open) VALUES (?, ?)")
	if err != nil {
		log.Printf("error while preparing query db: %v", err)
		return -1, err
	}

	res, err := statement.Exec(time.Now().Format(time.DateTime), "true")
	if err != nil {
		log.Printf("error with executing the statement: %v", err)
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}

	return id, nil
}

// Returns current session id. If there is no session open or present, then it's being created.
func GetSessionId() (int64, error) {
	db, err := InitDB()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	row := db.QueryRow(`SELECT id FROM sessions WHERE is_open = ?`, "true")
	var (
		id int64 = -1
	)

	err = row.Scan(&id)

	if err == sql.ErrNoRows {
		id, errAdd := addNewSession(db)
		if errAdd != nil {
			return -1, errAdd
		}
		return id, nil
	}

	return id, err
}

// Finishes a session with given id. Means setting is_open to false
// and closed_at to time.Now(DateTime).
func FinishSession(id int64) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := `UPDATE sessions SET is_open = ?, closed_at = ? WHERE id = ?`
	_, err = db.Exec(sql, "false", time.Now().Format(time.DateTime), id)
	if err != nil {
		return err
	}

	return nil
}

func CreateTotal(st *models.SessionTotal) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	sql := `INSERT INTO totals (session_id, total, recipient, amount) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(sql, st.SessionId, st.Total, st.Recipient, st.Amount)
	if err != nil {
		return err
	}

	return nil
}

func getAllCheckIdsWithSessionId(db *sql.DB, sessionId int64) ([]int64, error) {

	sql := `SELECT check_id FROM checks_and_sessions WHERE session_id = ?`

	rows, err := db.Query(sql, sessionId)
	if err != nil {
		return nil, err
	}
	var checkIds []int64
	var id int64
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return checkIds, err
		}
		checkIds = append(checkIds, id)
	}

	return checkIds, nil
}

func getCheckWithId(db *sql.DB, checkId int64) (*models.Check, error) {

	sql := `SELECT Name, Owner FROM checks WHERE id = ?`

	row := db.QueryRow(sql, checkId)

	c := models.Check{}
	if err := row.Scan(&c.Name, &c.Owner); err != nil {
		return nil, err
	}

	return &c, nil
}

func getItemsForCheckId(db *sql.DB, checkId int64) ([]models.Item, error) {
	sql := `SELECT * FROM items WHERE check_id = ?`
	rows, err := db.Query(sql, checkId)
	if err != nil {
		return nil, err
	}
	var items []models.Item
	for rows.Next() {
		i := models.Item{}
		if err := rows.Scan(&i.Id, &i.CheckId, &i.Name, &i.Owner, &i.Price); err != nil {
			return items, err
		}
		items = append(items, i)
	}

	return items, nil
}

func GetCheckWithItemsForId(checkId int64) (*models.CheckWithItems, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	check, err := getCheckWithId(db, checkId)
	if err != nil {
		return nil, err
	}

	items, err := getItemsForCheckId(db, checkId)
	if err != nil {
		return nil, err
	}
	res := &models.CheckWithItems{Id: checkId}
	res.SetCheck(check)
	res.SetItems(items)

	return res, nil
}

func GetAllChecksWithItemsForSesssionId(sessionId int64) ([]*models.CheckWithItems, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	checkIds, err := getAllCheckIdsWithSessionId(db, sessionId)
	if err != nil {
		return nil, err
	}

	var checksWithItems []*models.CheckWithItems
	for _, checkId := range checkIds {
		check, err := getCheckWithId(db, checkId)
		if err != nil {
			return nil, err
		}

		items, err := getItemsForCheckId(db, checkId)
		if err != nil {
			return nil, err
		}
		c := &models.CheckWithItems{Id: checkId}
		c.SetCheck(check)
		c.SetItems(items)
		checksWithItems = append(checksWithItems, c)
	}

	return checksWithItems, nil
}

// Returns a slice of session-totals from table 'totals'.
// Function may return not all totals, and an error.
func GetAllSessionTotals() ([]*models.SessionTotal, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sql := `SELECT * FROM totals`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	var totals []*models.SessionTotal

	for rows.Next() {
		var id int64
		total := new(models.SessionTotal)
		if err := rows.Scan(&id, &total.SessionId, &total.Total, &total.Recipient, &total.Amount); err != nil {
			return totals, err
		}

		s, err := getSessionById(db, total.SessionId)
		if err != nil {
			return totals, err
		}

		total.SetSession(s)
		totals = append(totals, total)
	}

	return totals, nil
}

func getSessionById(db *sql.DB, sessionId int64) (*models.Session, error) {
	var (
		id        int64
		opened_at string
		closed_at string
		is_open   string
	)
	sql := `SELECT * FROM sessions WHERE id = ?`
	row := db.QueryRow(sql, sessionId)

	if err := row.Scan(&id, &opened_at, &closed_at, &is_open); err != nil {
		return nil, err
	}

	s := new(models.Session)

	s.Id = id

	if openedAt, err := time.Parse(time.DateTime, opened_at); err == nil {
		s.OpenedAt = &openedAt
	} else {
		return nil, err
	}

	if closedAt, err := time.Parse(time.DateTime, closed_at); err == nil {
		s.ClosedAt = &closedAt
	} else {
		return nil, err
	}

	if isOpen, err := strconv.ParseBool(is_open); err == nil {
		s.IsOpen = isOpen
	} else {
		return nil, err
	}

	return s, nil
}

// Returns [models.Session] with provided id from 'sessions' table.
// Public wrapper for getSessionById, that sets, up a db connection and passes it to getSessionById.
func GetSessionById(id int64) (*models.Session, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return getSessionById(db, id)
}

func alterCheck(db *sql.DB, check *models.Check) error {
	if check.Name == "" && check.Owner == "" {
		return fmt.Errorf("expected at least one non-empty param, but provided: %+v", *check)
	}

	var sql string
	var err error
	switch {
	case check.Name == "" && check.Owner != "":
		sql = `UPDATE checks SET Owner = ? WHERE id = ?`
		_, err = db.Exec(sql, check.Owner, check.Id)
	case check.Name != "" && check.Owner == "":
		sql = `UPDATE checks SET Name = ? WHERE id = ?`
		_, err = db.Exec(sql, check.Name, check.Id)
	default:
		sql = `UPDATE checks SET Name = ?, Owner = ? WHERE id = ?`
		_, err = db.Exec(sql, check.Name, check.Owner, check.Id)
	}

	return err
}

func EditCheckName(checkId int64, newName string) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	newCheck := &models.Check{Id: checkId, Name: newName}

	return alterCheck(db, newCheck)
}

func EditCheckOwner(checkId int64, newOwner string) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	newCheck := &models.Check{Id: checkId, Owner: newOwner}

	return alterCheck(db, newCheck)
}

// TODO change params to Check struct
func EditCheck(checkId int64, newName string, newOwner string) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	newCheck := &models.Check{Id: checkId, Name: newName, Owner: newOwner}

	return alterCheck(db, newCheck)
}

func alterItem(db *sql.DB, item *models.Item) error {
	if item.Id == 0 {
		return errors.New("item.Id must be set, but provided")
	}

	if item.Name == "" && item.Owner == "" && item.Price == 0 {
		return fmt.Errorf("expected at least one non-empty param, but provided: %+v", *item)
	}

	// TODO: is this any good???
	var sql string
	var err error
	switch {
	case item.Name == "" && item.Owner != "" && item.Price != 0:
		sql = `UPDATE items SET Owner = ?, Price = ? WHERE id = ?`
		_, err = db.Exec(sql, item.Owner, item.Price, item.Id)
	case item.Name != "" && item.Owner == "" && item.Price != 0:
		sql = `UPDATE items SET Name = ?, Price = ? WHERE id = ?`
		_, err = db.Exec(sql, item.Name, item.Price, item.Id)
	case item.Name != "" && item.Owner != "" && item.Price == 0:
		sql = `UPDATE items SET Name = ?, Owner = ? WHERE id = ?`
		_, err = db.Exec(sql, item.Name, item.Owner, item.Id)
	default:
		sql = `UPDATE items SET Name = ?, Owner = ?, Price = ? WHERE id = ?`
		_, err = db.Exec(sql, item.Name, item.Owner, item.Price, item.Id)
	}

	return err
}

// TODO change params to Item struct
// EditItem changes item with given id in db.
func EditItem(itemId int64, newName string, newOwner string, newPrice float64) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	newItem := &models.Item{Id: itemId, Name: newName, Owner: newOwner, Price: newPrice}
	// perhaps select item before altering. if nothing's new, then do nothing
	return alterItem(db, newItem)
}
