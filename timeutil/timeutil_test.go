package timeutil

import "testing"
import "fmt"
import "time"

func TestNextWeekDay(t *testing.T) {

	tt, err := time.Parse("2006-01-02 15:04", "2014-10-16 22:15")
	//t, err := time.Parse("2011-11-24 13:50", "2013-26-10 22:15")
	if err != nil {
		t.Error(err)
	}

	nextWeekDay := NextWeekDay(tt, time.Thursday)

	if nextWeekDay.Weekday() != tt.Weekday() {
		t.Errorf("Error Weekday do not matchup")
	}

	if nextWeekDay.Day() != 23 {
		t.Errorf("Error NextWeekDay did not return correct day")
	}

	fmt.Println(tt, tt.Weekday(), nextWeekDay, nextWeekDay.Weekday())

	fmt.Println(NextMonth(tt))

}

func TestNormalFormat(t *testing.T) {

	const normal = "2006-01-02 15:04"

	tn, err := time.Parse(normal, "2011-02-27 04:40")

	if err != nil {
		t.Error(err)
	}

	if NormalFormat(tn) != "2011-02-27 04:40" {
		t.Error("Format is incorrect")
	}

	if tn.Year() != 2011 {
		t.Error("year is wrong")
	}

	if tn.Day() != 27 {
		t.Error("day is wrong")
	}

	if tn.Month().String() != "February" {
		t.Error("Month  is wrong =>", tn.Month().String())
	}

}
