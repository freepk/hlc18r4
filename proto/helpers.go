package proto

import (
	"github.com/freepk/dictionary"
	"gitlab.com/freepk/hlc18r4/parse"
)

var (
	FnameDict    = dictionary.NewDictionary(256)
	SnameDict    = dictionary.NewDictionary(2048)
	CountryDict  = dictionary.NewDictionary(256)
	CityDict     = dictionary.NewDictionary(2048)
	InterestDict = dictionary.NewDictionary(256)
)

func ParseFname(b []byte) ([]byte, uint8, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	x, err := FnameDict.Identify(v)
	if err != nil {
		return b, 0, false
	}
	return t, uint8(x), true
}

func ParseSname(b []byte) ([]byte, uint16, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	x, err := SnameDict.Identify(v)
	if err != nil {
		return b, 0, false
	}
	return t, uint16(x), true
}

func ParseSex(b []byte) ([]byte, SexEnum, bool) {
	t := parse.ParseSpaces(b)
	if len(t) < 3 {
		return b, 0, false
	}
	switch string(t[:3]) {
	case `"m"`:
		return t[3:], MaleSex, true
	case `"f"`:
		return t[3:], FemaleSex, true
	}
	return b, 0, false
}

func ParseCountry(b []byte) ([]byte, uint8, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	x, err := CountryDict.Identify(v)
	if err != nil {
		return b, 0, false
	}
	return t, uint8(x), true
}

func ParseCity(b []byte) ([]byte, uint16, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	x, err := CityDict.Identify(v)
	if err != nil {
		return b, 0, false
	}
	return t, uint16(x), true
}

func ParseStatus(b []byte) ([]byte, StatusEnum, bool) {
	t := parse.ParseSpaces(b)
	switch {
	case len(t) > 38 && string(t[:38]) == `"\u0437\u0430\u043d\u044f\u0442\u044b"`:
		return t[38:], BusyStatus, true
	case len(t) > 50 && string(t[:50]) == `"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b"`:
		return t[50:], FreeStatus, true
	case len(t) > 57 && string(t[:57]) == `"\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e"`:
		return t[57:], ComplicatedStatus, true
	}
	return b, 0, false
}

func ParseInterest(b []byte) ([]byte, uint8, bool) {
	t, v, ok := parse.ParseQuoted(b)
	if !ok {
		return b, 0, false
	}
	x, err := InterestDict.Identify(v)
	if err != nil {
		return b, 0, false
	}
	return t, uint8(x), true
}
