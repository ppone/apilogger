package regexutil

import (
	"reflect"
	"testing"

	"../db/sqlconstants"
	"strings"
)

const TEST_SQLITE3_CREATE_STATEMENT_A = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY)"
const TEST_SQLITE3_CREATE_STATEMENT_B = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_C = "CREATE TABLE FOO ( X INTEGER, PRIMARY KEY Y, Z BETA NOT NULL)"
const TEST_SQLITE3_CREATE_STATEMENT_D = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY, Z BETA NOT NULL, T DEFAULT  CURRENT_TIMESTAMP)"

const TEST_SQLITE_SELECT_STATEMENT_A = "SELECT * FROM FOO WHERE X = 2 AND Y = 3"
const TEST_SQLITE_SELECT_STATEMENT_B = "SELECT NAME, ADDRESS FROM FOOBAR WHERE ID IS NOT NULL;"

const TEST_SQLITE_SELECT_STATEMENT_C = "SELECT NAME,ADDRESS FROM FOO;"

func TestTableNameFromSelect(t *testing.T) {
	tableName, err := TableNameFromSelect(TEST_SQLITE_SELECT_STATEMENT_A)

	if err != nil {
		t.Error(err)
	}

	if tableName != "FOO" {
		t.Error("table name does not equal ", "FOO")

	}

	tableName, err = TableNameFromSelect(TEST_SQLITE_SELECT_STATEMENT_B)

	if err != nil {
		t.Error(err)
	}

	if tableName != "FOOBAR" {
		t.Error("table name does not equal ", "FOO")
	}

	tableName, err = TableNameFromSelect(TEST_SQLITE_SELECT_STATEMENT_C)

	if err != nil {
		t.Error(err)
	}

	if tableName != "FOO" {
		t.Error("table name does not equal ", "FOO")
		t.Error("table name actually is ", tableName)
	}

}

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

func TestX(t *testing.T) {

	replaceFunctionsInCreate, err := sqlconstants.CreateStatementFunctionsToReplace()

	if err != nil {
		t.Error(err)
	}

	statementToParse := sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA

	for _, replaceFunction := range replaceFunctionsInCreate {
		statementToParse = strings.Replace(statementToParse, replaceFunction, "", -1)
	}

	_, columns, err := firstPassParseCreateStatement(statementToParse)
	if err != nil {
		t.Error(err)
	}

	_, err = secondPassParseCreateStatement(columns)

	if err != nil {
		t.Error(err)
	}

}

func TestSecondPassParseCreateStatement(t *testing.T) {

	replaceFunctionsInCreate, err := sqlconstants.CreateStatementFunctionsToReplace()

	if err != nil {
		t.Error(err)
	}

	statementToParse := sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA

	for _, replaceFunction := range replaceFunctionsInCreate {
		statementToParse = strings.Replace(statementToParse, replaceFunction, "", -1)
	}

	_, columns, err := firstPassParseCreateStatement(statementToParse)
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

	//CASE A

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

	//CASE B

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

	//CASE C

	table, columnArrayMap, err = ParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_C)

	if err == nil {
		t.Error("ERROR should be thrown")
	}

	//CASE D
	table, columnArrayMap, err = ParseCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_D)

	if err != nil {
		t.Error(err)
	}

	if table != "FOO" {
		t.Error("error table name is incorrect")
	}

	cam = []map[string]string{}
	c1 = map[string]string{"": "X INTEGER", "columntype": "INTEGER", "columnname": "X", "constraints": ""}
	c2 = map[string]string{"": "Y TEXT PRIMARY KEY", "columntype": "TEXT", "columnname": "Y", "constraints": "PRIMARY KEY"}
	c3 = map[string]string{"": "Z BETA NOT NULL", "columntype": "BETA", "columnname": "Z", "constraints": "NOT NULL"}
	c4 := map[string]string{"": "T DEFAULT  CURRENT_TIMESTAMP", "columntype": "TIMESTAMP", "columnname": "T", "constraints": "DEFAULT  CURRENT_TIMESTAMP"}

	cam = append(cam, c1)
	cam = append(cam, c2)
	cam = append(cam, c3)
	cam = append(cam, c4)

	if !reflect.DeepEqual(columnArrayMap, cam) {
		t.Error("ParseCreateStatement created unexpected results => test case ->", TEST_SQLITE3_CREATE_STATEMENT_D)
		t.Error("test data =>", cam)
		t.Error("returned data =>", columnArrayMap)
	}

}
