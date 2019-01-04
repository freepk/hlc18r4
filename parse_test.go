package main

import (
	"bytes"
	"testing"
)

func TestParseSymbol(t *testing.T) {
	if b, ok := parseSymbol([]byte(""), '{'); ok || !bytes.Equal(b, []byte("")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" "), '{'); ok || !bytes.Equal(b, []byte(" ")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte("x{"), '{'); ok || !bytes.Equal(b, []byte("x{")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" x{"), '{'); ok || !bytes.Equal(b, []byte(" x{")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte("{"), '{'); !ok || !bytes.Equal(b, []byte("")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" {"), '{'); !ok || !bytes.Equal(b, []byte("")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte("{x"), '{'); !ok || !bytes.Equal(b, []byte("x")) {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" {x"), '{'); !ok || !bytes.Equal(b, []byte("x")) {
		t.Fail()
	}
}

func TestParseInt(t *testing.T) {
	if x, b, ok := parseInt([]byte("")); x != 0 || ok || !bytes.Equal(b, []byte("")) {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" ")); x != 0 || ok || !bytes.Equal(b, []byte(" ")) {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" a1")); x != 0 || ok || !bytes.Equal(b, []byte(" a1")) {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 1")); x != 1 || !ok || !bytes.Equal(b, []byte("")) {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 1 ")); x != 1 || !ok || !bytes.Equal(b, []byte(" ")) {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 12")); x != 12 || !ok || !bytes.Equal(b, []byte("")) {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 12 ")); x != 12 || !ok || !bytes.Equal(b, []byte(" ")) {
		t.Fail()
	}
}

func BenchmarkParseSymbol(b *testing.B) {
	s := []byte("         {")
	for i := 0; i < b.N; i++ {
		parseSymbol(s, '{')
	}
}

func BenchmarkParseInt(b *testing.B) {
	s := []byte("         1")
	for i := 0; i < b.N; i++ {
		parseInt(s)
	}
}
