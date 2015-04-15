package fixedbucket

/*
SELECT date('now','start of month','+1 month','-1 day');
SELECT DATE()

CREATE TABLE FIXEDBUCKETS(NAME TEXT NOT NULL UNIQUE PRIMARY KEY, VOLUME INTEGER DEFAULT 0, CAPACITY INTEGER, TIMEFRAME INTEGER, CREATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, UPDATED_TIMESTAMP INTEGER DEFAULT (strftime('%s','now')), START_DATE TEXT DEFAULT (date('now')), START_TIME TEXT, USE_START_DATE INTEGER);"

FIXEDBUCKETS STARTDATE, STARTTIME


NNN days
NNN hours
NNN minutes
NNN.NNNN seconds
NNN months
NNN years
start of month
start of year
start of day
weekday N
unixepoch
localtime
utc

SELECT STRTIME('%s','2014-07-04' || '03:40')


SELECT STRTIME('%s', START_DATE || START_TIME)


TIMEFRAME => SECONDS TILL VOLUME RESETS TO ZERO

UPDATING VOLUME FIRST SET VOLUME = 0, THEN SET VOLUME = NEW VALUE

FIRE TRIGGER EVERY X SECONDS FROM START_DATE OR UPDATED_DATE

FIRE TRIGGER EVERY X WEEK AT 10AM

FIRE TRIGGER EVERY X MONTH AT 9AM
 SELECT DATETIME(SELECT DATE('now') || '09:00',+1 month');

FIRE TRIGGER EVERY X YEAR AT 10AM
 SELECT DATE('now' || '10:00' ,'+12 month');

FIRE TRIGGER EVERY X DAYS AT
 SELECT DATE('now' || '10:00' ,'+X DAY');

*/

import "errors"
import "../db/sqlconstants"
import "../db"
import "../db/godata"

import "strconv"
import "time"
import "fmt"

//import "../timeutil"

type FixedBucket interface {
	Capacity() int
	Empty() (bool, error)
	Fill(amount int) error
	Full() (bool, error)
	GoSelect() (godata.GoRow, error)
	Increase() error
	LastUpdated() time.Time
	Name() string
	PrintRawSQL() (string, error)
	RefreshBucket() error
	Volume() (int, error)
	StartDateTime() time.Time
	LastUpdatedTimestamp() time.Time
}

type fixedBucket struct {
	capacity             int
	createdTimeStamp     time.Time
	lastUpdatedTimeStamp time.Time
	name                 string
	startDateTime        time.Time
	timeDuration         string
	timeFrequency        int
	volume               int
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

	err = conn.InitTable(sqlconstants.SQLITE3_CREATE_FIXED_BUCKETS, sqlconstants.SQLITE3_FIXEDBUCKET_NAME)

	if err != nil {
		return nil, err
	}

	return conn, err
}

func LoadBucket(name string) (FixedBucket, error) {
	conn, err := newConnection()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_FIXEDBUCKETS, name)

	if err != nil {
		return nil, err
	}

	rows := goSel.GetRows()

	if len(rows) == 0 {
		return nil, errors.New("Could not load bucket, as bucket does not seem to exist in the database")
	}

	capacity := ((rows[0])["CAPACITY"]).(int)
	volume := ((rows[0])["VOLUME"]).(int)
	timeFrequency := ((rows[0])["TIMEFRAME_FREQUENCY"]).(int)
	timeDuration := ((rows[0])["TIMEFRAME_DURATION"]).(string)
	startDateTime := ((rows[0])["START_DATETIME"]).(time.Time)
	createdTimeStamp := ((rows[0])["CREATED_TIMESTAMP"]).(time.Time)
	lastUpdatedTimeStamp := ((rows[0])["UPDATED_TIMESTAMP"]).(time.Time)

	// 	type fixedBucket struct {
	// 	capacity             int
	// 	createdTimeStamp     time.Time
	// 	lastUpdatedTimeStamp time.Time
	// 	name                 string
	// 	startDateTime        time.Time
	// 	timeDuration         string
	// 	timeFrequency        int
	// 	volume               int
	// }

	return &fixedBucket{capacity, createdTimeStamp, lastUpdatedTimeStamp, name, startDateTime, timeDuration, timeFrequency, volume}, nil

}

func NewBucket(name string, capacity, timeFrequency int, timeDuration string, startDateTime time.Time, loadBucketIfExists bool) (FixedBucket, error) {
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

	fmt.Println("UP IN HERE")
	_, err = conn.Insert(sqlconstants.SQLITE3_INSERT_INTO_FIXEDBUCKET, name, 0, capacity, timeFrequency, timeDuration, startDateTime)

	if err != nil {
		return nil, err
	}

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_FIXEDBUCKETS, name)
	if err != nil {
		return nil, err
	}

	rows := goSel.GetRows()

	// NAME TEXT
	// VOLUME INTEGER
	// CAPACITY INTEGER
	// TIMEFRAME_FREQUENCY INTEGER
	// TIMEFRAME_DURATION TEXT
	// START_DATETIME TIMESTAMP
	// CREATED_TIMESTAMP TIMESTAMP
	// UPDATED_TIMESTAMP TIMESTAMP

	// 	type fixedBucket struct {
	// 	name          string
	// 	capacity      int
	// 	volume        int
	// 	timeFrequency int
	// 	timeDuration  string
	// 	lastUpdated   time.Time
	// 	startDateTime time.Time
	// }

	createdTimeStamp := ((rows[0])["CREATED_TIMESTAMP"]).(time.Time)
	lastUpdatedTimestamp := ((rows[0])["UPDATED_TIMESTAMP"]).(time.Time)

	return &fixedBucket{capacity, createdTimeStamp, lastUpdatedTimestamp, name, startDateTime, timeDuration, timeFrequency, 0}, nil
}

func (L *fixedBucket) Capacity() int {

	return L.capacity
}

func (L *fixedBucket) Name() string {

	return L.name
}

func (L *fixedBucket) TimeFrequency() int {

	return L.timeFrequency
}

func (L *fixedBucket) TimeDuration() string {

	return L.timeDuration
}

func (L *fixedBucket) StartDateTime() time.Time {

	return L.startDateTime
}

func (L *fixedBucket) LastUpdatedTimestamp() time.Time {

	return L.lastUpdatedTimeStamp
}

func (L *fixedBucket) CreatedTimestamp() time.Time {

	return L.createdTimeStamp
}

func (L *fixedBucket) Empty() (bool, error) {
	err := L.RefreshBucket()
	if err != nil {
		return false, err
	}
	return L.volume == 0, nil
}

func (L *fixedBucket) Full() (bool, error) {
	err := L.RefreshBucket()
	if err != nil {
		return false, err
	}
	return L.volume == L.capacity, nil
}

func (L *fixedBucket) Volume() (int, error) {

	err := L.RefreshBucket()
	if err != nil {
		return -1, err
	}
	return L.volume, nil
}

//
//
func (L *fixedBucket) RefreshBucket() error {
	conn, err := newConnection()
	var timeNOW time.Time

	defer conn.Close()

	if err != nil {
		return err
	}

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_FIXEDBUCKETS_W_TIMENOW, L.name)

	if err != nil {
		return err
	}

	rows := goSel.GetRows()
	timeNOW = ((rows[0])["TIMENOW"]).(time.Time)
	L.startDateTime = ((rows[0])["START_DATETIME"]).(time.Time)
	L.lastUpdatedTimeStamp = ((rows[0])["UPDATED_TIMESTAMP"]).(time.Time)

	fmt.Println(timeNOW.Day())

	// err = conn.Update(sqlconstants.SQLITE3_REFRESH_BUCKET, L.name, L.name)
	// if err != nil {
	// 	return err
	// }

	return nil

}

func (L *fixedBucket) LastUpdated() time.Time {
	return L.lastUpdatedTimeStamp
}

func (L *fixedBucket) Increase() error {

	err := L.Fill(1)
	if err != nil {
		return err
	}
	return nil
}

func (L *fixedBucket) Fill(amount int) error {
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

func (L *fixedBucket) GoSelect() (godata.GoRow, error) {
	conn, err := newConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	goSel, err := conn.Select(sqlconstants.SQLITE3_SELECT_ALL_FIXEDBUCKETS, L.name)

	if err != nil {
		return nil, err
	}

	rows := goSel.GetRows()

	if len(rows) == 0 {
		return nil, errors.New("Database Error := row not found")
	}

	return goSel.GetRows()[0], nil
}

func (L *fixedBucket) PrintRawSQL() (string, error) {

	row, err := L.GoSelect()

	if err != nil {
		return "", err
	}

	// NAME TEXT
	// VOLUME INTEGER
	// CAPACITY INTEGER
	// TIMEFRAME_FREQUENCY INTEGER
	// TIMEFRAME_DURATION TEXT
	// START_DATETIME TIMESTAMP
	// CREATED_TIMESTAMP TIMESTAMP
	// UPDATED_TIMESTAMP TIMESTAMP

	name := (row["NAME"]).(string)
	capacity := strconv.Itoa((row["CAPACITY"]).(int))
	volume := strconv.Itoa((row["VOLUME"]).(int))
	timeFrequency := strconv.Itoa((row["TIMEFRAME_FREQUENCY"]).(int))
	timeDuration := (row["TIMEFRAME_DURATION"]).(string)

	createdTimestamp := (row["CREATED_TIMESTAMP"]).(time.Time).String()
	lastUpdated := (row["UPDATED_TIMESTAMP"]).(time.Time).String()
	startDateTime := (row["START_DATETIME"]).(time.Time).String()

	s := "NAME => " + name + ", VOLUME => " + volume + ", capacity => " + capacity + ", timeFrequency => " + timeFrequency + ", timeDuration => " + timeDuration + ", lastUpdated => " + lastUpdated + ", createdTimestamp => " + createdTimestamp + ", startDateTime " + startDateTime

	return s, nil

}
