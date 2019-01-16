package parse

import (
	"testing"
)

func TestParseSpaces(t *testing.T) {
	if b := ParseSpaces([]byte("")); string(b) != "" {
		t.Fail()
	}
	if b := ParseSpaces([]byte("  ")); string(b) != "" {
		t.Fail()
	}
	if b := ParseSpaces([]byte("a")); string(b) != "a" {
		t.Fail()
	}
	if b := ParseSpaces([]byte(" a")); string(b) != "a" {
		t.Fail()
	}
	if b := ParseSpaces([]byte("  a")); string(b) != "a" {
		t.Fail()
	}
}

func TestParseSymbol(t *testing.T) {
	if b, ok := ParseSymbol([]byte(""), '{'); ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte(" "), '{'); ok || string(b) != " " {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte("x{"), '{'); ok || string(b) != "x{" {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte(" x{"), '{'); ok || string(b) != " x{" {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte("{"), '{'); !ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte(" {"), '{'); !ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte("{x"), '{'); !ok || string(b) != "x" {
		t.Fail()
	}
	if b, ok := ParseSymbol([]byte(" {x"), '{'); !ok || string(b) != "x" {
		t.Fail()
	}
}

func TestParseInt(t *testing.T) {
	if x, v, ok := ParseInt([]byte("")); v != 0 || ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseInt([]byte(" ")); v != 0 || ok || string(x) != " " {
		t.Fail()
	}
	if x, v, ok := ParseInt([]byte(" a1")); v != 0 || ok || string(x) != " a1" {
		t.Fail()
	}
	if x, v, ok := ParseInt([]byte(" 1")); v != 1 || !ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseInt([]byte(" 1 ")); v != 1 || !ok || string(x) != " " {
		t.Fail()
	}
	if x, v, ok := ParseInt([]byte(" 12")); v != 12 || !ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseInt([]byte(" 12 ")); v != 12 || !ok || string(x) != " " {
		t.Fail()
	}
}

func TestParseQuoted(t *testing.T) {
	if x, v, ok := ParseQuoted([]byte("")); v != nil || ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseQuoted([]byte(" ")); v != nil || ok || string(x) != " " {
		t.Fail()
	}
	if x, v, ok := ParseQuoted([]byte(" aa\"")); v != nil || ok || string(x) != " aa\"" {
		t.Fail()
	}
	if x, v, ok := ParseQuoted([]byte(" \"aa\"")); string(v) != "aa" || !ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseQuoted([]byte(" \"aa\" ")); string(v) != "aa" || !ok || string(x) != " " {
		t.Fail()
	}
}

func BenchmarkParseSpaces(b *testing.B) {
	x := []byte("         ")
	for i := 0; i < b.N; i++ {
		ParseSpaces(x)
	}
}

func BenchmarkParseSymbol(b *testing.B) {
	x := []byte("         {")
	for i := 0; i < b.N; i++ {
		ParseSymbol(x, '{')
	}
}

func BenchmarkParseInt(b *testing.B) {
	x := []byte("         1234567")
	for i := 0; i < b.N; i++ {
		ParseInt(x)
	}
}

func BenchmarkParseQuoted(b *testing.B) {
	x := []byte("         \"1234567\"")
	for i := 0; i < b.N; i++ {
		ParseQuoted(x)
	}
}
