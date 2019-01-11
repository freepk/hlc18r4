package lookup

import (
	"fmt"
	"sync"

	"github.com/freepk/hashtab"
	"github.com/spaolacci/murmur3"
)

type Lookup struct {
	sync.RWMutex
	k int
	h *hashtab.HashTab
	a [][]byte
	c []int
}

func NewLookup(power uint8) *Lookup {
	h := hashtab.NewHashTab(power)
	if h == nil {
		return nil
	}
	a := make([][]byte, 1, h.Size())
	c := make([]int, 1, h.Size())
	return &Lookup{k: 1, h: h, a: a, c: c}
}

func (l *Lookup) GetOrGen(v []byte) (int, bool) {
	h := murmur3.Sum64(v)
	l.Lock()
	k, ok := l.h.Get(h)
	if ok {
		l.c[k]++
		l.Unlock()
		return int(k), true
	}
	x := make([]byte, len(v))
	copy(x, v)
	k = uint64(l.k)
	l.h.Set(h, k)
	l.k++
	l.a = append(l.a, x)
	l.c = append(l.c, 1)
	l.Unlock()
	return int(k), false
}

func (l *Lookup) Get(k int) ([]byte, bool) {
	if k < 1 {
		return nil, false
	}
	if k < len(l.a) {
		l.c[k]++
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
		fmt.Println(i, l.c[i], string(l.a[i]))
	}
}
