package db

import _ "github.com/mattn/go-sqlite3"
import "../sqlite"
import "errors"

func NewConnection(dbDriver string) (Connection, error) {

	if dbDriver == "sqlite3" {
		return sqlite.NewConnection()
	}

	return nil, errors.New("DB Vendor: " + dbDriver + " not supported")
}

type Connection interface {
	Insert(insertStatement string, data ...interface{}) error
	InitTable(createStatement string, tableName string) error
	CheckTableExists(tableName string) (bool, error)
	Insert(insertStatement string, data ...interface{}) error
	Close() error
}
