package db

import _ "github.com/mattn/go-sqlite3"

import "./sqlite"
import "errors"
import "./godata"

func NewConnection(dbDriver string) (Connection, error) {

	if dbDriver == "sqlite3" {
		return sqlite.NewConnection()
	}

	return nil, errors.New("DB Vendor: " + dbDriver + " not supported")
}

type Connection interface {
	CheckTableExists(tableName string) (bool, error)
	Close() error
	InitTable(createStatement string, tableName string) error
	Insert(insertStatement string, data ...interface{}) (int64, error)
	MetaGoTable(name string) (*godata.GoMetaTable, error)
	Select(query string, data ...interface{}) (*godata.GoSelect, error)
}
