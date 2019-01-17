package backup

import (
	"gitlab.com/freepk/hlc18r4/service"
	"testing"
)

func TestRestore(t *testing.T) {
	rep, err := Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fail()
	}
	svc := service.NewAccountsService(rep)
	_ = svc
}
