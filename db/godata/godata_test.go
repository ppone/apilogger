package godata

import (
	"fmt"
	"testing"
)

const TEST_SQLITE3_CREATE_STATEMENT_A = "CREATE TABLE FOO ( X INTEGER, Y TEXT PRIMARY KEY)"

func TestNewMetaTableFromCreateStatement(t *testing.T) {
	meta, err := NewMetaTableFromCreateStatement(TEST_SQLITE3_CREATE_STATEMENT_A)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(meta.columns)
}
