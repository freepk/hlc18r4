package main

import (
	//"fmt"
	"os"
	"testing"
)

func TestJsonFindObj(t *testing.T) {
	if _, _, ok := jsonFindObj([]byte("")); ok {
		t.Fail()
	}
	if _, _, ok := jsonFindObj([]byte("{")); ok {
		t.Fail()
	}
	if _, _, ok := jsonFindObj([]byte("{{}")); ok {
		t.Fail()
	}
	if i, j, ok := jsonFindObj([]byte("{}")); !ok || i != 0 || j != 1 {
		t.Fail()
	}
}

func TestJsonFindQuoted(t *testing.T) {
	if _, _, ok := jsonFindQuoted([]byte("")); ok {
		t.Fail()
	}
	if _, _, ok := jsonFindQuoted([]byte("\"")); ok {
		t.Fail()
	}
	if i, j, ok := jsonFindQuoted([]byte("\"\"")); !ok || i != 0 || j != 1 {
		t.Fail()
	}
}

func TestJsonReadObj(t *testing.T) {
	return
	f, err := os.Open("test.big.json")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer f.Close()
	b := make([]byte, 8192)
	u := make([]byte, 8192)
	jsonReadObj(f, b, 0, func(b []byte) error {
		unquote(u, b)
		return nil
	})
}
