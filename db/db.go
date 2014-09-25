package db

import _ "github.com/mattn/go-sqlite3"
import "../sqlite"
import "errors"
import "strings"

type GoData struct {
	goValue interface{}
}

type GoColumn struct {
	name    string
	sqlType int
	goType  int
}

type GoRow struct {
	row []GoData
}

type GoRows struct {
	rows []GoRow
}
type GoColumns struct {
	columns []GoColumn
}

type GoMetaTable struct {
	name     string
	dbVendor string
	columns  *GoColumns
}

type GoTable struct {
	metatable *GoMetaTable
	rows      *GoRow
}

func NewGoRow() *GoRow {
	return &GoRow{}
}

func NewGoColumns() *GoColumns {
	return &GoColumns{}
}

func NewGoColumnsFromCreateStatement(stmt string) error {
	if !strings.Contains(stmt, "create table") || !strings.con

}

func (c *GoColumns) Add(name string, sqlType, goType int) *GoColumns {
	data := GoColumn{name, sqlType, goType}
	c.columns = append(c.columns, data)
	return c
}

func (r *GoRow) Add(value interface{}) *GoRow {
	data := GoData{value}
	r.row = append(r.row, data)
	return r
}

func NewGoTable(name, dbVendor string, columns *GoColumns, rows *GoRow) *GoTable {
	return &GoTable{name, dbVendor, columns, rows}

}

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
