package godata

import "testing"

const TEST_SQLITE3_CREATE_STATEMENT_A = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY)"
const TEST_SQLITE3_CREATE_STATEMENT_B = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_C = "CREATE TABLE FOO ( X INTEGER, PRIMARY KEY Y, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_D = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL, T DEFAULT  CURRENT_TIMESTAMP)"

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

	columns := meta.columns.GetAll()

	if columns[0].name != "X" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[0].name)
	}

	if columns[0].sqlType != "INTEGER" {
		t.Error("Error column sql type is wrong")
	}

	if columns[0].goType != "int" {
		t.Error("Error column go type  is wrong")
	}

	if columns[1].name != "Y" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[1].name)
	}

	if columns[1].sqlType != "TEXT" {
		t.Error("Error column sql type is wrong")
	}

	if columns[1].goType != "string" {
		t.Error("Error column go type  is wrong")
	}

	//TEST B

	meta, err = NewMetaTableFromCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_B)
	if err != nil {
		t.Error(err)
		return
	}

	if meta.name != "FOO" {
		t.Error("table name not stored properly")
	}

	if meta.dbVendor != "sqlite3" {
		t.Error("table name not stored properly")
	}

	columns = meta.columns.GetAll()

	if columns[0].name != "X" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[0].name)
	}

	if columns[0].sqlType != "INTEGER" {
		t.Error("Error column sql type is wrong")
	}

	if columns[0].goType != "int" {
		t.Error("Error column go type  is wrong")
	}

	if columns[1].name != "Y" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[1].name)
	}

	if columns[1].sqlType != "TEXT" {
		t.Error("Error column sql type is wrong")
	}

	if columns[1].goType != "string" {
		t.Error("Error column go type  is wrong")
	}

	if columns[2].name != "Z" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[0].name)
	}

	if columns[2].sqlType != "BETA" {
		t.Error("Error column sql type is wrong")
	}

	if columns[2].goType != "" {
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
		return
	}

	if meta.name != "FOO" {
		t.Error("table name not stored properly")
	}

	if meta.dbVendor != "sqlite3" {
		t.Error("table name not stored properly")
	}

	columns = meta.columns.GetAll()

	if columns[0].name != "X" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[0].name)
	}

	if columns[0].sqlType != "INTEGER" {
		t.Error("Error column sql type is wrong")
	}

	if columns[0].goType != "int" {
		t.Error("Error column go type  is wrong")
	}

	if columns[1].name != "Y" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[1].name)
	}

	if columns[1].sqlType != "TEXT" {
		t.Error("Error column sql type is wrong")
	}

	if columns[1].goType != "string" {
		t.Error("Error column go type  is wrong")
	}

	if columns[2].name != "Z" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[0].name)
	}

	if columns[2].sqlType != "BETA" {
		t.Error("Error column sql type is wrong")
	}

	if columns[2].goType != "" {
		t.Error("Error column go type  is wrong")
	}

	if columns[3].name != "T" {
		t.Error("Error column name is wrong")
		t.Error("Column name is => " + columns[0].name)
	}

	if columns[3].sqlType != "TIMESTAMP" {
		t.Error("Error column sql type is wrong")
	}

	if columns[3].goType != "time.Time" {
		t.Error("Error column go type  is wrong")
	}

}
