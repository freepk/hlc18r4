package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestLookupGetOrSetIndex(t *testing.T) {
	l := NewLookup(100, 40)
	for i := 0; i < 50; i++ {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(i))
		index := l.GetIndexOrSet(b)
		item := l.GetItem(index)
		if !bytes.Equal(b, item) {
			t.Fail()
		}
	}
}
