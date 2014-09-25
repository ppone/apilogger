package sqlconstants

import (
	"fmt"
	"testing"
)

func TestGoType(t *testing.T) {
	gotype, err := GoType(SQLITE3_TYPE_INTEGER)
	if err != nil {
		t.Error(err)
		return
	}
	if gotype != "int" {
		t.Error("Mapping of int is not proper")
		return
	}

	fmt.Println("CONSTANT VALUES => ", SQLITE3_TYPE_NIL, SQLITE3_TYPE_INTEGER)
}
