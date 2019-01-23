package service

import (
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/index"
	"testing"
)

func TestIndexer(t *testing.T) {
	t.Log("Restore")
	rep, err := backup.Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fail()
	}
	t.Log("Create index")
	ir := NewAccountsIndexer(rep)
	ix := index.NewInverted(ir)
	ix.Rebuild()
}
