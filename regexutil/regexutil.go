package regexutil

import (
	"../db/sqlconstants"
	"regexp"
	"strings"
)
import "errors"

const TABLE = "table"
const COLUMNS = "columns"
const COLUMNNAME = "columnname"
const COLUMNTYPE = "columntype"
const CONSTRAINTS = "constraints"
const FIRSTWORD = "firstword"
const SQLITE3_CREATE_TABLE_FIRST_PASS_PARSER = `^(?i)\s*create\s+table\s+(?P<table>[a-z]*)\s*\((?P<columns>\s*(?:[a-z]+\s*.*\s*,\s*)*(?:[a-z]+\s*[^,]*\s*)+\s*)\)\s*[;]?\s*$`
const SQLITE3_CREATE_TABLE_SECOND_PASS_COLUMN_PARSER = `^(?i)\s*(?P<columnname>\w+)\s+(?P<columntype>INTEGER|TEXT|REAL|BLOB|NULL)?(?P<constraints>.*)$`
const SQLITE3_GET_BACK_FIRST_WORD = `^(?i)\s*(?P<firstword>\w+)\b*.*$`
const SQLITE3_DEFAULT_TIMESTAMP = `^(?i)\s*(?:DEFAULT\s*CURRENT_TIMESTAMP$|DEFAULT\s*CURRENT_TIMESTAMP\s+.*$)`
const SQLITE3_GET_TABLE_FROM_SELECT = `^(?i)\s*select\s+.*from\s+(?P<tablename>[a-z]+)(?:;|\s+.*;?)?$`

//const SQLITE3_GET_TABLE_FROM_SELECT = `^(?i)\s*select\s+.*from\s+(?P<tablename>[a-z]+)\s+.*$`

func getTableNameFromSelect() string {
	if sqlconstants.CurrentVendor() == "sqlite3" {
		return SQLITE3_GET_TABLE_FROM_SELECT
	}

	return ""
}

func getDefaultTimestamp() string {
	if sqlconstants.CurrentVendor() == "sqlite3" {
		return SQLITE3_DEFAULT_TIMESTAMP
	}

	return ""
}

func getBackFirstWord() string {
	if sqlconstants.CurrentVendor() == "sqlite3" {
		return SQLITE3_GET_BACK_FIRST_WORD
	}

	return ""
}

func createTableFirstPassParser() string {
	if sqlconstants.CurrentVendor() == "sqlite3" {
		return SQLITE3_CREATE_TABLE_FIRST_PASS_PARSER
	}

	return ""
}

func createTableSecondPassColumnParser() string {
	if sqlconstants.CurrentVendor() == "sqlite3" {
		return SQLITE3_CREATE_TABLE_SECOND_PASS_COLUMN_PARSER
	}

	return ""
}

func TableNameFromSelect(selectStmt string) (string, error) {

	stmt := strings.TrimRight(selectStmt, ";")
	m, err := ParseAndReturnNameGroupValueMap(getTableNameFromSelect(), stmt)

	if err != nil {
		return "", err
	}

	tableName, ok := m["tablename"]

	if !ok {
		return "", errors.New("Error: could not get table name from select statment")
	}

	return tableName, nil

}

func ParseAndReturnNameGroupValueMap(regexParser, statementToParse string) (map[string]string, error) {
	re, err := regexp.Compile(regexParser)

	if err != nil {
		return nil, err
	}

	groupvalues := re.FindStringSubmatch(statementToParse)
	groupnames := re.SubexpNames()

	if len(groupvalues) == 0 {
		return nil, errors.New("Statement did not match parser")
	}

	if len(groupnames) != len(groupvalues) {
		return nil, errors.New("Mismatch between group name and group values")
	}

	groupnamevaluemap := make(map[string]string)

	for i, v := range groupnames {
		groupnamevaluemap[v] = groupvalues[i]
	}

	return groupnamevaluemap, nil

}

func fillCustomColumnType(parsedColumnMap map[string]string) error {

	isTimestamp, err := regexp.MatchString(getDefaultTimestamp(), parsedColumnMap["constraints"])
	if err != nil {
		return err
	}

	if isTimestamp {
		parsedColumnMap[COLUMNTYPE] = "TIMESTAMP"

	} else if parsedColumnMap[COLUMNTYPE] == "" {
		firstWordinConstraints, err := ParseAndReturnNameGroupValueMap(getBackFirstWord(), parsedColumnMap["constraints"])
		if err != nil {
			return err
		}

		isconstraint, err := sqlconstants.IsSQLConstraint(firstWordinConstraints["firstword"])

		if err != nil {
			return err
		}

		if !isconstraint {
			parsedColumnMap[COLUMNTYPE] = firstWordinConstraints["firstword"]
			parsedColumnMap[CONSTRAINTS] = strings.Trim(strings.Replace(parsedColumnMap[CONSTRAINTS], firstWordinConstraints["firstword"], "", -1), " ")
		}

	}

	return nil

}

func ParseCreateStatement(statementToParse string) (string, []map[string]string, error) {

	replaceFunctionsInCreate, err := sqlconstants.CreateStatementFunctionsToReplace()

	if err != nil {
		return "", nil, err
	}

	for _, replaceFunction := range replaceFunctionsInCreate {
		statementToParse = strings.Replace(statementToParse, replaceFunction, "", -1)
	}

	table, columns, err := firstPassParseCreateStatement(statementToParse)
	if err != nil {
		return "", nil, err
	}
	columnArrayMap, err := secondPassParseCreateStatement(columns)

	if err != nil {
		return "", nil, err
	}

	return table, columnArrayMap, nil
}

func firstPassParseCreateStatement(statementToParse string) (string, string, error) {
	m, err := ParseAndReturnNameGroupValueMap(createTableFirstPassParser(), statementToParse)
	if err != nil {
		return "", "", err
	}

	return m[TABLE], m[COLUMNS], nil

}

func secondPassParseCreateStatement(columnsToParse string) ([]map[string]string, error) {

	columns := strings.Split(columnsToParse, ",")

	var columnmaparray []map[string]string

	for _, v := range columns {
		parsedColumnMap, err := ParseAndReturnNameGroupValueMap(createTableSecondPassColumnParser(), strings.Trim(v, " "))
		if err != nil {
			return nil, err
		}

		parsedColumnMap[COLUMNTYPE] = strings.Trim(parsedColumnMap[COLUMNTYPE], " ")
		parsedColumnMap[COLUMNNAME] = strings.Trim(parsedColumnMap[COLUMNNAME], " ")
		parsedColumnMap[CONSTRAINTS] = strings.Trim(parsedColumnMap[CONSTRAINTS], " ")

		if ok, err := sqlconstants.IsSQLConstraint(parsedColumnMap[COLUMNNAME]); err == nil {

			if ok {
				return nil, errors.New("column name cannot be blank")
			}

		} else {
			return nil, err
		}

		err = fillCustomColumnType(parsedColumnMap)

		if err != nil {
			return nil, err
		}

		columnmaparray = append(columnmaparray, parsedColumnMap)

	}

	return columnmaparray, nil

}
