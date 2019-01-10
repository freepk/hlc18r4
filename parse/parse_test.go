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
	if x, b, ok := ParseInt([]byte("")); x != 0 || ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := ParseInt([]byte(" ")); x != 0 || ok || string(b) != " " {
		t.Fail()
	}
	if x, b, ok := ParseInt([]byte(" a1")); x != 0 || ok || string(b) != " a1" {
		t.Fail()
	}
	if x, b, ok := ParseInt([]byte(" 1")); x != 1 || !ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := ParseInt([]byte(" 1 ")); x != 1 || !ok || string(b) != " " {
		t.Fail()
	}
	if x, b, ok := ParseInt([]byte(" 12")); x != 12 || !ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := ParseInt([]byte(" 12 ")); x != 12 || !ok || string(b) != " " {
		t.Fail()
	}
}

func TestParseQuoted(t *testing.T) {
	if x, b, ok := ParseQuoted([]byte("")); x != nil || ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := ParseQuoted([]byte(" ")); x != nil || ok || string(b) != " " {
		t.Fail()
	}
	if x, b, ok := ParseQuoted([]byte(" aa\"")); x != nil || ok || string(b) != " aa\"" {
		t.Fail()
	}
	if x, b, ok := ParseQuoted([]byte(" \"aa\"")); string(x) != "aa" || !ok || string(b) != "" {
		t.Fail()
	}
	if x, b, ok := ParseQuoted([]byte(" \"aa\" ")); string(x) != "aa" || !ok || string(b) != " " {
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
