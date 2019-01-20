package parse

import (
	"testing"
)

func TestSkipSpaces(t *testing.T) {
	if b := SkipSpaces([]byte("")); string(b) != "" {
		t.Fail()
	}
	if b := SkipSpaces([]byte("  ")); string(b) != "" {
		t.Fail()
	}
	if b := SkipSpaces([]byte("a")); string(b) != "a" {
		t.Fail()
	}
	if b := SkipSpaces([]byte(" a")); string(b) != "a" {
		t.Fail()
	}
	if b := SkipSpaces([]byte("  a")); string(b) != "a" {
		t.Fail()
	}
}

func TestSkipSymbol(t *testing.T) {
	if b, ok := SkipSymbol([]byte(""), '{'); ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte(" "), '{'); ok || string(b) != " " {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte("x{"), '{'); ok || string(b) != "x{" {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte(" x{"), '{'); ok || string(b) != " x{" {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte("{"), '{'); !ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte(" {"), '{'); !ok || string(b) != "" {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte("{x"), '{'); !ok || string(b) != "x" {
		t.Fail()
	}
	if b, ok := SkipSymbol([]byte(" {x"), '{'); !ok || string(b) != "x" {
		t.Fail()
	}
}

func TestParseNumber(t *testing.T) {
	if x, v, ok := ParseNumber([]byte("")); v != nil || ok || string(x) != "" {
		t.Fail()
	}
	if x, v, ok := ParseNumber([]byte(" ")); v != nil || ok || string(x) != " " {
		t.Fail()
	}
	if x, v, ok := ParseNumber([]byte(" 12\"")); string(v) != "12" || !ok || string(x) != "\"" {
		t.Fail()
	}
	if x, v, ok := ParseNumber([]byte(" \"12\"")); v != nil || ok || string(x) != " \"12\"" {
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
	if x, v, ok := ParseQuoted([]byte(`"\u041b\u0435\u0431\u0435\u0442\u0430\u0442\u0435\u0432"}`)); string(v) != `\u041b\u0435\u0431\u0435\u0442\u0430\u0442\u0435\u0432` || !ok || string(x) != "}" {
		t.Fail()
	}
	return
}

func BenchmarkSkipSpaces(b *testing.B) {
	x := []byte("         ")
	for i := 0; i < b.N; i++ {
		SkipSpaces(x)
	}
}

func BenchmarkSkipSymbol(b *testing.B) {
	x := []byte("         { ")
	for i := 0; i < b.N; i++ {
		SkipSymbol(x, '{')
	}
}

func BenchmarkParseNumber(b *testing.B) {
	x := []byte("         12345678 ")
	for i := 0; i < b.N; i++ {
		ParseNumber(x)
	}
}

func BenchmarkParseInt(b *testing.B) {
	x := []byte("         12345678 ")
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
