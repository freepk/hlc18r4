package backup

import (
	"testing"
)

func TestRestore(t *testing.T) {
	rep, err := Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fail()
	}
	_ = rep
}
