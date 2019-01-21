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
	interests := NewInvertedIndex(rep, DefaultParts, InterestsTokens)

	total, grow := 0, 0

	t.Log("Frist pass")
	total, grow = interests.Rebuild()

	t.Log("Second pass")

	total, grow = interests.Rebuild()
	t.Log("Interests", total, grow)
}
