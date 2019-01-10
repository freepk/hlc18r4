package lookup

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLookupGetKeyOrSet(t *testing.T) {
	l := NewLookup(100)
	for i := 1; i <= 50; i++ {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(i))
		k, ok := l.GetKeyOrSet(b)
		if ok {
			t.Fail()
		}
		v, ok := l.GetValue(k)
		if !ok {
			t.Fail()
		}
		if !bytes.Equal(b, v) {
			t.Fail()
		}
	}
}
