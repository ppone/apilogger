package godata

import "testing"

const TEST_SQLITE3_CREATE_STATEMENT_A = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY)"
const TEST_SQLITE3_CREATE_STATEMENT_B = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_C = "CREATE TABLE FOO ( X INTEGER, PRIMARY KEY Y, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_D = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL, T DEFAULT  CURRENT_TIMESTAMP)"

func TestPointerConvertor(t *testing.T) {

	var m interface{}
	s := "hello world"
	m = &s

	_, err := PointerConvertor(m)

	if err != nil {
		t.Error(err)
	}

	var n interface{}

	var u int64
	u = 2342342341111
	n = &u

	_, err = PointerConvertor(n)

	if err == nil {
		t.Error(err)
	}

}

func TestNewMetaTableFromCreateStatement(t *testing.T) {

	//TEST A
	meta, err := NewMetaTableFromCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_A)
	if err != nil {
		t.Error(err)
	}

	if meta.name != "FOO" {
		t.Error("table name not stored properly")
	}

	if meta.dbVendor != "sqlite3" {
		t.Error("table name not stored properly")
	}

	columns := meta.columns

	colName := "X"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "INTEGER" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "int" {
		t.Error("Error column go type  is wrong")
	}

	colName = "Y"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "TEXT" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "string" {
		t.Error("Error column go type  is wrong")
	}

	//TEST B

	meta, err = NewMetaTableFromCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_B)

	if err != nil {
		t.Error(err)
	}

	if meta.name != "FOO" {
		t.Error("table name not stored properly")
	}

	if meta.dbVendor != "sqlite3" {
		t.Error("table name not stored properly")
	}

	columns = meta.columns

	colName = "X"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "INTEGER" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "int" {
		t.Error("Error column go type  is wrong")
	}

	colName = "Y"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "TEXT" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "string" {
		t.Error("Error column go type  is wrong")
	}

	colName = "Z"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "BETA" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "" {
		t.Error("Error column go type  is wrong")
	}

	//TEST C

	meta, err = NewMetaTableFromCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_C)
	if err == nil {
		t.Error(err)
		return
	}

	//TEST D
	meta, err = NewMetaTableFromCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_D)

	if err != nil {
		t.Error(err)
	}

	if meta.name != "FOO" {
		t.Error("table name not stored properly")
	}

	if meta.dbVendor != "sqlite3" {
		t.Error("table name not stored properly")
	}

	columns = meta.columns

	colName = "X"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "INTEGER" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "int" {
		t.Error("Error column go type  is wrong")
	}

	colName = "Y"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "TEXT" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "string" {
		t.Error("Error column go type  is wrong")
	}

	colName = "Z"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "BETA" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "" {
		t.Error("Error column go type  is wrong")
	}

	colName = "T"

	if _, ok := columns[colName]; !ok {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + colName)
	}

	if columns[colName].SqlType() != "TIMESTAMP" {
		t.Error("Error column sql type is wrong")
	}

	if columns[colName].GoType() != "time.Time" {
		t.Error("Error column go type  is wrong")
	}

}
