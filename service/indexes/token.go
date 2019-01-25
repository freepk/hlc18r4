package indexes

import (
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
	PremiumNullToken    = NullToken
	PremiumNotNullToken = NotNullToken
	PremiumNowToken     = 4
)

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
