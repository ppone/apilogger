package fixedbucket

import "testing"
import "time"
import "fmt"
import "../db/sqlconstants"
import "reflect"

func TestNewBucket(t *testing.T) {

	timenow := time.Now()
	time.Sleep(120 * time.Second)
	ns, err := NewBucket("YELP2", 1000, 1, "DAYS", timenow, true)

	if err != nil {
		t.Fatal(err)
	}

	s, err := ns.PrintRawSQL()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(s)

	conn, err := newConnection()

	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = conn.Drop(sqlconstants.SQLITE3_DROP_FIXEDBUCKETS)

	if err != nil {
		t.Error(err)
	}

}

func TestLoadBucket(t *testing.T) {

	timenow := time.Now()
	ns, err := NewBucket("YELP", 1000, 1, "DAYS", timenow, true)

	if err != nil {
		t.Fatal(err)
	}

	vs, err := LoadBucket("YELP")

	if err != nil {
		t.Fatal(err)
	}

	s, err := ns.PrintRawSQL()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(s)
	fmt.Println()

	if reflect.DeepEqual(ns, vs) {
		t.Error("Values are not equal")
	}

	conn, err := newConnection()

	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = conn.Drop(sqlconstants.SQLITE3_DROP_FIXEDBUCKETS)

	if err != nil {
		t.Error(err)
	}

}
