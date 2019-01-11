package lookup

import (
	"fmt"
	"sync"

	"github.com/spaolacci/murmur3"
	"github.com/freepk/hashtab"
)

type Lookup struct {
	sync.RWMutex
	k int
	h *hashtab.HashTab
	a [][]byte
}

func NewLookup(power uint8) (*Lookup, error) {
	h, err := hashtab.NewHashTab(power)
	if err != nil {
		return nil, err
	}
	a := make([][]byte, 1, h.Size())
	return &Lookup{k: 1, a: a, h: h}, nil
}

func (l *Lookup) GetOrGen(v []byte) (int, bool) {
	h := murmur3.Sum64(v)
	l.Lock()
	k, ok := l.h.Get(h)
	if ok {
		l.Unlock()
		return int(k), true
	}
	x := make([]byte, len(v))
	copy(x, v)
	k = uint64(l.k)
	l.h.Set(h, k)
	l.k++
	l.a = append(l.a, x)
	l.Unlock()
	return int(k), false
}

func (l *Lookup) Get(k int) ([]byte, bool) {
	if k < 1 {
		return nil, false
	}
	if k < len(l.a) {
		return l.a[k], true
	}
	return nil, false
}

func (l *Lookup) LastKey() int {
	return l.k
}

func (l *Lookup) Print() {
	k := l.k
	for i := 0; i < k; i++ {
		fmt.Println(i, string(l.a[i]))
	}
}

