package indexes

import (
	"github.com/freepk/iterator"
	"gitlab.com/freepk/hlc18r4/backup"
	"testing"
)

func TestIndexer(t *testing.T) {
	t.Log("Restore")
	rep, err := backup.Restore("../../tmp/data/data.zip")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Create index")
	index := NewDefaultIndex(rep)
	t.Log("Rebuild")
	index.Rebuild()
	it := iterator.Iterator(index.Country(NotNullToken))
	it = iterator.NewInterIter(it, index.Sex(MaleToken))
	limit := 20
	for limit > 0 {
		limit--
		pseudo, ok := it.Next()
		if !ok {
			break
		}
		id := 2000000 - pseudo
		t.Log(limit, id)
	}
}
