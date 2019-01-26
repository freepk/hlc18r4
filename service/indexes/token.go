package indexes

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"gitlab.com/freepk/hlc18r4/dictionary"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

const (
	NullToken = iota
	NotNullToken
)

const (
	MaleToken = iota
	FemaleToken
)

const (
	SingleToken = iota
	InRelToken
	ComplToken
	NotSingleToken
	NotInRelToken
	NotComplToken
)

const (
	PremiumNowToken = 4
)

var (
	phoneCodeDict   = dictionary.NewDictionary(4)
	emailDomainDict = dictionary.NewDictionary(4)
)

var currentTime int

func init() {
	if options, err := ioutil.ReadFile("tmp/data/options.txt"); err == nil {
		if _, now, ok := parse.ParseInt(options); ok {
			currentTime = now
		}
	}
	log.Println("Current time", currentTime)
}

func GetNullToken(b []byte) (int, bool) {
	switch string(b) {
	case `0`:
		return NotNullToken, true
	case `1`:
		return NullToken, true
	}
	return 0, false
}

func GetSexToken(b []byte) (int, bool) {
	switch string(b) {
	case `m`:
		return MaleToken, true
	case `f`:
		return FemaleToken, true
	}
	return 0, false
}

func GetStatusToken(b []byte) (int, bool) {
	switch string(b) {
	case `свободны`:
		return SingleToken, true
	case `заняты`:
		return InRelToken, true
	case `всё сложно`:
		return ComplToken, true
	}
	return 0, false
}

func GetNotStatusToken(b []byte) (int, bool) {
	switch string(b) {
	case `свободны`:
		return NotSingleToken, true
	case `заняты`:
		return NotInRelToken, true
	case `всё сложно`:
		return NotComplToken, true
	}
	return 0, false
}

func GetFnameToken(b []byte) (int, bool) {
	return proto.GetFnameToken(b)
}

func GetSnameToken(b []byte) (int, bool) {
	return proto.GetSnameToken(b)
}

func GetCountryToken(b []byte) (int, bool) {
	return proto.GetCountryToken(b)
}

func GetCityToken(b []byte) (int, bool) {
	return proto.GetCityToken(b)
}

func GetInterestToken(b []byte) (int, bool) {
	return proto.GetInterestToken(b)
}

func birthYearToken(year int) int {
	return year - 1950
}

func birthYearTokenTS(ts int) int {
	year := time.Unix(int64(ts), 0).UTC().Year()
	return birthYearToken(year)
}

func GetBirthYearToken(b []byte) (int, bool) {
	if _, year, ok := parse.ParseInt(b); ok {
		return birthYearToken(year), true
	}
	return 0, false
}

func GetBirthYearTokenTS(b []byte) (int, bool) {
	if _, ts, ok := parse.ParseInt(b); ok {
		return birthYearTokenTS(ts), true
	}
	return 0, false
}

func premiumNow(b []byte) bool {
	if _, premium, ok := parse.ParseInt(b); ok && premium > currentTime {
		return true
	}
	return false
}

func phoneCode(b []byte) ([]byte, bool) {
	if len(b) > 5 && b[1] == '(' && b[5] == ')' {
		return b[2:5], true
	}
	return nil, false
}

func phoneCodeToken(b []byte) int {
	token, _ := phoneCodeDict.AddToken(b)
	return token
}

func GetPhoneCodeToken(b []byte) (int, bool) {
	return phoneCodeDict.Token(b)
}

func emailDomain(b []byte) ([]byte, bool) {
	if tilda := bytes.IndexByte(b, 0x40) + 1; tilda > 0 {
		return b[tilda:], true
	}
	return nil, false
}

func emailDomainToken(b []byte) int {
	token, _ := emailDomainDict.AddToken(b)
	return token
}

func GetEmailDomainToken(b []byte) (int, bool) {
	return emailDomainDict.Token(b)
}
