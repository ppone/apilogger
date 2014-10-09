package leakybucket

import "errors"
import "../db/sqlconstants"
import "../db"

type LeakyBucket interface {
	Capacity() int
	Name() string
	TimeFrame() int
	Empty() (bool, error)
	Full() (bool, error)
	Increase() error
	Fill(amount int) error
}

type leakyBucket struct {
	name      string
	capacity  int
	volume    int
	timeFrame int //unix epoch seconds
}

func newConnection() (db.Connection, error) {
	conn, err := db.NewConnection(sqlconstants.CurrentVendor())

	return conn, err
}

func initDB() (db.Connection, error) {
	conn, err := newConnection()

	if err != nil {
		return nil, err
	}

	err = conn.InitTable(sqlconstants.SQLITE3_CREATE_BUCKET_SCHEMA, sqlconstants.SQLITE3_BUCKET_NAME)

	if err != nil {
		return nil, err
	}

	return conn, err
}

func LoadBucket(name string) (LeakyBucket, error) {
	conn, err := newConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_BUCKETS, name)
	if err != nil {
		return nil, err
	}

	rows := goSel.GetRows()

	if len(rows) == 0 {
		return nil, errors.New("Could not load bucket, as bucket does not seem to exist in the database")
	}

	capacity := ((rows[0])["CAPACITY"]).(int)
	volume := ((rows[0])["VOLUME"]).(int)
	timeFrame := ((rows[0])["TIMEFRAME"]).(int)

	return &leakyBucket{name, capacity, volume, timeFrame}, nil

}

func NewBucket(name string, capacity, timeFrame int) (LeakyBucket, error) {
	conn, err := initDB()

	defer conn.Close()

	conn.Insert(sqlconstants.SQLITE3_INSERT_INTO_BUCKET, name, capacity, 0, timeFrame)

	if err != nil {
		return nil, err
	}

	return &leakyBucket{name, capacity, 0, timeFrame}, nil
}

func (L *leakyBucket) Capacity() int {

	return L.capacity
}

func (L *leakyBucket) Name() string {

	return L.name
}

func (L *leakyBucket) TimeFrame() int {

	return L.timeFrame
}

func (L *leakyBucket) Empty() (bool, error) {
	err := L.RefreshBucket()
	if err != nil {
		return false, err
	}
	return L.volume == 0, nil
}

func (L *leakyBucket) Full() (bool, error) {
	err := L.RefreshBucket()
	if err != nil {
		return false, err
	}
	return L.volume == L.capacity, nil
}

func (L *leakyBucket) Volume() (int, error) {

	err := L.RefreshBucket()
	if err != nil {
		return -1, err
	}
	return L.volume, nil
}

//seconds
//volume = (volume - capacity / timeToRefresh) * (timeNow - timeNow)
func (L *leakyBucket) RefreshBucket() error {
	conn, err := newConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.Update(sqlconstants.SQLITE3_REFRESH_BUCKET, L.name)
	if err != nil {
		return err
	}

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_BUCKETS, L.name)
	if err != nil {
		return err
	}

	rows := goSel.GetRows()
	L.capacity = ((rows[0])["CAPACITY"]).(int)
	L.volume = ((rows[0])["VOLUME"]).(int)
	L.timeFrame = ((rows[0])["TIMEFRAME"]).(int)

	return nil

}

func (L *leakyBucket) Increase() error {

	err := L.Fill(1)
	if err != nil {
		return err
	}
	return nil
}

func (L *leakyBucket) Fill(amount int) error {
	err := L.RefreshBucket()
	if err != nil {
		return err
	}

	if L.volume+amount > L.capacity {
		return errors.New("ERROR => Bucket is full; cannot add capacity")
	}

	conn, err := newConnection()
	if err != nil {
		return nil
	}

	err = conn.Update(sqlconstants.SQLITE3_FILL_BUCKET, L.name)
	if err != nil {
		return err
	}

	L.volume += amount

	return nil

}
