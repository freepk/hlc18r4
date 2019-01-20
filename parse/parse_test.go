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

func TestParseNumbers(t *testing.T) {
	if x, v, ok := ParseNumbers([]byte("")); v != nil || ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseNumbers([]byte(" ")); v != nil || ok || string(x) != " " {
		t.Fail()
	}
	if x, v, ok := ParseNumbers([]byte(" 12\"")); string(v) != "12" || !ok || string(x) != "\"" {
		t.Fail()
	}
	if x, v, ok := ParseNumbers([]byte(" \"12\"")); v != nil || ok || string(x) != " \"12\"" {
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

//func BenchmarkUnquoteInplace(b *testing.B) {
//	x := []byte("\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b")
//	for i := 0; i < b.N; i++ {
//		UnquoteInplace(x)
//	}
//}

//func BenchmarkAtoiNocheck(b *testing.B) {
//	x := []byte("1234567")
//	for i := 0; i < b.N; i++ {
//		AtoiNocheck(x)
//	}
//}

func BenchmarkParseSpaces(b *testing.B) {
	x := []byte("         ")
	for i := 0; i < b.N; i++ {
		ParseSpaces(x)
	}
}

func BenchmarkParseSymbol(b *testing.B) {
	x := []byte("         { ")
	for i := 0; i < b.N; i++ {
		ParseSymbol(x, '{')
	}
}

func BenchmarkParseNumbers(b *testing.B) {
	x := []byte("         123 ")
	for i := 0; i < b.N; i++ {
		ParseSymbol(x, '{')
	}
}

func BenchmarkParseQuoted(b *testing.B) {
	x := []byte("         \"1234567\"")
	for i := 0; i < b.N; i++ {
		ParseQuoted(x)
	}
}
