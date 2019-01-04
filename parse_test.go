package main

import (
	"testing"
)

func TestParseSpaces(t *testing.T) {
	if b := parseSpaces([]byte("")); string(b) != "" {
		t.Fail()
	}
	if b := parseSpaces([]byte("  ")); string(b) != "" {
		t.Fail()
	}
	if b := parseSpaces([]byte("a")); string(b) != "a" {
		t.Fail()
	}
	if b := parseSpaces([]byte(" a")); string(b) != "a" {
		t.Fail()
	}
	if b := parseSpaces([]byte("  a")); string(b) != "a" {
		t.Fail()
	}
}

func TestParseSymbol(t *testing.T) {
	if b, ok := parseSymbol([]byte(""), '{'); ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" "), '{'); ok || string(b) != " " {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte("x{"), '{'); ok || string(b) != "x{" {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" x{"), '{'); ok || string(b) != " x{" {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte("{"), '{'); !ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" {"), '{'); !ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte("{x"), '{'); !ok || string(b) != "x" {
		t.Fail()
	}
	if b, ok := parseSymbol([]byte(" {x"), '{'); !ok || string(b) != "x" {
		t.Fail()
	}
}

func TestParseInt(t *testing.T) {
	if x, b, ok := parseInt([]byte("")); x != 0 || ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" ")); x != 0 || ok || string(b) != " " {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" a1")); x != 0 || ok || string(b) != " a1" {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 1")); x != 1 || !ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 1 ")); x != 1 || !ok || string(b) != " " {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 12")); x != 12 || !ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := parseInt([]byte(" 12 ")); x != 12 || !ok || string(b) != " " {
		t.Fail()
	}
}

func BenchmarkParseSpaces(b *testing.B) {
	x := []byte("         ")
	for i := 0; i < b.N; i++ {
		parseSpaces(x)
	}
}

func BenchmarkParseSymbol(b *testing.B) {
	x := []byte("         {")
	for i := 0; i < b.N; i++ {
		parseSymbol(x, '{')
	}
}

func BenchmarkParseInt(b *testing.B) {
	x := []byte("         1234567")
	for i := 0; i < b.N; i++ {
		parseInt(x)
	}
}
