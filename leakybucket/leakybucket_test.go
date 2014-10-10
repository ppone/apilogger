package leakybucket

import "testing"
import "../db/sqlconstants"
import "time"

// import "fmt"

func TestNewBucket(t *testing.T) {
	s, err := NewBucket("A11", 10000, 3600, true)

	if err != nil {
		t.Error(err)
	}

	row, err := s.GoSelect()

	if row["NAME"].(string) != "A11" {
		t.Error("NAME not saved properly in database")
	}

	if s.Name() != "A11" {
		t.Error("NAME not saved properly in memory")
	}

	if row["CAPACITY"].(int) != 10000 {
		t.Error("CAPACITY not saved properly in database")
	}

	if s.Capacity() != 10000 {
		t.Error("CAPACITY not saved properly in memory")
	}

	if row["TIMEFRAME"].(int) != 3600 {
		t.Error("TIMEFRAME not saved properly in database")
	}

	if s.TimeFrame() != 3600 {
		t.Error("TIMEFRAME not saved properly in memory")
	}

	if row["VOLUME"].(int) != 0 {
		t.Error("VOLUME. not saved properly in database")
	}

	v, err := s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 0 {
		t.Error("VOLUME. not saved properly in memory")
	}

	if row["UPDATED_TIMESTAMP"].(int) <= 0 {
		t.Error("UPDATED_TIMESTAMP not saved properly in database")
	}

	if s.LastUpdated() <= 0 {
		t.Error("UPDATED_TIMESTAMP not saved properly in database")
	}

	conn, err := newConnection()

	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = conn.Drop(sqlconstants.SQLITE3_DROP_BUCKETS)

	if err != nil {
		t.Error(err)
	}

}

func TestRefresh(t *testing.T) {
	s, err := NewBucket("A11", 10000, 3600, true)

	if err != nil {
		t.Error(err)
	}

	lastUpdatedTimestamp1 := s.LastUpdated()

	const secondsToWait = 2

	time.Sleep(secondsToWait * time.Second)

	err = s.RefreshBucket()

	lastUpdatedTimestamp2 := s.LastUpdated()

	if lastUpdatedTimestamp2-lastUpdatedTimestamp1 != secondsToWait {
		t.Error("Refresh did not update lastUpdatedTimestamp")
	}

	if err != nil {
		t.Error(err)
	}

	err = s.RefreshBucket()
	if err != nil {
		t.Error(err)
	}
	err = s.RefreshBucket()
	if err != nil {
		t.Error(err)
	}
	err = s.RefreshBucket()
	if err != nil {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}
	err = s.RefreshBucket()
	if err != nil {
		t.Error(err)
	}
	err = s.RefreshBucket()
	if err != nil {
		t.Error(err)
	}

	conn, err := newConnection()

	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = conn.Drop(sqlconstants.SQLITE3_DROP_BUCKETS)

	if err != nil {
		t.Error(err)
	}

}

func TestFillVolume(t *testing.T) {

	//Threshold is 10000 API calls every hour or 60 minutes or 3600 seconds.

	//1800 API calls every hour or 1800 every 3600 seconds or 1 API every 2 seconds.

	//So Volume decreases every 2 seconds.

	s, err := NewBucket("A11", 1800, 3600, true)

	if err != nil {
		t.Error(err)
		return
	}

	v, err := s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 0 {
		t.Error("volume should be 0")
	}

	err = s.Increase()

	if err != nil {
		t.Error(err)
	}

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 1 {
		t.Error("volume should be 1")
	}

	err = s.Increase()

	if err != nil {
		t.Error(err)
	}

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 2 {
		t.Error("volume should be 2")
	}

	err = s.Increase()

	if err != nil {
		t.Error(err)
	}

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 3 {
		t.Error("volume should be 3")
	}

	const secondsToWait = 2

	time.Sleep(secondsToWait * time.Second)

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 2 {
		t.Error("volume should be 2")
	}

	time.Sleep(secondsToWait * time.Second)

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 1 {
		t.Error("volume should be 1")
	}

	time.Sleep(secondsToWait * time.Second)

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 0 {
		t.Error("volume should be 0")
	}

	time.Sleep(secondsToWait * time.Second)

	v, err = s.Volume()

	if err != nil {
		t.Error(err)
	}

	if v != 0 {
		t.Error("volume should be 0")
	}

	conn, err := newConnection()

	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = conn.Drop(sqlconstants.SQLITE3_DROP_BUCKETS)

	if err != nil {
		t.Error(err)
	}

}

func TestFullEmpty(t *testing.T) {

	s, err := NewBucket("A11", 3, 3600, true)

	if err != nil {
		t.Error(err)
	}

	empty, err := s.Empty()

	if err != nil {
		t.Error(err)
	}

	if !empty {
		t.Error("Error: Bucket should be empty")
	}

	err = s.Increase()

	if err != nil {
		t.Error(err)
	}

	err = s.Increase()

	if err != nil {
		t.Error(err)
	}

	err = s.Increase()

	if err != nil {
		t.Error(err)
	}

	full, err := s.Full()

	if !full {
		t.Error("Error: Bucket should be full")
	}

	conn, err := newConnection()

	defer conn.Close()

	if err != nil {
		t.Error(err)
	}

	err = conn.Drop(sqlconstants.SQLITE3_DROP_BUCKETS)

	if err != nil {
		t.Error(err)
	}

}
