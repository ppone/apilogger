package sqlite

import "database/sql"
import "../sqlconstants"
import _ "github.com/mattn/go-sqlite3"
import "fmt"
import "errors"

type connection struct {
	db     *sql.DB
	dbName string
}

func (Conn *connection) Close() error {
	return Conn.Close()
}

func (Conn *connection) CheckTableExists(tableName string) (bool, error) {
	rows, err := Conn.db.Query(check_if_table_exists, "table", tableName)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return rows.Next(), nil

}

func (Conn *connection) CreateTable(createStatement string, tableName string) error {
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

	if err != nil {
		return err
	}

	defer stmt.Close()

	return nil
}
