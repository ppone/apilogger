package leakybucket

import "testing"
import "../db/sqlconstants"
import "fmt"

func TestNewBucket(t *testing.T) {
	s, err := NewBucket("A11", 10000, 3600)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(s.Name())

	conn, err := newConnection()
	if err != nil {
		t.Error(err)
	}

	goData, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_BUCKETS, s.Name())

	if err != nil {
		t.Error(err)
	}

	fmt.Println(goData.GetRows()[0]["CREATED_TIMESTAMP"])

	err = conn.Drop(sqlconstants.SQLITE3_DROP_BUCKETS)

	if err != nil {
		t.Error(err)
	}

	err = conn.Close()

	if err != nil {
		t.Error(err)
	}

}
