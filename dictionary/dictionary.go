package dictionary

import (
	"sync"

	"github.com/spaolacci/murmur3"
	"gitlab.com/freepk/hlc18r4/parse"
)

type Dictionary struct {
	sync.RWMutex
	tokens map[uint64]int
	values [][]byte
}

func NewDictionary(reserved int) *Dictionary {
	tokens := make(map[uint64]int)
	values := make([][]byte, reserved)
	return &Dictionary{tokens: tokens, values: values}
}

func (dict *Dictionary) AddToken(value []byte) (int, bool) {
	hash := murmur3.Sum64(value)
	dict.RLock()
	if token, ok := dict.tokens[hash]; ok {
		dict.RUnlock()
		return token, true
	}
	dict.RUnlock()
	dict.Lock()
	if token, ok := dict.tokens[hash]; ok {
		dict.Unlock()
		return token, true
	}
	token := len(dict.values)
	dict.tokens[hash] = token
	value = parse.Unquote(value)
	hash = murmur3.Sum64(value)
	dict.tokens[hash] = token
	dict.values = append(dict.values, value)
	dict.Unlock()
	return token, false
}

func (dict *Dictionary) Token(value []byte) (int, bool) {
	hash := murmur3.Sum64(value)
	dict.RLock()
	token, ok := dict.tokens[hash]
	dict.RUnlock()
	return token, ok
}

func (dict *Dictionary) Value(token int) ([]byte, bool) {
	if token > 0 && token < len(dict.values) {
		return dict.values[token], true
	}
	return nil, false
}

func (dict *Dictionary) Len() int {
	return len(dict.values)
}
