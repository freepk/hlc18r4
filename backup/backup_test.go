package backup

import (
	"testing"
)

func TestRestore(t *testing.T) {
	db, err := Restore("../data/")
	if err != nil {
		t.Fatal(err)
	}
	db.Ping()
	db.BuildIndexes()
}
