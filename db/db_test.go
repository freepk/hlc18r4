package db

import (
	"testing"
)

func TestDBRestore(t *testing.T) {
	db := NewDB()
	db.Restore("../data/data.zip")
	db.printStats()
}
