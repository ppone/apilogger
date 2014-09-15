package sqlite

//import "fmt"
import "testing"

func Test_Open_DB(t *testing.T) {
	conn, err := NewConnection("")
	err = conn.createTable(create_table_throttle, "throttle")
	if err != nil {
		t.Error(err)
	}
}
