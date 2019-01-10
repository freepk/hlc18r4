package hashmap

import (
	"sync"

	"github.com/spaolacci/murmur3"
)

type HashMap struct {
	sync.RWMutex
	m map[uint64]uint32
}

func NewHashMap() *HashMap {
	return &HashMap{m: make(map[uint64]uint32)}
}

func (m *HashMap) GetOrSet(k []byte, v int) (int, bool) {
	h := murmur3.Sum64(k)
	m.Lock()
	a, ok := m.m[h]
	if ok {
		m.Unlock()
		return int(a), ok
	}
	m.m[h] = uint32(v)
	m.Unlock()
	return int(v), false
}
