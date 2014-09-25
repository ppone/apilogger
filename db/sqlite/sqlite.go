package sqlite

import "database/sql"
import "../sqlconstants"
import _ "github.com/mattn/go-sqlite3"
import "fmt"
import "errors"

type connection struct {
	db *sql.DB
}

func NewConnection() (*connection, error) {

	db, err := sql.Open(sqlconstants.CurrentVendor(), sqlconstants.SQLITE3_CONNECTION_STRING)

	if err != nil {
		return nil, err
	}

	return &connection{db}, nil

}

func (Conn *connection) Close() error {
	return Conn.db.Close()
}

func (Conn *connection) CheckTableExists(tableName string) (bool, error) {
	rows, err := Conn.db.Query(sqlconstants.SQLITE3_CHECK_IF_TABLE_EXISTS, "table", tableName)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return rows.Next(), nil

}

/*
func getColumNamesFromSelectQuery(query string) ([]sqlData, error) {

}*/
/*
func (Conn *connection) Select(query string) ([]sqlData, error) {
	rows, err := Conn.db.Query(query, data)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//cols, err := rows.Columns()
	_, err = rows.Columns()

	if err != nil {
		return nil, err
	}

	return nil, nil

}*/

func (Conn *connection) InitTable(createStatement string, tableName string) error {
	if createStatement == "" {
		return errors.New("Create Statement cannot be empty")

	}

	rowsExist, err := Conn.CheckTableExists(tableName)

	if err != nil {
		return err
	}

	if !rowsExist {

		_, err = Conn.db.Exec(createStatement)
		if err != nil {
			return err
		}

		rowsExist, err = Conn.CheckTableExists(tableName)

		if err != nil {
			return err
		}

		if rowsExist != true {
			return errors.New("Error => Could not create table " + tableName)
		}
	}

	return nil
}

func (Conn *connection) Insert(insertStatement string, data ...interface{}) error {

	if insertStatement == "" {
		return errors.New("Insert Statement cannot be blank")

	}

	if len(data) == 0 {
		return errors.New("data cannot be empty")

	}

	tx, err := Conn.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}

	stmt, err := tx.Prepare(insertStatement)

	if err != nil {
		fmt.Println("Could not prepare statment = >", err)
		return err
	}

	_, err = stmt.Exec(data)

	err = tx.Commit()

	defer stmt.Close()

	return err

}
