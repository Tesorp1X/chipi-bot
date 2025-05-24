package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

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

func CreateDB() (*sql.DB, error) {
	dsnURI := os.Getenv("DB_PATH")
	db, err := sql.Open("sqlite3", dsnURI)
	if err != nil {
		log.Printf("error while opening db: %v", err)
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
			Type:       "STRING",
			IsNullable: false,
		},
		{
			Name:       "Owner",
			Type:       "STRING",
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
			Type: "STRING",
		},
		{
			Name: "Owner",
			Type: "STRING",
		},
		{
			Name: "Price",
			Type: "FLOAT",
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
			Name: "check_id",
			Type: "INTEGER",

			IsForeignKey: true,
			RefTableName: "checks",
			RefFieldName: "id",
		},
		{
			Name: "total",
			Type: "FLOAT",
		},
		{
			Name: "recipient",
			Type: "STRING",
		},
		{
			Name: "amount",
			Type: "FLOAT",
		},
	}
	if err = createTable(db, "totals", totalsField...); err != nil {
		log.Printf("error couldn't create a db: %v", err)
		return nil, err
	}

	return db, nil
}

// adds a check in db and returns id of that chec if no error occured.
func AddCheck(c *models.Check) (int64, error) {
	db, err := CreateDB()
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

	return id, nil
}

// adds items to db and returns whatever error happened
func AddItems(items ...models.Item) error {
	db, err := CreateDB()
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
