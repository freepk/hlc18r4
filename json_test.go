package main

import "testing"

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
