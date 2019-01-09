package lookup

import (
	"sync"
)

type Lookup struct {
	sync.RWMutex
	k int
	a [][]byte
	m map[string]int
}

func NewLookup(num int) *Lookup {
	l := &Lookup{
		k: 0,
		a: make([][]byte, 0, num),
		m: make(map[string]int, num)}
	l.GetKeyOrSet([]byte{})
	return l
}

func (l *Lookup) GetKeyOrSet(v []byte) int {
	l.Lock()
	k, ok := l.m[string(v)]
	if ok {
		l.Unlock()
		return k
	}
	x := make([]byte, len(v))
	copy(x, v)
	k = l.k
	l.m[string(x)] = k
	l.k++
	l.a = append(l.a, x)
	l.Unlock()
	return k
}

func (l *Lookup) GetValue(k int) []byte {
	if k < 1 {
		return nil
	}
	if k < len(l.a) {
		return l.a[k]
	}
	return nil
}
