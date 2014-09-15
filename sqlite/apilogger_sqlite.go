package sqlite

import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "fmt"
import "errors"

const file = "throttle.db"

const check_if_table_exists = "SELECT name FROM sqlite_master WHERE type=? AND name=?;"
const create_table_throttle = "create table throttle(id integer not null primary key AUTOINCREMENT,throttle_message text, throttle_response_timestamp text, insert_timestamp default CURRENT_TIMESTAMP);"

type connection struct {
	db     *sql.DB
	dbName string
}

func NewConnection(fileName string) (*connection, error) {

	if fileName == "" {

		fileName = file
	}

	db, err := sql.Open("sqlite3", "file:"+file)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &connection{db, fileName}, nil

}

func (Conn *connection) Close() error {
	return Conn.Close()
}

func (Conn *connection) checkTableExists(tableName string) (bool, error) {
	rows, err := Conn.db.Query(check_if_table_exists, "table", tableName)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return rows.Next(), nil

}

func (Conn *connection) createTable(createStatement string, tableName string) error {
	if createStatement == "" {
		return errors.New("Create Statement cannot be empty")

	}

	rowsExist, err := Conn.checkTableExists(tableName)

	if err != nil {
		return err
	}

	if !rowsExist {

		_, err = Conn.db.Exec(createStatement)
		if err != nil {
			return err
		}

		rowsExist, err = Conn.checkTableExists(tableName)

		if err != nil {
			return err
		}

		if rowsExist != true {
			return errors.New("Error => Could not create table " + tableName)
		}
	}

	return nil
}

/*
func (Conn *connection) InsertAPIThrottle(throttle_message, throttle_response_timestamp string) error {
	db, err := open_db()
	defer db.Close()

	if err != nil {
		return err
	}
	err = create_db(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}

	stmt, err := tx.Prepare("insert into throttle(throttle_message,throttle_response_timestamp) values(?, ?)")

	if err != nil {
		fmt.Println("Could not prepare statment = >", err)
		return "", err
	}

	defer stmt.Close()

	return nil
} */
