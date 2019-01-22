package dictionary

import (
	"sync"

	"gitlab.com/freepk/hlc18r4/parse"
)

type Dictionary struct {
	sync.RWMutex
	tokens map[string]int
	values [][]byte
}

func NewDictionary(reserved int) *Dictionary {
	tokens := make(map[string]int)
	values := make([][]byte, reserved)
	return &Dictionary{tokens: tokens, values: values}
}

func (dict *Dictionary) AddToken(value []byte) (int, bool) {
	dict.RLock()
	if token, ok := dict.tokens[string(value)]; ok {
		dict.RUnlock()
		return token, true
	}
	dict.RUnlock()
	dict.Lock()
	if token, ok := dict.tokens[string(value)]; ok {
		dict.Unlock()
		return token, true
	}
	token := len(dict.values)
	dict.tokens[string(value)] = token
	value = parse.Unquote(value)
	dict.tokens[string(value)] = token
	dict.values = append(dict.values, value)
	dict.Unlock()
	return token, false
}

func (dict *Dictionary) Token(value []byte) (int, bool) {
	dict.RLock()
	token, ok := dict.tokens[string(value)]
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
