package sqlconstants

import (
	"testing"
	"time"
)

func TestGoType(t *testing.T) {
	gotype, err := GoType(SQLITE3_TYPE_INTEGER)
	if err != nil {
		t.Error(err)
		return
	}
	if gotype != "int" {
		t.Error("Mapping of int is not proper")
		return
	}

}

func TestFireTriggerEveryXDaysFrom(t *testing.T) {

	const normal = "2006-01-02 15:04"

	tn, err := time.Parse(normal, "2011-02-27 04:40")

	if err != nil {
		t.Error(err)
	}

	s, err := FireTriggerEveryXDaysFrom(3, tn)

	if err != nil {
		t.Error("Error =>", err)
	}

	if s != "DATETIME('2011-02-27 04:40','+3 days')" {
		t.Error("Expected SQL did not match ")
	}

}

func TestFireTriggerEveryXMonthsFrom(t *testing.T) {

	const normal = "2006-01-02 15:04"

	tn, err := time.Parse(normal, "2011-02-27 04:40")

	if err != nil {
		t.Error(err)
	}

	s, err := FireTriggerEveryXMonthsFrom(3, tn)

	if err != nil {
		t.Error("Error =>", err)
	}

	if s != "DATETIME('2011-02-27 04:40','+3 months')" {
		t.Error("Expected SQL did not match ")
	}

}

func TestFireTriggerEveryXWeeksFrom(t *testing.T) {

	const normal = "2006-01-02 15:04"

	tn, err := time.Parse(normal, "2011-02-27 04:40")

	if err != nil {
		t.Error(err)
	}

	s, err := FireTriggerEveryXWeeksFrom(3, tn)

	if err != nil {
		t.Error("Error =>", err)
	}

	if s != "DATETIME('2011-02-27 04:40','+21 days')" {
		t.Error("Expected SQL did not match ")
	}

}
