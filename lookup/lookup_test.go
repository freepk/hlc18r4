package lookup

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLookupGetKeyOrSet(t *testing.T) {
	l := NewLookup(7)
	if l == nil {
		t.Fail()
	}
	for i := 1; i <= 50; i++ {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(i))
		k, ok := l.GetOrGen(b)
		if ok {
			t.Fail()
		}
		v, ok := l.Get(k)
		if !ok {
			t.Fail()
		}
		if !bytes.Equal(b, v) {
			t.Fail()
		}
	}
}
