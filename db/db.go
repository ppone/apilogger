package db

import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "../sqlite"

func NewConnection(fileName string) (Connection, error) {

	return sqlite.NewConnection(fileName)

}

type Connection interface {
	Insert(insertStatement string, data ...interface{}) error
	CreateTable(createStatement string, tableName string) error
	CheckTableExists(tableName string) (bool, error)
	Close() error
}
