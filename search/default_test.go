package search

import (
	"gitlab.com/freepk/hlc18r4/backup"
	"testing"
)

func TestDefaultIndex(t *testing.T) {
	rep, err := backup.Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	index := NewDefaultIndex(rep)
	index.Rebuild()
}
