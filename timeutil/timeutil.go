package timeutil

import "time"

func NextXthWeekFromDate(t time.Time, x int) time.Time {
	x = x * 7
	return t.AddDate(0, 0, x)
}

func NextXthMonthFromDate(t time.Time, x int) time.Time {
	return t.AddDate(0, x, 0)
}

func NextXthYearFromDate(t time.Time, x int) time.Time {
	return t.AddDate(x, 0, 0)
}

func NextXthDayFromDate(t time.Time, x int) time.Time {
	return t.AddDate(0, 0, x)
}

func NormalFormat(t time.Time) string {
	// returns string in normal format
	//YYYY-DD-MM HH:MM
	const normal = "2006-01-02 15:04"
	return t.Format(normal)

}
