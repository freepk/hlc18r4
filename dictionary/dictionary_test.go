package dictionary

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestDictionary(t *testing.T) {
	dict := NewDictionary()
	for i := 1; i <= 50; i++ {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(i))
		id, ok := dict.AddToken(buf)
		if ok {
			t.Fail()
			return
		}
		val, ok := dict.Value(id)
		if !ok || !bytes.Equal(buf, val) {
			t.Fail()
			return
		}
	}
}
