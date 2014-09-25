package sqlite

import (
	"fmt"
	"testing"
)

const SQLITE3_TEST_TABLE_NAME = "FOO"
const SQLITE3_TEST_CREATE_SCHEMA = "CREATE TABLE FOO (NAME TEXT, ADDRESS TEXT, ID INTEGER, CREATED_TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"
const SQLITE3_TEST_DROP_SCHEMA = "DROP TABLE FOO"
const SQLITE3_TEST_INSERT = "INSERT INTO FOO (NAME,ADDRESS) VALUES (?,?);"
const SQLITE3_TEST_SELECT = "SELECT NAME,ADDRESS FROM FOO;"

func TestNewConnection(t *testing.T) {
	_, err := NewConnection()
	if err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}
	err = conn.Close()
	if err != nil {
		t.Error(conn)
	}

}

func TestCheckTableExists(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}
	exists, err := conn.CheckTableExists(SQLITE3_TEST_TABLE_NAME)
	if err != nil {
		t.Error(err)
	}

	if exists != false {
		t.Error("error => table " + SQLITE3_TEST_TABLE_NAME + "should not exist")
	}

	err = conn.Close()
	if err != nil {
		t.Error(conn)
	}
}

func TestInitTable(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}

	err = conn.InitTable(SQLITE3_TEST_CREATE_SCHEMA, SQLITE3_TEST_TABLE_NAME)

	if err != nil {
		t.Error(err)
	}

	_, err = conn.db.Exec(SQLITE3_TEST_DROP_SCHEMA)

	if err != nil {
		t.Error(err)
	}

	err = conn.Close()

	if err != nil {
		t.Error(err)
	}

}

func TestFoo(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	err = conn.InitTable(SQLITE3_TEST_CREATE_SCHEMA, SQLITE3_TEST_TABLE_NAME)

	if err != nil {
		t.Error(err)
	}

	err = conn.Insert(SQLITE3_TEST_INSERT, "parag", "austin")

	if err != nil {
		t.Error(err)
	}

	rows, err := conn.db.Query(SQLITE3_TEST_SELECT)

	defer rows.Close()

	cols, err := rows.Columns()

	if err != nil {
		t.Error(err)
	}
	fmt.Println("columns are : ", cols)

}
