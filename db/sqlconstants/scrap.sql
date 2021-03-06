UPDATE BUCKETS
SET VOLUME =
  (SELECT CASE
              WHEN expr > 0 THEN expr
              ELSE 0
          END
   FROM
     (SELECT CAST(ROUND(VOLUME - ((CAPACITY*1.0)/TIMEFRAME)*(STRFTIME('%s','now') - UPDATED_TIMESTAMP)) AS INT) AS expr
      FROM BUCKETS
      WHERE NAME = ?) b),
    UPDATED_TIMESTAMP = STRFTIME('%s', 'now')
WHERE NAME = ?;


SELECT CASE
           WHEN expr > 0 THEN expr
           ELSE 0
       END
FROM
  ( SELECT CAST(ROUND(VOLUME - ((CAPACITY*1.0)/TIMEFRAME)*(STRFTIME('%s','now') - UPDATED_TIMESTAMP)) AS INT) AS expr
   FROM BUCKETS
   WHERE NAME = ?) b), UPDATED_TIMESTAMP = STRFTIME('%s', 'now')
WHERE NAME = ?;

 const SQLITE3_CREATE_FIXED_BUCKETS =
CREATE TABLE FIXEDBUCKETS (NAME TEXT NOT NULL UNIQUE PRIMARY KEY, VOLUME INTEGER DEFAULT 0, CAPACITY INTEGER, TIMEFRAME_FREQUENCY INTEGER NOT NULL, TIMEFRAME_DURATION TEXT NOT NULL, START_DATETIME TIMESTAMP DEFAULT CURRENT_TIMESTAMP, CREATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, UPDATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP);

 ;

IF START_DATETIME_SET THEN START_DATETIME IS NOT updated ;

CREATE TABLE 
FIXEDBUCKETS (
NAME TEXT NOT NULL UNIQUE PRIMARY KEY, 
VOLUME INTEGER DEFAULT 0, 
CAPACITY INTEGER, 
TIMEFRAME_FREQUENCY INTEGER NOT NULL, 
TIMEFRAME_DURATION TEXT NOT NULL, 
START_DATETIME TIMESTAMP DEFAULT CURRENT_TIMESTAMP,CREATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, UPDATED_TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP);


type fixedBucket struct {
	capacity             int
	createdTimeStamp     time.Time
	lastUpdatedTimeStamp time.Time
	name                 string
	startDateTime        time.Time
	timeDuration         string
	timeFrequency        int
	volume               int
};



DATETIME(START_DATETIME, '+3 DAYS') but INSTEAD IS
SELECT CASE WHEN CT > 0 THEN EXPR
SELECT COUNT(*) AS CT
FROM FIXEDBUCKETS
WHERE DATETIME('now') >= DATETIME(START_DATETIME,'+3 days')
  AND NAME = ?
  AND START_DATETIME
  UPDATE START_DATETIME
  SET START_DATETIME = DATETIME(START_DATETIME, '+3 DAYS')
  UPDATE START_DATETIME
  SET


  2 QUERIES 

  SELECT CURRENT_TIMESTAMP, START_DATETIME, TIMEFRAME_FREQUENCY,TIMEFRAME_DURATION, VOLUME FROM FIXEDBUCKETS WHERE NAME = ?

  UPDATE FIXEDBUCKETS SET START_DATETIME = ? AND VOLUME = ? AND UPDATED_TIMESTAMP = DATETIME('NOW') WHERE NAME = ?