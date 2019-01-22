package proto

import (
	"gitlab.com/freepk/hlc18r4/dictionary"
	"gitlab.com/freepk/hlc18r4/parse"
)

var (
	fnameDict    = dictionary.NewDictionary(4)
	snameDict    = dictionary.NewDictionary(4)
	countryDict  = dictionary.NewDictionary(4)
	cityDict     = dictionary.NewDictionary(4)
	interestDict = dictionary.NewDictionary(4)
)

const (
	NullToken    = 1
	NotNullToken = 2
)

const (
	MaleSex   = 'm'
	FemaleSex = 'f'
)

const (
	SingleStatus    = 's'
	InRelStatus     = 'r'
	ComplStatus     = 'c'
	NotSingleStatus = '!' + SingleStatus
	NotInRelStatus  = '!' + InRelStatus
	NotComplStatus  = '!' + ComplStatus
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
		return MaleSex, true
	case `f`:
		return FemaleSex, true
	}
	return 0, false
}

func parseSex(b []byte) ([]byte, uint8, bool) {
	tail := parse.SkipSpaces(b)
	if len(tail) > 3 {
		switch string(tail[:3]) {
		case `"m"`:
			return tail[3:], MaleSex, true
		case `"f"`:
			return tail[3:], FemaleSex, true
		}
	}
	return b, 0, false
}

func StatusToken(b []byte) (int, bool) {
	switch string(b) {
	case `свободны`:
		return SingleStatus, true
	case `заняты`:
		return InRelStatus, true
	case `всё сложно`:
		return ComplStatus, true
	}
	return 0, false
}

func parseStatus(b []byte) ([]byte, uint8, bool) {
	tail := parse.SkipSpaces(b)
	switch {
	case len(tail) > 50 && string(tail[:50]) == `"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b"`:
		return tail[50:], SingleStatus, true
	case len(tail) > 38 && string(tail[:38]) == `"\u0437\u0430\u043d\u044f\u0442\u044b"`:
		return tail[38:], InRelStatus, true
	case len(tail) > 57 && string(tail[:57]) == `"\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e"`:
		return tail[57:], ComplStatus, true
	}
	return b, 0, false
}

func IsNullToken(b []byte) (int, bool) {
	switch string(b) {
	case `0`:
		return NotNullToken, true
	case `1`:
		return NullToken, true
	}
	return 0, false
}
