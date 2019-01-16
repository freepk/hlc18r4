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

type Like struct {
	ID uint32
	TS uint32
}

type Account struct {
	ID     uint32
	Birth  uint32
	Joined uint32
	//Email  string
	Fname uint8
	Sname uint16
	//Phone     []byte
	Sex     SexEnum
	Country uint8
	City    uint16
	Status  StatusEnum
	//Interests []uint8
	//PremiumFinish uint32
	//PremiumPeriod
	//Likes []Like
}

func (a *Account) Reset() {
	a.ID = 0
	a.Birth = 0
	a.Joined = 0
	//a.Email = ""
	a.Fname = 0
	a.Sname = 0
	//a.Phone = a.Phone[:0]
	a.Sex = 0
	a.Country = 0
	a.City = 0
	a.Status = 0
	//a.Premium.Start = 0
	//a.Premium.Finish = 0
	//a.Interests = a.Interests[:0]
	//a.Likes = a.Likes[:0]
}

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
	case MaleSexStr:
		return t[3:], MaleSex, true
	case FemaleSexStr:
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
	case len(t) > BusyStatusLen && string(t[:BusyStatusLen]) == BusyStatusStr:
		return t[BusyStatusLen:], BusyStatus, true
	case len(t) > FreeStatusLen && string(t[:FreeStatusLen]) == FreeStatusStr:
		return t[FreeStatusLen:], FreeStatus, true
	case len(t) > ComplicatedStatusLen && string(t[:ComplicatedStatusLen]) == ComplicatedStatusStr:
		return t[ComplicatedStatusLen:], ComplicatedStatus, true
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

func (a *Account) UnmarshalJSON(buf []byte) ([]byte, bool) {
	var tail []byte
	var ok bool

	if tail, ok = parse.ParseSymbol(buf, '{'); !ok {
		return buf, false
	}
	for {
		tail = parse.ParseSpaces(tail)
		switch {
		case len(tail) > IdLen && string(tail[:IdLen]) == IdKey:
			if tail, a.ID, ok = parse.ParseUint32(tail[IdLen:]); !ok {
				return buf, false
			}
		case len(tail) > BirthLen && string(tail[:BirthLen]) == BirthKey:
			if tail, a.Birth, ok = parse.ParseUint32(tail[BirthLen:]); !ok {
				return buf, false
			}
		case len(tail) > JoinedLen && string(tail[:JoinedLen]) == JoinedKey:
			if tail, a.Joined, ok = parse.ParseUint32(tail[JoinedLen:]); !ok {
				return buf, false
			}
		//case len(tail) > EmailLen && string(tail[:EmailLen]) == EmailKey:
		//	if tail, a.Email, ok = parse.ParseString(tail[EmailLen:]); !ok {
		//		return buf, false
		//	}
		case len(tail) > FnameLen && string(tail[:FnameLen]) == FnameKey:
			if tail, a.Fname, ok = ParseFname(tail[FnameLen:]); !ok {
				return buf, false
			}
		case len(tail) > SnameLen && string(tail[:SnameLen]) == SnameKey:
			if tail, a.Sname, ok = ParseSname(tail[SnameLen:]); !ok {
				return buf, false
			}
		//case len(tail) > PhoneLen && string(tail[:PhoneLen]) == PhoneKey:
		//	if tail, a.Phone, ok = parse.ParseQuoted(tail[PhoneLen:]); !ok {
		//		return buf, false
		//	}
		case len(tail) > SexLen && string(tail[:SexLen]) == SexKey:
			if tail, a.Sex, ok = ParseSex(tail[SexLen:]); !ok {
				return buf, false
			}
		case len(tail) > CountryLen && string(tail[:CountryLen]) == CountryKey:
			if tail, a.Country, ok = ParseCountry(tail[CountryLen:]); !ok {
				return buf, false
			}
		case len(tail) > CityLen && string(tail[:CityLen]) == CityKey:
			if tail, a.City, ok = ParseCity(tail[CityLen:]); !ok {
				return buf, false
			}
		case len(tail) > StatusLen && string(tail[:StatusLen]) == StatusKey:
			if tail, a.Status, ok = ParseStatus(tail[StatusLen:]); !ok {
				return buf, false
			}
			//case len(tail) > PremiumLen && string(tail[:PremiumLen]) == PremiumKey:
			//	var premium struct {
			//		Start  int
			//		Finish int
			//	}
			//	if tail, ok = parse.ParseSymbol(tail[PremiumLen:], '{'); !ok {
			//		return buf, false
			//	}
			//	for {
			//		tail = parse.ParseSpaces(tail)
			//		switch {
			//		case len(tail) > StartLen && string(tail[:StartLen]) == StartKey:
			//			if tail, premium.Start, ok = parse.ParseInt(tail[StartLen:]); !ok {
			//				return buf, false
			//			}
			//		case len(tail) > FinishLen && string(tail[:FinishLen]) == FinishKey:
			//			if tail, premium.Finish, ok = parse.ParseInt(tail[FinishLen:]); !ok {
			//				return buf, false
			//			}
			//		}
			//		if tail, ok = parse.ParseSymbol(tail, ','); !ok {
			//			break
			//		}
			//	}
			//	if tail, ok = parse.ParseSymbol(tail, '}'); !ok {
			//		return buf, false
			//	}
			//	a.Premium = premium
			//case len(tail) > InterestsLen && string(tail[:InterestsLen]) == InterestsKey:
			//	if tail, ok = parse.ParseSymbol(tail[InterestsLen:], '['); !ok {
			//		return buf, false
			//	}
			//	for {
			//		var interest uint8
			//		if tail, interest, ok = ParseInterest(tail); !ok {
			//			return buf, false
			//		}
			//		a.Interests = append(a.Interests, interest)
			//		if tail, ok = parse.ParseSymbol(tail, ','); !ok {
			//			break
			//		}
			//	}
			//	if tail, ok = parse.ParseSymbol(tail, ']'); !ok {
			//		return buf, false
			//	}
			//case len(tail) > LikesLen && string(tail[:LikesLen]) == LikesKey:
			//	if tail, ok = parse.ParseSymbol(tail[LikesLen:], '['); !ok {
			//		return buf, false
			//	}
			//	for {
			//		var like struct {
			//			ID int
			//			TS int
			//		}
			//		if tail, ok = parse.ParseSymbol(tail, '{'); !ok {
			//			return buf, false
			//		}
			//		for {
			//			tail = parse.ParseSpaces(tail)
			//			switch {
			//			case len(tail) > IdLen && string(tail[:IdLen]) == IdKey:
			//				if tail, like.ID, ok = parse.ParseInt(tail[IdLen:]); !ok {
			//					return buf, false
			//				}
			//			case len(tail) > TsLen && string(tail[:TsLen]) == TsKey:
			//				if tail, like.TS, ok = parse.ParseInt(tail[TsLen:]); !ok {
			//					return buf, false
			//				}
			//			}
			//			if tail, ok = parse.ParseSymbol(tail, ','); !ok {
			//				break
			//			}
			//		}
			//		if tail, ok = parse.ParseSymbol(tail, '}'); !ok {
			//			return buf, false
			//		}
			//		a.Likes = append(a.Likes, like)
			//		if tail, ok = parse.ParseSymbol(tail, ','); !ok {
			//			break
			//		}
			//	}
			//	if tail, ok = parse.ParseSymbol(tail, ']'); !ok {
			//		return buf, false
			//	}
		default:
		}
		if tail, ok = parse.ParseSymbol(tail, ','); !ok {
			break
		}

	}
	if tail, ok = parse.ParseSymbol(tail, '}'); !ok {
		return buf, false
	}
	return tail, true
}
