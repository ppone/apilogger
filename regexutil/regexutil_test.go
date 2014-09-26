package regexutil

import (
	"fmt"
	"reflect"
	"testing"

	"../db/sqlconstants"
)

const TEST_SQLITE3_CREATE_STATEMENT_A = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY)"
const TEST_SQLITE3_CREATE_STATEMENT_B = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_C = "CREATE TABLE FOO ( X INTEGER, PRIMARY KEY Y, Z BETA NOT NULL)"

func TestReturnNameGroupValueMap(t *testing.T) {
	_, err := ParseAndReturnNameGroupValueMap(SQLITE3_CREATE_TABLE_FIRST_PASS_PARSER, TEST_SQLITE3_CREATE_STATEMENT_A)

	if err != nil {
		t.Error(err)
	}

}

func TestFirstPassParseCreateStatement(t *testing.T) {
	_, _, err := firstPassParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_A)
	if err != nil {
		t.Error(err)
	}

}

func TestSecondPassParseCreateStatement(t *testing.T) {
	_, columns, err := firstPassParseCreateStatement(sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA)
	if err != nil {
		t.Error(err)
	}

	_, err = secondPassParseCreateStatement(columns)

	if err != nil {
		t.Error(err)
	}

	_, columns, err = firstPassParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_B)
	if err != nil {
		t.Error(err)
	}

	_, err = secondPassParseCreateStatement(columns)

	if err != nil {
		t.Error(err)
	}

	_, columns, err = firstPassParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_C)
	if err != nil {
		t.Error(err)
	}

	_, err = secondPassParseCreateStatement(columns)

	if err == nil {
		t.Error("Error should have been caught")
	}

}

func TestParseCreateStatement(t *testing.T) {

	table, columnArrayMap, err := ParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_A)

	if err != nil {
		t.Error(err)
	}

	if table != "FOO" {
		t.Error("error table name is incorrect")
	}

	cam := []map[string]string{}
	c1 := map[string]string{"": "X INTEGER", "columntype": "INTEGER", "columnname": "X", "constraints": ""}
	c2 := map[string]string{"": "Y TEXT PRIMARY KEY", "columntype": "TEXT", "columnname": "Y", "constraints": "PRIMARY KEY"}

	cam = append(cam, c1)
	cam = append(cam, c2)

	if !reflect.DeepEqual(columnArrayMap, cam) {
		t.Error("ParseCreateStatement created unexpected results => test case ->", TEST_SQLITE3_CREATE_STATEMENT_A)
		t.Error("test data =>", cam)
		t.Error("returned data =>", columnArrayMap)
	}

	table, columnArrayMap, err = ParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_B)

	if err != nil {
		t.Error(err)
	}

	if table != "FOO" {
		t.Error("error table name is incorrect")
	}

	cam = []map[string]string{}
	c1 = map[string]string{"": "X INTEGER", "columntype": "INTEGER", "columnname": "X", "constraints": ""}
	c2 = map[string]string{"": "Y TEXT PRIMARY KEY", "columntype": "TEXT", "columnname": "Y", "constraints": "PRIMARY KEY"}
	c3 := map[string]string{"": "Z BETA NOT NULL", "columntype": "BETA", "columnname": "Z", "constraints": "NOT NULL"}

	cam = append(cam, c1)
	cam = append(cam, c2)
	cam = append(cam, c3)

	if !reflect.DeepEqual(columnArrayMap, cam) {
		t.Error("ParseCreateStatement created unexpected results => test case ->", TEST_SQLITE3_CREATE_STATEMENT_B)
		t.Error("test data =>", cam)
		t.Error("returned data =>", columnArrayMap)
	}

	table, columnArrayMap, err = ParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_C)

	if err == nil {
		t.Error("ERROR should be thrown")
	}

	fmt.Println(err)

}
