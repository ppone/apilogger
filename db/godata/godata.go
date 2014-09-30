package godata

import "../../regexutil"
import "../sqlconstants"

type GoData struct {
	goValue interface{}
}

type GoColumn struct {
	name    string
	sqlType string
	goType  string
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

func (c *GoColumns) Add(name string, sqlType, goType string) *GoColumns {
	data := GoColumn{name, sqlType, goType}
	c.columns = append(c.columns, data)
	return c
}

func (r *GoRow) Add(value interface{}) *GoRow {
	data := GoData{value}
	r.row = append(r.row, data)
	return r
}

func NewMetaTableFromCreateStatement(stmt string) (*GoMetaTable, error) {

	table, columnArrayMap, err := regexutil.ParseCreateStatement(stmt)
	cols := NewGoColumns()

	if err != nil {
		return nil, err
	}

	for _, column := range columnArrayMap {
		columnname := column[regexutil.COLUMNNAME]
		sqltype := column[regexutil.COLUMNTYPE]
		gotype, err := sqlconstants.GoType(sqltype)
		if err != nil {
			return nil, err
		}
		cols.Add(columnname, sqltype, gotype)

	}

	return &GoMetaTable{table, sqlconstants.CurrentVendor(), cols}, nil

}
