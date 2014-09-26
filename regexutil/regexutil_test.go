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

	//fmt.Println("Pass One == > table => ", table, " , columns => ", columns)

}

func TestSecondPassParseCreateStatement(t *testing.T) {
	_, columns, err := firstPassParseCreateStatement(sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA)
	if err != nil {
		t.Error(err)
	}

	//fmt.Println("Pass Two == > table => ", table, " , columns => ", columns)

	_, err = secondPassParseCreateStatement(columns)

	if err != nil {
		t.Error(err)
	}

	_, columns, err = firstPassParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_B)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("new pass.")
	//fmt.Println("Pass Two == > table => ", table, " , columns => ", columns)

	c, err := secondPassParseCreateStatement(columns)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(c)

	_, columns, err = firstPassParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_C)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("new pass 2.")
	//fmt.Println("Pass Two == > table => ", table, " , columns => ", columns)

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

	fmt.Println(table, columnArrayMap)

	cam := []map[string]string{{"": "X INTEGER", "columntype": "INTEGER", "columnname": "X", "constraints": ""}, "": "Y TEXT PRIMARY KEY", "columntype": "TEXT", "columnname": "Y", "constraints": "PRIMARY KEY"}

	if !reflect.DeepEqual(columnArrayMap, cam) {
		fmt.Println(cam)
		fmt.Println(columnArrayMap)
	}

}
