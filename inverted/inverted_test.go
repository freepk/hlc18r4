package inverted

import (
	"testing"

	"gitlab.com/freepk/hlc18r4/backup"
)

func TestRebuild(t *testing.T) {
	rep, err := backup.Restore("../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Rebuilding interests")
	ii := NewInvertedIndex(rep, DefaultParts, InterestsTokens)
	total, grow := ii.Rebuild()
	t.Log("total", total, grow)
	total, grow = ii.Rebuild()
	t.Log("total", total, grow)
}
