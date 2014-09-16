package leakybucket

import "time"
import "errors"
import "../db"

const createBucketSchema := "CREATE TABLE BUCKETS(NAME TEXT NOT NULL UNIQUE PRIMARY KEY, VOLUME INTEGER, CAPACITY INTEGER, TIMEFRAME INTEGER, CREATED_TIMESTAMP DEFAULT CURRENT_TIMESTAMP, UPDATED_TIMESTAMP TEXT)"

type LeakyBucket interface {
	Capacity() int
	Fill(amount int) error
	Volume() int
	UpdateBucket()
	WhenWillBucketBeEmpty() (time, error)
	Empty() bool
}

type leakyBucket struct {
	name      string
	capacity  int
	volume    int
	timeFrame int //seconds 
}

func NewBucket(name string, capacity, timeFrame int) LeakyBucket {
	return &leakyBucket{capcity, 0, timeFrame, capacity / timeFrame}
}

func (L *leakyBucket) Capacity() int {
	L.UpdateBucket()
	return L.capacity

}

func (L *leakyBucket) Empty() bool {
	return L.volume == 0
}

func (L *leakyBucket) Volume() int {
	return L.volume == 0
}


func (L *leakyBucket) UpdateBucket() {

}

func (L *leakyBucket) Fill(amount int) error {
	L.UpdateBucket()
	if L.volume+amount > L.capacity {
		return errors.New("ERROR => Bucket is full; cannot add capacity")
	}

	L.volume += amount

	return nil

}
