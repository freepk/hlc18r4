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
	ii := NewInvertedIndex(rep, InterestsTokens, DefaultParts)
	total := ii.Rebuild()
	t.Log("IndexedAccountsTotal", total)
	total = ii.Rebuild()
	t.Log("IndexedAccountsTotal", total)

}
