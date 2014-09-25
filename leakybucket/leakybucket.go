package leakybucket

import "errors"
import "apilogger/db"
import "apilogger/db/sqlconstants"

type LeakyBucket interface {
	Capacity() int
	Fill(amount int) error
	Volume() int
	UpdateBucket()
	WhenWillBucketBeEmpty() (time, error)
	Empty() bools
}

type leakyBucket struct {
	name      string
	capacity  int
	volume    int
	timeFrame int //seconds
}

func newConnection() (db.Connection, error) {
	conn, err := db.New(sqlconstants.CURRENT_VENDOR)

	if err != nil {
		return nil, err
	}
}

func initDB() (db.Connection, error) {
	conn, err := newConnection()

	err = conn.InitTable(sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA, sqlconstants.SQLITE3_BUCKET_NAME)

	if err != nil {
		return err
	}

	defer conn.Close()

}

func LoadBucket(name string) (LeakeyBucket, error) {

}

func NewBucket(name string, capacity, timeFrame int) (LeakyBucket, error) {
	conn, err := newConnection()

	conn.Insert(sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA, name, capacity, 0, timeFrame)

	if err != nil {
		return nil, err
	}

	return &leakyBucket{name, capacity, volume, timeFrame}, nil
}

func (L *leakyBucket) Capacity() int {
	return L.capacity
}

func (L *leakyBucket) Empty() bool {
	L.RefreshBucket()
	return L.volume == 0
}

func (L *leakyBucket) Full() bool {
	L.RefreshBucket()
	return L.volume == L.capacity
}

func (L *leakyBucket) Volume() int {
	L.RefreshBucket()
	return L.volume
}

//seconds
//capacity / timeToRefresh ()  * timeNow - timeNow

func (L *leakyBucket) RefreshBucket() {

}

func (L *leakyBucket) Increase() error {
	err := L.Fill(1)
	if err != nil {
		return err
	}
	return nil
}

func (L *leakyBucket) Fill(amount int) error {
	L.UpdateBucket()
	if L.volume+amount > L.capacity {
		return errors.New("ERROR => Bucket is full; cannot add capacity")
	}

	L.volume += amount

	return nil

}
