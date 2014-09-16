package db

import "../sqlite"

const DB_TYPE_SQLITE3 = "sqlite3"
const SQLITE3_FILE_URL = "throttle.db"
const SQLITE3_CONNECTION_STRING = "file:" + SQLITE3_FILE_URL

func NewConnection(fileName string) (Connection, error) {

	if fileName == "" {

		fileName = file
	}

	db, err := sql.Open(DB_TYPE_SQLITE3, SQLITE3_CONNECTION_STRING)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &connection{db, fileName}, nil

}

type Connection interface {
	Insert(insertStatement string, data ...interface{}) error
	CreateTable(createStatement string, tableName string) error
	CheckTableExists(tableName string) (bool, error)
	Close() error
}
