package proto

import (
	"gitlab.com/freepk/hlc18r4/dictionary"
	"gitlab.com/freepk/hlc18r4/parse"
)

var (
	fnameDict    = dictionary.NewDictionary()
	snameDict    = dictionary.NewDictionary()
	countryDict  = dictionary.NewDictionary()
	cityDict     = dictionary.NewDictionary()
	interestDict = dictionary.NewDictionary()
)

func parseFname(b []byte) ([]byte, uint8, bool) {
	tail, value, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	token, _ := fnameDict.AddToken(value)
	return tail, uint8(token), true
}

func FnameToken(value []byte) (int, bool) {
	return fnameDict.Token(value)
}

func parseSname(b []byte) ([]byte, uint16, bool) {
	tail, value, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	token, _ := snameDict.AddToken(value)
	return tail, uint16(token), true
}

func SnameToken(b []byte) (int, bool) {
	return snameDict.Token(b)
}

func parseCountry(b []byte) ([]byte, uint8, bool) {
	tail, value, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	token, _ := countryDict.AddToken(value)
	return tail, uint8(token), true
}

func CountryToken(b []byte) (int, bool) {
	return countryDict.Token(b)
}

func parseCity(b []byte) ([]byte, uint16, bool) {
	tail, value, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	token, _ := cityDict.AddToken(value)
	return tail, uint16(token), true
}

func CityToken(b []byte) (int, bool) {
	return cityDict.Token(b)
}

func parseInterest(b []byte) ([]byte, uint8, bool) {
	tail, value, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	token, _ := interestDict.AddToken(value)
	return tail, uint8(token), true
}

func InterestToken(b []byte) (int, bool) {
	return interestDict.Token(b)
}

func SexToken(b []byte) (int, bool) {
	switch string(b) {
	case `m`:
		return 1, true
	case `f`:
		return 2, true
	}
	return 0, false
}

func parseSex(b []byte) ([]byte, uint8, bool) {
	tail := parse.SkipSpaces(b)
	if len(tail) > 3 {
		switch string(tail[:3]) {
		case `"m"`:
			return tail[3:], 1, true
		case `"f"`:
			return tail[3:], 2, true
		}
	}
	return b, 0, false
}

func StatusToken(b []byte) (int, bool) {
	switch string(b) {
	case `свободны`:
		return 1, true
	case `заняты`:
		return 2, true
	case `все сложно`:
		return 3, true
	}
	return 0, false
}

func parseStatus(b []byte) ([]byte, uint8, bool) {
	tail := parse.SkipSpaces(b)
	switch {
	case len(tail) > 50 && string(tail[:50]) == `"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b"`:
		return tail[50:], 1, true
	case len(tail) > 38 && string(tail[:38]) == `"\u0437\u0430\u043d\u044f\u0442\u044b"`:
		return tail[38:], 2, true
	case len(tail) > 57 && string(tail[:57]) == `"\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e"`:
		return tail[57:], 3, true
	}
	return b, 0, false
}