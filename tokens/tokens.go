package tokens

import (
	"github.com/freepk/dictionary"
)

const (
	Null = iota
	NotNull
)

const (
	MaleSex = iota
	FemaleSex
)

const (
	SingleStatus = iota
	InRelStatus
	ComplStatus
	NotSingleStatus
	NotInRelStatus
	NotComplStatus
)

const (
	PremiumNow = 4
)

const (
	EpochYear = 1950
)

var (
	fnameDict       = dictionary.NewDictionary(4)
	snameDict       = dictionary.NewDictionary(4)
	countryDict     = dictionary.NewDictionary(4)
	cityDict        = dictionary.NewDictionary(4)
	interestDict    = dictionary.NewDictionary(4)
	phoneCodeDict   = dictionary.NewDictionary(4)
	emailDomainDict = dictionary.NewDictionary(4)
)

func Sex(b []byte) (int, bool) {
	switch string(b) {
	case `m`:
		return MaleSex, true
	case `f`:
		return FemaleSex, true
	}
	return 0, false
}

func Status(b []byte) (int, bool) {
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

func NotStatus(b []byte) (int, bool) {
	switch string(b) {
	case `свободны`:
		return NotSingleStatus, true
	case `заняты`:
		return NotInRelStatus, true
	case `всё сложно`:
		return NotComplStatus, true
	}
	return 0, false
}

func AddFname(b []byte) int {
	k, _ := fnameDict.AddKey(b)
	return k
}

func AddSname(b []byte) int {
	k, _ := snameDict.AddKey(b)
	return k
}

func AddCountry(b []byte) int {
	k, _ := countryDict.AddKey(b)
	return k
}

func AddCity(b []byte) int {
	k, _ := cityDict.AddKey(b)
	return k
}

func AddInterest(b []byte) int {
	k, _ := interestDict.AddKey(b)
	return k
}

func AddPhoneCode(b []byte) int {
	k, _ := phoneCodeDict.AddKey(b)
	return k
}

func AddEmailDomain(b []byte) int {
	k, _ := emailDomainDict.AddKey(b)
	return k
}

func Fname(b []byte) (int, bool) {
	return fnameDict.Key(b)
}

func Sname(b []byte) (int, bool) {
	return snameDict.Key(b)
}

func Country(b []byte) (int, bool) {
	return countryDict.Key(b)
}

func City(b []byte) (int, bool) {
	return cityDict.Key(b)
}

func Interest(b []byte) (int, bool) {
	return interestDict.Key(b)
}

func PhoneCode(b []byte) (int, bool) {
	return phoneCodeDict.Key(b)
}

func EmailDomain(b []byte) (int, bool) {
	return emailDomainDict.Key(b)
}

//func Year(b []byte) (int, bool) {
//	return 0, false
//}

//func YearTS(b []byte) (int, bool) {
//	return 0, false
//}

func FnameVal(k int) ([]byte, bool) {
	return fnameDict.Val(k)
}

func SnameVal(k int) ([]byte, bool) {
	return snameDict.Val(k)
}

func CountryVal(k int) ([]byte, bool) {
	return countryDict.Val(k)
}

func CityVal(k int) ([]byte, bool) {
	return cityDict.Val(k)
}

func InterestVal(k int) ([]byte, bool) {
	return interestDict.Val(k)
}
