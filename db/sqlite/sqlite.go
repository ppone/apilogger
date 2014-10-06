package sqlite

import "database/sql"

import "../sqlconstants"
import _ "github.com/mattn/go-sqlite3"
import "fmt"
import "errors"
import "../godata"
import "../../regexutil"

type connection struct {
	db *sql.DB
}

func NewConnection() (*connection, error) {

	db, err := sql.Open(sqlconstants.CurrentVendor(), sqlconstants.SQLITE3_CONNECTION_STRING)

	if err != nil {
		return nil, err
	}

	return &connection{db}, nil

}

func (Conn *connection) Close() error {
	return Conn.db.Close()
}

func (Conn *connection) CheckTableExists(tableName string) (bool, error) {
	rows, err := Conn.db.Query(sqlconstants.SQLITE3_CHECK_IF_TABLE_EXISTS, "table", tableName)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return rows.Next(), nil

}

func (Conn *connection) Select(query string, data ...interface{}) (*godata.GoSelect, error) {

	tableName, err := regexutil.TableNameFromSelect(query)

	if err != nil {
		return nil, err
	}

	if len(tableName) == 0 {
		return nil, errors.New("Could not figure out table name from select")
	}

	goMetaTable, err := Conn.MetaGoTable(tableName)

	if err != nil {
		return nil, err
	}

	goMetaTableColumns := goMetaTable.GetColumns()

	if goMetaTableColumns == nil {
		return nil, errors.New("Error Meta Table Columns are nil ")
	}

	resultRows, err := Conn.db.Query(query)

	defer resultRows.Close()

	if err != nil {
		return nil, err
	}

	columnNamesInResultRows, err := resultRows.Columns()

	if err != nil {
		return nil, err
	}

	selectColumns := godata.NewGoColumns()

	selectColumnCount := len(columnNamesInResultRows)

	for _, columnName := range columnNamesInResultRows {
		tableColumnType, ok := goMetaTableColumns[columnName]
		if !ok {
			return nil, errors.New("Column " + columnName + " , not found in hashmap")

		}

		selectColumns[columnName] = tableColumnType
	}

	selectRows := godata.NewGoRows()

	for resultRows.Next() {
		selectRow := godata.NewGoRow()
		results := []interface{}{}

		for _, columnName := range columnNamesInResultRows {
			tableColumn, ok := selectColumns[columnName]

			if !ok {

				return nil, errors.New("Column " + columnName + " , not found in hashmap")

			}

			results, err = godata.AppendData(results, tableColumn.GoType())

			if err != nil {
				return nil, err
			}

		}

		err := resultRows.Scan(results...)

		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < selectColumnCount; i++ {
			temp, err := godata.PointerConvertor(results[i])
			if err != nil {
				return nil, err
			}
			selectRow[columnNamesInResultRows[i]] = temp
		}

		selectRows = append(selectRows, selectRow)

	}

	return godata.NewGoSelect(selectColumns, selectRows), nil

	//sreturn nil, nil

}

func (Conn *connection) InitTable(createStatement string, tableName string) error {
	if createStatement == "" {
		return errors.New("Create Statement cannot be empty")

	}

	rowsExist, err := Conn.CheckTableExists(tableName)

	if err != nil {
		return err
	}

	if !rowsExist {

		_, err = Conn.db.Exec(createStatement)
		if err != nil {
			return err
		}

		rowsExist, err = Conn.CheckTableExists(tableName)

		if err != nil {
			return err
		}

		if rowsExist != true {
			return errors.New("Error => Could not create table " + tableName)
		}
	}

	return nil
}

func (Conn *connection) Insert(insertStatement string, data ...interface{}) (int64, error) {

	if insertStatement == "" {
		return -1, errors.New("Insert Statement cannot be blank")

	}

	if len(data) == 0 {
		return -1, errors.New("data cannot be empty")

	}

	tx, err := Conn.db.Begin()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	stmt, err := tx.Prepare(insertStatement)

	if err != nil {
		fmt.Println("Could not prepare statment = >", err)
		return -1, err
	}

	results, err := stmt.Exec(data...)

	defer stmt.Close()

	if err != nil {
		return -1, err
	}

	err = tx.Commit()

	n, err := results.LastInsertId()

	if err != nil {
		return -1, err
	}

	return n, err

}

func (Conn *connection) MetaGoTable(tableName string) (*godata.GoMetaTable, error) {

	//fetch create statement from db
	rows, err := Conn.db.Query(sqlconstants.SQLITE3_GET_SCHEMA, "table", tableName)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("ERROR: could not get schema for table " + tableName + " , in database")
	}
	var schema string

	rows.Scan(&schema)

	gmeta, err := godata.NewMetaTableFromCreateStatement(schema)

	if err != nil {
		return nil, err
	}
	return gmeta, nil

}
