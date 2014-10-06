package db

import (
	"fmt"
	"testing"
)

func TestNewConnection(t *testing.T) {
	c, err := NewConnection("sqlite3")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(c)
}
