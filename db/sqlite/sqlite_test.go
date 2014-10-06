package sqlite

import (
	"fmt"
	"reflect"
	"testing"

	"../godata"
)

const SQLITE3_TEST_TABLE_NAME = "FOO"
const SQLITE3_TEST_CREATE_SCHEMA = "CREATE TABLE FOO (NAME TEXT, ADDRESS TEXT, ID INTEGER, CREATED_TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"
const SQLITE3_TEST_DROP_SCHEMA = "DROP TABLE FOO"
const SQLITE3_TEST_INSERT = "INSERT INTO FOO (NAME,ADDRESS) VALUES (?,?);"
const SQLITE3_TEST_SELECT = "SELECT NAME,ADDRESS FROM FOO;"
const SQLITE3_TEST_CREATE_SCHEMA_FAIL = "CREATE TABLE FOO (NAME TEXT, ADDRESS TEXT, ID2 INTEGER, CREATED_TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"

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

func TestInsert(t *testing.T) {
	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	err = conn.InitTable(SQLITE3_TEST_CREATE_SCHEMA, SQLITE3_TEST_TABLE_NAME)

	if err != nil {
		t.Error(err)
	}

	n, err := conn.Insert(SQLITE3_TEST_INSERT, "parag", "austin")

	if err != nil {
		t.Error(err)
	}

	if n != 1 {
		t.Error("Row was not inserted properly ; n = ", n)
	}

	_, err = conn.db.Exec(SQLITE3_TEST_DROP_SCHEMA)

	if err != nil {
		t.Error(err)
	}

}

func TestSelect(t *testing.T) {

	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}

	defer conn.Close()

	err = conn.InitTable(SQLITE3_TEST_CREATE_SCHEMA, SQLITE3_TEST_TABLE_NAME)

	if err != nil {
		t.Error(err)
	}

	n, err := conn.Insert(SQLITE3_TEST_INSERT, "parag", "austinYO")

	if n != 1 {
		t.Error("statement is not inserted n => ", n)
	}

	//fmt.Println("n => ", n)

	if err != nil {
		t.Error(err)
	}

	goselect, err := conn.Select(SQLITE3_TEST_SELECT)

	if err != nil {
		t.Error(err)
	}

	//_:= goselect.GetColumns()
	//ows := goselect.GetRows()

	fmt.Println("YO YO", goselect)

	_, err = conn.db.Exec(SQLITE3_TEST_DROP_SCHEMA)

	if err != nil {
		t.Error(err)
	}

}

func TestMetaGoTable(t *testing.T) {

	conn, err := NewConnection()
	if err != nil {
		t.Error(err)
	}

	err = conn.InitTable(SQLITE3_TEST_CREATE_SCHEMA, SQLITE3_TEST_TABLE_NAME)

	if err != nil {
		t.Error(err)
	}

	gmeta, err := conn.MetaGoTable(SQLITE3_TEST_TABLE_NAME)

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

	gtest, err := godata.NewMetaTableFromCreateStatement(SQLITE3_TEST_CREATE_SCHEMA)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(gmeta, gtest) {
		t.Error("gmeta and gtest are not the same")
		t.Error("gmeta column values are ", gmeta.GetColumns())
		t.Error("gtest column values are ", gtest.GetColumns())
	}

}
