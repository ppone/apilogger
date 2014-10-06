package godata

import (
	"errors"
	"time"

	"../../regexutil"
)
import "../sqlconstants"

type GoType struct {
	sqlType string
	goType  string
}

type GoRow map[string]interface{}

type GoRows []GoRow

type GoColumns map[string]GoType

type GoSelect struct {
	columns GoColumns
	rows    GoRows
}

type GoMetaTable struct {
	name     string
	dbVendor string
	columns  GoColumns
}

type GoTable struct {
	metatable GoMetaTable
	rows      GoRow
}

func NewGoRow() GoRow {
	return GoRow{}
}

func NewGoRows() GoRows {
	return GoRows{}
}

func NewGoColumns() GoColumns {
	return GoColumns{}
}

func NewGoSelect(columns GoColumns, rows GoRows) *GoSelect {
	return &GoSelect{columns, rows}
}

func (s GoSelect) GetColumns() GoColumns {
	return s.columns
}

func (s GoSelect) GetRows() []GoRow {
	return s.rows
}

func (c GoColumns) Add(name, sqlType, goType string) GoColumns {
	t := GoType{sqlType, goType}
	c[name] = t
	return c
}

func (t GoType) GoType() string {
	return t.goType
}

func (t GoType) SqlType() string {
	return t.sqlType
}

func (m GoMetaTable) GetColumns() GoColumns {
	return m.columns
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

func AppendData(dataArray []interface{}, dataType string) ([]interface{}, error) {
	t := new(string)
	switch dataType {
	case "string":
		r := append(dataArray, t)
		return r, nil
	case "int":
		return append(dataArray, new(int)), nil
	case "float64":
		return append(dataArray, new(float64)), nil
	case "bool":
		return append(dataArray, new(bool)), nil
	case "nil":
		return append(dataArray, nil), nil
	case "time.Time":
		return append(dataArray, new(time.Time)), nil
	}

	return nil, errors.New("type is not recognized")

}

func PointerConvertor(data interface{}) (interface{}, error) {
	switch data.(type) {
	case *string:
		return *(data.(*string)), nil
	case *int:
		return *(data.(*int)), nil
	case *float64:
		return *(data.(*float64)), nil
	case *bool:
		return *(data.(*bool)), nil
	case nil:
		return nil, nil
	case *time.Time:
		return *(data.(*time.Time)), nil
	}

	return nil, errors.New("data type is not recognized")

}
