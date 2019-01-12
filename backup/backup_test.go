package backup

import (
	"testing"
)

func TestRestore(t *testing.T) {
	db, err := Restore("../data/")
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	t.Log(db.State())
}
