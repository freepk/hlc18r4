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

func parseFname(b []byte) ([]byte, uint8, bool) {
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

func parseSname(b []byte) ([]byte, uint16, bool) {
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

func parseSex(b []byte) ([]byte, SexEnum, bool) {
	t := parse.ParseSpaces(b)
	if len(t) < 3 {
		return b, 0, false
	}
	switch string(t[:3]) {
	case MaleSexStr:
		return t[3:], MaleSex, true
	case FemaleSexStr:
		return t[3:], FemaleSex, true
	}
	return b, 0, false
}

func parseCountry(b []byte) ([]byte, uint8, bool) {
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

func parseCity(b []byte) ([]byte, uint16, bool) {
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

func parseStatus(b []byte) ([]byte, StatusEnum, bool) {
	t := parse.ParseSpaces(b)
	switch {
	case len(t) > BusyStatusLen && string(t[:BusyStatusLen]) == BusyStatusStr:
		return t[BusyStatusLen:], BusyStatus, true
	case len(t) > FreeStatusLen && string(t[:FreeStatusLen]) == FreeStatusStr:
		return t[FreeStatusLen:], FreeStatus, true
	case len(t) > ComplicatedStatusLen && string(t[:ComplicatedStatusLen]) == ComplicatedStatusStr:
		return t[ComplicatedStatusLen:], ComplicatedStatus, true
	}
	return b, 0, false
}

//func parseInterest(b []byte) ([]byte, byte, bool) {
//	t, v, ok := parse.ParseQuoted(b)
//	if !ok {
//		return b, 0, false
//	}
//	x, err := InterestDict.Identify(v)
//	if err != nil {
//		return b, 0, false
//	}
//	return t, byte(x), true
//}
