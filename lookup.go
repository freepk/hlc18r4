package main

import (
	"sync"
)

type Lookup struct {
	sync.RWMutex
	offset int
	index  int
	items  [][]byte
	lookup map[string]int
}

func NewLookup(num, offset int) *Lookup {
	return &Lookup{offset: offset,
		items:  make([][]byte, 0, num),
		lookup: make(map[string]int, num)}
}

func (l *Lookup) GetIndex(item []byte) int {
	l.RLock()
	index, ok := l.lookup[string(item)]
	l.RUnlock()
	if ok {
		return index
	}
	return -1
}

func (l *Lookup) GetIndexOrSet(item []byte) int {
	l.Lock()
	index, ok := l.lookup[string(item)]
	if ok {
		l.Unlock()
		return index
	}
	temp := make([]byte, len(item))
	copy(temp, item)
	index = l.index + l.offset
	l.lookup[string(item)] = index
	l.index++
	l.items = append(l.items, temp)
	l.Unlock()
	return index
}

func (l *Lookup) GetItem(index int) []byte {
	if index-l.offset < 0 {
		return nil
	}
	index -= l.offset
	l.RLock()
	if index < len(l.items) {
		l.RUnlock()
		return l.items[index]
	}
	l.RUnlock()
	return nil
}

func (l *Lookup) GetItemNoLock(index int) []byte {
	if index-l.offset < 0 {
		return nil
	}
	index -= l.offset
	if index < len(l.items) {
		return l.items[index]
	}
	return nil
}
