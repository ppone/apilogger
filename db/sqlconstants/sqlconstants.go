package sqlconstants

import "errors"
import "strings"
import "strconv"
import "time"
import "../../timeutil"

const CURRENT_VENDOR = SQLITE3
const SQLITE3_FILE_URL = "throttle.db"
const SQLITE3_CONNECTION_STRING = "file:" + SQLITE3_FILE_URL
const SQLITE3_CHECK_IF_TABLE_EXISTS = "SELECT NAME FROM SQLITE_MASTER WHERE TYPE=? AND NAME=?;"
const SQLITE3_CREATE_BUCKET_SCHEMA = "CREATE TABLE BUCKETS(NAME TEXT NOT NULL UNIQUE PRIMARY KEY, VOLUME INTEGER DEFAULT 0, CAPACITY INTEGER, TIMEFRAME INTEGER, CREATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, UPDATED_TIMESTAMP INTEGER DEFAULT (strftime('%s','now')))"
const SQLITE3_INSERT_INTO_BUCKET = "INSERT INTO BUCKETS (NAME, VOLUME, CAPACITY, TIMEFRAME) VALUES (?, ?, ?, ?)"
const SQLITE3_UPDATE_BUCKET = "UPDATE BUCKETS SET (VOLUME,CAPACITY,TIMEFRAME, UPDATED_TIMESTAMP) VALUES (?, ?, ?, ?)"

const SQLITE3_REFRESH_BUCKET = "UPDATE BUCKETS SET VOLUME = (SELECT CASE WHEN expr > 0 THEN expr ELSE 0 END FROM (SELECT CAST(ROUND(VOLUME - ((CAPACITY*1.0)/TIMEFRAME)*(STRFTIME('%s','now') - UPDATED_TIMESTAMP)) AS INT) AS expr FROM BUCKETS WHERE NAME = ?) b), UPDATED_TIMESTAMP = STRFTIME('%s', 'now') WHERE NAME = ?;"
const SQLITE3_LOAD_BUCKET = "SELECT VOLUME, CAPACITY, TIMEFRAME FROM BUCKETS WHERE NAME = ?"
const SQLITE3_DELETE_BUCKET = "DELETE BUCKETS WHERE NAME = ?"
const SQLITE3_SELECT_ALL_BUCKETS = "SELECT * FROM BUCKETS where NAME = ?"
const SQLITE3_BUCKET_NAME = "BUCKETS"
const SQLITE3_GET_SCHEMA = "SELECT SQL FROM SQLITE_MASTER WHERE TYPE=? AND NAME=?"
const SQLITE3_FILL_BUCKET = "UPDATE BUCKETS SET VOLUME = VOLUME + ? WHERE NAME = ?"
const SQLITE3_DROP_BUCKETS = "DROP TABLE BUCKETS"

const SQLITE3_CREATE_FIXED_BUCKETS = "CREATE TABLE FIXEDBUCKETS(NAME TEXT NOT NULL UNIQUE PRIMARY KEY, VOLUME INTEGER DEFAULT 0, CAPACITY INTEGER, TIMEFRAME_FREQUENCY INTEGER NOT NULL, TIMEFRAME_DURATION TEXT NOT NULL, START_DATETIME TIMESTAMP, CREATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, UPDATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP);"

const SQLITE3_FIXEDBUCKET_NAME = "FIXEDBUCKETS"
const SQLITE3_SELECT_ALL_FIXEDBUCKETS = "SELECT * FROM FIXEDBUCKETS where NAME = ?"
const SQLITE3_SELECT_ALL_FIXEDBUCKETS_W_TIMENOW = "SELECT datetime('now') as TIMENOW, * FROM FIXEDBUCKETS where NAME = ?"

//_, err = conn.Insert(sqlconstants.SQLITE3_INSERT_INTO_BUCKET, name, 0, capacity, timeFrequency, timeDuration, startDateTime)

const SQLITE3_INSERT_INTO_FIXEDBUCKET = "INSERT INTO FIXEDBUCKETS (NAME, VOLUME, CAPACITY, TIMEFRAME_FREQUENCY, TIMEFRAME_DURATION, START_DATETIME) VALUES (?, ?, ?, ?, ?, ?)"

const SQLITE3_REFRESH_FIXEDBUCKET = "UPDATE FIXBUCKETS SET VOLUME = (SELECT )"

const SQLITE3_DROP_FIXEDBUCKETS = "DROP TABLE FIXEDBUCKETS"

const (
	SQLITE3_TYPE_NULL      = ""
	SQLITE3_TYPE_INTEGER   = "INTEGER"
	SQLITE3_TYPE_TEXT      = "TEXT"
	SQLITE3_TYPE_REAL      = "REAL"
	SQLITE3_TYPE_BLOB      = "BLOG"
	SQLITE3_TYPE_BOOL      = "BOOL"
	SQLITE3_TYPE_TIMESTAMP = "TIMESTAMP"
	//add postgres, plus other types here
)

const (
	SQLITE3 = iota
	POSTGRESQL
	MYSQL
	MONGODB
)

func FireTriggerEveryXDaysFrom(triggerDay int, startTimeDate time.Time) (string, error) {

	if CurrentVendor() != "sqlite3" {
		return "", errors.New("Vendor not recognized ")
	}
	baseSQL := "DATETIME('" + timeutil.NormalFormat(startTimeDate) + "',"
	baseSQL = baseSQL + "'+" + strconv.Itoa(triggerDay) + " days')"

	return baseSQL, nil
}

func FireTriggerEveryXMonthsFrom(triggerMonth int, startTimeDate time.Time) (string, error) {

	if CurrentVendor() != "sqlite3" {
		return "", errors.New("Vendor not recognized ")
	}
	baseSQL := "DATETIME('" + timeutil.NormalFormat(startTimeDate) + "',"
	baseSQL = baseSQL + "'+" + strconv.Itoa(triggerMonth) + " months')"

	return baseSQL, nil
}

func FireTriggerEveryXWeeksFrom(triggerWeek int, startTimeDate time.Time) (string, error) {

	if CurrentVendor() != "sqlite3" {
		return "", errors.New("Vendor not recognized ")
	}
	triggerWeek = triggerWeek * 7

	baseSQL := "DATETIME('" + timeutil.NormalFormat(startTimeDate) + "',"
	baseSQL = baseSQL + "'+" + strconv.Itoa(triggerWeek) + " days')"

	return baseSQL, nil
}

func CreateStatementFunctionsToReplace() ([]string, error) {
	if CURRENT_VENDOR == SQLITE3 {
		return []string{"(strftime('%s','now'))"}, nil
	}

	return nil, errors.New("Error: Current Vendor is not support in CreateStatementFunctionsToReplace")
}

func CurrentVendor() string {
	if CURRENT_VENDOR == SQLITE3 {
		return "sqlite3"
	} else if CURRENT_VENDOR == POSTGRESQL {
		return "postgres"
	} else if CURRENT_VENDOR == MYSQL {
		return "mysql"
	} else if CURRENT_VENDOR == MONGODB {
		return "mongodb"
	}

	return ""
}

func IsSQLConstraint(constraint string) (bool, error) {
	c := strings.ToUpper(constraint)
	if CURRENT_VENDOR == SQLITE3 {
		switch c {
		case "PRIMARY":
			return true, nil
		case "DEFAULT":
			return true, nil
		case "NOT":
			return true, nil
		case "UNIQUE":
			return true, nil
		case "CHECK":
			return true, nil
		case "REFERENCES":
			return true, nil
		case "COLLATE":
			return true, nil
		default:
			return false, nil
		}

	}

	return false, errors.New("Go type no recongnized for current db vendor => " + CurrentVendor())

}

func GoType(sqlType string) (string, error) {
	if CURRENT_VENDOR == SQLITE3 {
		switch sqlType {
		case SQLITE3_TYPE_NULL:
			return "nil", nil
		case SQLITE3_TYPE_INTEGER:
			return "int", nil
		case SQLITE3_TYPE_REAL:
			return "float64", nil
		case SQLITE3_TYPE_BOOL:
			return "bool", nil
		case SQLITE3_TYPE_TEXT:
			return "string", nil
		case SQLITE3_TYPE_TIMESTAMP:
			return "time.Time", nil
		default:
			return "", nil

		}
	}

	return "", errors.New("Go type no recongnized for current db vendor => " + CurrentVendor())

}

func SQLType(goType string) (string, error) {
	if CURRENT_VENDOR == SQLITE3 {
		switch goType {
		case "nil":
			return SQLITE3_TYPE_NULL, nil
		case "int":
			return SQLITE3_TYPE_INTEGER, nil
		case "float64":
			return SQLITE3_TYPE_REAL, nil
		case "bool":
			return SQLITE3_TYPE_BOOL, nil
		case "string":
			return SQLITE3_TYPE_TEXT, nil
		case "time.Time":
			return SQLITE3_TYPE_TIMESTAMP, nil

		}
	}

	return "", errors.New("Go type no recongnized for current db vendor => " + CurrentVendor())

}
