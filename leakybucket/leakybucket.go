package leakybucket

import "errors"
import "../db/sqlconstants"
import "../db"
import "../db/godata"

import "strconv"
import "time"

type LeakyBucket interface {
	Capacity() int
	Name() string
	TimeFrame() int
	Empty() (bool, error)
	Full() (bool, error)
	Increase() error
	Fill(amount int) error
	LastUpdated() int
	GoSelect() (godata.GoRow, error)
	PrintRawSQL() (string, error)
	RefreshBucket() error
	Volume() (int, error)
}

type leakyBucket struct {
	name        string
	capacity    int
	volume      int
	timeFrame   int //unix epoch seconds
	lastUpdated int //unix epoch seconds
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
	lastUpdated := ((rows[0])["UPDATED_TIMESTAMP"]).(int)

	return &leakyBucket{name, capacity, volume, timeFrame, lastUpdated}, nil

}

func NewBucket(name string, capacity, timeFrame int, loadBucketIfExists bool) (LeakyBucket, error) {
	conn, err := initDB()

	defer conn.Close()

	if loadBucketIfExists {
		l, err := LoadBucket(name)

		if err != nil {
			if err.Error() != "Could not load bucket, as bucket does not seem to exist in the database" {

				return nil, err

			}

		}

		if err == nil {
			return l, nil
		}

	}

	_, err = conn.Insert(sqlconstants.SQLITE3_INSERT_INTO_BUCKET, name, 0, capacity, timeFrame)

	if err != nil {
		return nil, err
	}

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_BUCKETS, name)
	if err != nil {
		return nil, err
	}

	rows := goSel.GetRows()

	lastUpdated := ((rows[0])["UPDATED_TIMESTAMP"]).(int)

	return &leakyBucket{name, capacity, 0, timeFrame, lastUpdated}, nil
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
	defer conn.Close()
	if err != nil {
		return err
	}

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
	L.lastUpdated = ((rows[0])["UPDATED_TIMESTAMP"]).(int)

	return nil

}

func (L *leakyBucket) LastUpdated() int {
	return L.lastUpdated
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
	defer conn.Close()

	err = conn.Update(sqlconstants.SQLITE3_FILL_BUCKET, amount, L.name)
	if err != nil {
		return err
	}

	L.volume += amount

	return nil

}

func (L *leakyBucket) GoSelect() (godata.GoRow, error) {
	conn, err := newConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_BUCKETS, L.name)

	if err != nil {
		return nil, err
	}

	rows := goSel.GetRows()

	if len(rows) == 0 {
		return nil, errors.New("Database Error := row not found")
	}

	return goSel.GetRows()[0], nil
}

func (L *leakyBucket) PrintRawSQL() (string, error) {

	row, err := L.GoSelect()

	if err != nil {
		return "", err
	}

	name := (row["NAME"]).(string)
	volume := strconv.Itoa((row["VOLUME"]).(int))
	timeFrame := strconv.Itoa((row["TIMEFRAME"]).(int))
	lastUpdated := strconv.Itoa((row["UPDATED_TIMESTAMP"]).(int))
	capacity := strconv.Itoa((row["CAPACITY"]).(int))
	createdTimestamp := (row["CREATED_TIMESTAMP"]).(time.Time).String()

	s := "NAME => " + name + ", VOLUME => " + volume + ", timeFrame => " + timeFrame + ", lastUpdated => " + lastUpdated + ", capacity => " + capacity + ", createdTimestamp => " + createdTimestamp

	return s, nil

}
