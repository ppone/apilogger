package regexutil

import (
	"fmt"
	"regexp"
	"strings"

	"../db/sqlconstants"
)
import "errors"

const SQLITE3_TABLE = "table"
const SQLITE3_COLUMNS = "columns"
const SQLITE3_COLUMNNAME = "columnname"
const SQLITE3_COLUMNTYPE = "columntype"
const SQLITE3_CONSTRAINTS = "constraints"
const FIRSTWORD = "firstword"
const SQLITE3_CREATE_TABLE_FIRST_PASS_PARSER = `^(?i)\s*create\s+table\s+(?P<table>[a-z]*)\s*\((?P<columns>\s*(?:[a-z]+\s*.*\s*,\s*)*(?:[a-z]+\s*[^,]*\s*)+\s*)\)\s*[;]?\s*$`
const SQLITE3_CREATE_TABLE_SECOND_PASS_COLUMN_PARSER = `^(?i)\s*(?P<columnname>\w+)\s+(?P<columntype>INTEGER|TEXT|REAL|BLOB|NULL)?(?P<constraints>.*)$`
const SQLITE3_GET_BACK_FIRST_WORD = `^(?i)\s*(?P<firstword>\w+)\b*.*$`

//const SQLITE3_CREATE_TABLE_SECOND_PASS_COLUMN_PARSER = `^(?i)\s*(?P<columnname>\w+)\s+.*?(?P<columntype>INTEGER|TEXT|REAL|BLOB|DEFAULT\s*CURRENT_TIMESTAMP|NULL).*$`
const tmp = `^(?i)\s*(?P<columnname>\w+)\s+([^ITRBDC]\w*\s+)*(?P<columntype>INTEGER|TEXT|REAL|BLOB|DEFAULT\s*CURRENT_TIMESTAMP).*$`

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
	if parsedColumnMap[SQLITE3_COLUMNTYPE] == "" {
		firstWordinConstraints, err := ParseAndReturnNameGroupValueMap(SQLITE3_GET_BACK_FIRST_WORD, parsedColumnMap["constraints"])
		if err != nil {
			return err
		}
		isconstraint, err := sqlconstants.IsSQLConstraint(firstWordinConstraints["firstword"])

		if err != nil {
			return err
		}

		if !isconstraint {
			fmt.Println("over riding column type with ", firstWordinConstraints["firstword"])
			parsedColumnMap[SQLITE3_COLUMNTYPE] = firstWordinConstraints["firstword"]
		}

	}

	return nil

}

func ParseCreateStatement(statementToParse string) (string, []map[string]string, error) {

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
	m, err := ParseAndReturnNameGroupValueMap(SQLITE3_CREATE_TABLE_FIRST_PASS_PARSER, statementToParse)
	if err != nil {
		return "", "", err
	}

	return m[SQLITE3_TABLE], m[SQLITE3_COLUMNS], nil

}

func secondPassParseCreateStatement(columnsToParse string) ([]map[string]string, error) {

	columns := strings.Split(columnsToParse, ",")

	var columnmaparray []map[string]string

	for _, v := range columns {
		parsedColumnMap, err := ParseAndReturnNameGroupValueMap(SQLITE3_CREATE_TABLE_SECOND_PASS_COLUMN_PARSER, strings.Trim(v, " "))
		if err != nil {
			return nil, err
		}

		parsedColumnMap[SQLITE3_COLUMNTYPE] = strings.Trim(parsedColumnMap[SQLITE3_COLUMNTYPE], " ")
		parsedColumnMap[SQLITE3_COLUMNNAME] = strings.Trim(parsedColumnMap[SQLITE3_COLUMNNAME], " ")
		parsedColumnMap[SQLITE3_CONSTRAINTS] = strings.Trim(parsedColumnMap[SQLITE3_CONSTRAINTS], " ")

		if ok, _ := sqlconstants.IsSQLConstraint(parsedColumnMap[SQLITE3_COLUMNNAME]); err == nil {

			if ok {
				return nil, errors.New("column with name = " + parsedColumnMap[SQLITE3_COLUMNNAME] + " is invalid")
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
