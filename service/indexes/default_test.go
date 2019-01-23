package indexes

import (
	"gitlab.com/freepk/hlc18r4/backup"
	"gitlab.com/freepk/hlc18r4/inverted"
	"testing"
)

func TestIndexer(t *testing.T) {
	t.Log("Restore")
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Create index")
	source := NewDefaultIndexer(rep)
	index := inverted.NewInverted(source)
	t.Log("Rebuild")
	index.Rebuild()
	t.Log("Rebuild")
	index.Rebuild()
	t.Log("Rebuild")
	index.Rebuild()
	t.Log("Rebuild")
	index.Rebuild()
}
