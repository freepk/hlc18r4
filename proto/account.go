package proto

import (
	"gitlab.com/freepk/hlc18r4/parse"
)

const EmailBufMaxLen = 40

const PhoneBufMaxLen = 15

const NumberBufMaxLen = 10

const InterestsMaxLen = 10

type EmailBuf struct {
	Len uint8
	Buf [EmailBufMaxLen]byte
}

type PhoneBuf [PhoneBufMaxLen]byte

type NumberBuf [NumberBufMaxLen]byte

type InterestsBuf [InterestsMaxLen]uint8

type Like struct {
	ID uint32
	TS uint32
}

type Account struct {
	ID            NumberBuf
	Birth         NumberBuf
	Joined        NumberBuf
	Email         EmailBuf
	Fname         uint8
	Sname         uint16
	Phone         PhoneBuf
	Sex           SexEnum
	Country       uint8
	City          uint16
	Status        StatusEnum
	PremiumStart  NumberBuf
	PremiumFinish NumberBuf
	Interests     InterestsBuf
	LikesTo       []Like
}

func (a *Account) reset() {
	var number NumberBuf
	var email EmailBuf
	var phone PhoneBuf
	var interests InterestsBuf

	a.ID = number
	a.Birth = number
	a.Joined = number
	a.Email = email
	a.Fname = 0
	a.Sname = 0
	a.Phone = phone
	a.Sex = 0
	a.Country = 0
	a.City = 0
	a.Status = 0
	a.PremiumStart = number
	a.PremiumFinish = number
	a.Interests = interests
	a.LikesTo = a.LikesTo[:0]
}

const (
	IDField        = 1
	BirthField     = 2
	JoinedField    = 4
	EmailField     = 8
	FnameField     = 16
	SnameField     = 32
	PhoneField     = 64
	SexField       = 128
	CountryField   = 256
	CityField      = 512
	StatusField    = 1024
	PremiumField   = 2048
	InterestsField = 4096
	LikesToField   = 8192
)

func (a *Account) MarshalToJSON(fields int, buf []byte) []byte {
	fields = (1 << 20) - 1
	buf = append(buf, '{')
	buf = append(buf, IdKey...)
	buf = append(buf, string(a.ID[:])...)
	if (fields & BirthField) == BirthField {
		buf = append(buf, ',')
		buf = append(buf, BirthKey...)
		buf = append(buf, string(a.Birth[:])...)
	}
	if (fields & JoinedField) == JoinedField {
		buf = append(buf, ',')
		buf = append(buf, JoinedKey...)
		buf = append(buf, string(a.Joined[:])...)
	}
	if (fields & EmailField) == EmailField {
		buf = append(buf, ',')
		buf = append(buf, EmailKey...)
		buf = append(buf, a.Email.Buf[:a.Email.Len]...)
	}
	if (fields&FnameField) == FnameField && a.Fname > 0 {
		fname, _ := FnameDict.Value(uint64(a.Fname))
		buf = append(buf, ',')
		buf = append(buf, FnameKey...)
		buf = append(buf, fname...)
	}
	if (fields&SnameField) == SnameField && a.Sname > 0 {
		sname, _ := SnameDict.Value(uint64(a.Fname))
		buf = append(buf, ',')
		buf = append(buf, SnameKey...)
		buf = append(buf, sname...)
	}
	if (fields&PhoneField) == PhoneField && a.Phone[0] > 0 {
		buf = append(buf, ',')
		buf = append(buf, PhoneKey...)
		buf = append(buf, a.Phone[:]...)
	}
	if (fields & SexField) == SexField {
		switch a.Sex {
		case MaleSex:
			buf = append(buf, `,"sex":"m"`...)
		case FemaleSex:
			buf = append(buf, `,"sex":"f"`...)
		}
	}
	if (fields&CountryField) == CountryField && a.Country > 0 {
		country, _ := CountryDict.Value(uint64(a.Country))
		buf = append(buf, ',')
		buf = append(buf, CountryKey...)
		buf = append(buf, country...)
	}
	if (fields&CityField) == CityField && a.City > 0 {
		city, _ := CountryDict.Value(uint64(a.City))
		buf = append(buf, ',')
		buf = append(buf, CityKey...)
		buf = append(buf, city...)
	}
	if (fields & StatusField) == StatusField {
		switch a.Status {
		case FreeStatus:
			buf = append(buf, `,"status":"свободны"`...)
		case BusyStatus:
			buf = append(buf, `,"status":"заняты"`...)
		case ComplicatedStatus:
			buf = append(buf, `,"status":"всё сложно"`...)
		}
	}
	if (fields&PremiumField) == PremiumField && a.PremiumFinish[0] > 0 {
		buf = append(buf, `',"premium":{"finish":`...)
		buf = append(buf, string(a.PremiumFinish[:])...)
		buf = append(buf, `,"start":`...)
		buf = append(buf, string(a.PremiumStart[:])...)
		buf = append(buf, '}')
	}
	buf = append(buf, '}')
	return buf
}

func (a *Account) UnmarshalJSON(buf []byte) ([]byte, bool) {
	var tail []byte
	var temp []byte
	var ok bool

	a.reset()

	if tail, ok = parse.ParseSymbol(buf, '{'); !ok {
		return buf, false
	}
	for {
		tail = parse.ParseSpaces(tail)
		switch {
		case len(tail) > IdLen && string(tail[:IdLen]) == IdKey:
			if tail, temp, ok = parse.ParseNumber(tail[IdLen:]); !ok {
				return buf, false
			}
			copy(a.ID[:], temp)
		case len(tail) > BirthLen && string(tail[:BirthLen]) == BirthKey:
			if tail, temp, ok = parse.ParseNumber(tail[BirthLen:]); !ok {
				return buf, false
			}
			copy(a.Birth[:], temp)
		case len(tail) > JoinedLen && string(tail[:JoinedLen]) == JoinedKey:
			if tail, temp, ok = parse.ParseNumber(tail[JoinedLen:]); !ok {
				return buf, false
			}
			copy(a.Joined[:], temp)
		case len(tail) > EmailLen && string(tail[:EmailLen]) == EmailKey:
			if tail, temp, ok = parse.ParseQuoted(tail[EmailLen:]); !ok {
				return buf, false
			}
			a.Email.Len = uint8(copy(a.Email.Buf[:], temp))
		case len(tail) > FnameLen && string(tail[:FnameLen]) == FnameKey:
			if tail, a.Fname, ok = parseFname(tail[FnameLen:]); !ok {
				return buf, false
			}
		case len(tail) > SnameLen && string(tail[:SnameLen]) == SnameKey:
			if tail, a.Sname, ok = parseSname(tail[SnameLen:]); !ok {
				return buf, false
			}
		case len(tail) > PhoneLen && string(tail[:PhoneLen]) == PhoneKey:
			if tail, temp, ok = parse.ParseQuoted(tail[PhoneLen:]); !ok {
				return buf, false
			}
			copy(a.Phone[:], temp)
		case len(tail) > SexLen && string(tail[:SexLen]) == SexKey:
			if tail, a.Sex, ok = parseSex(tail[SexLen:]); !ok {
				return buf, false
			}
		case len(tail) > CountryLen && string(tail[:CountryLen]) == CountryKey:
			if tail, a.Country, ok = parseCountry(tail[CountryLen:]); !ok {
				return buf, false
			}
		case len(tail) > CityLen && string(tail[:CityLen]) == CityKey:
			if tail, a.City, ok = parseCity(tail[CityLen:]); !ok {
				return buf, false
			}
		case len(tail) > StatusLen && string(tail[:StatusLen]) == StatusKey:
			if tail, a.Status, ok = parseStatus(tail[StatusLen:]); !ok {
				return buf, false
			}
		case len(tail) > PremiumLen && string(tail[:PremiumLen]) == PremiumKey:
			if tail, ok = parse.ParseSymbol(tail[PremiumLen:], '{'); !ok {
				return buf, false
			}
			for {
				tail = parse.ParseSpaces(tail)
				switch {
				case len(tail) > StartLen && string(tail[:StartLen]) == StartKey:
					if tail, temp, ok = parse.ParseNumber(tail[StartLen:]); !ok {
						return buf, false
					}
					copy(a.PremiumStart[:], temp)
				case len(tail) > FinishLen && string(tail[:FinishLen]) == FinishKey:
					if tail, temp, ok = parse.ParseNumber(tail[FinishLen:]); !ok {
						return buf, false
					}
					copy(a.PremiumFinish[:], temp)
				}
				if tail, ok = parse.ParseSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.ParseSymbol(tail, '}'); !ok {
				return buf, false
			}
		case len(tail) > InterestsLen && string(tail[:InterestsLen]) == InterestsKey:
			if tail, ok = parse.ParseSymbol(tail[InterestsLen:], '['); !ok {
				return buf, false
			}
			var i uint8
			var interest uint8
			for {
				if tail, interest, ok = parseInterest(tail); !ok {
					return buf, false
				}
				a.Interests[i] = interest
				i++
				if tail, ok = parse.ParseSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.ParseSymbol(tail, ']'); !ok {
				return buf, false
			}
		case len(tail) > LikesLen && string(tail[:LikesLen]) == LikesKey:
			if tail, ok = parse.ParseSymbol(tail[LikesLen:], '['); !ok {
				return buf, false
			}
			for {
				ID := 0
				TS := 0
				if tail, ok = parse.ParseSymbol(tail, '{'); !ok {
					return buf, false
				}
				for {
					tail = parse.ParseSpaces(tail)
					switch {
					case len(tail) > IdLen && string(tail[:IdLen]) == IdKey:
						if tail, ID, ok = parse.ParseInt(tail[IdLen:]); !ok {
							return buf, false
						}
					case len(tail) > TsLen && string(tail[:TsLen]) == TsKey:
						if tail, TS, ok = parse.ParseInt(tail[TsLen:]); !ok {
							return buf, false
						}
					}
					if tail, ok = parse.ParseSymbol(tail, ','); !ok {
						break
					}
				}
				if tail, ok = parse.ParseSymbol(tail, '}'); !ok {
					return buf, false
				}
				a.LikesTo = append(a.LikesTo, Like{ID: uint32(ID), TS: uint32(TS)})
				if tail, ok = parse.ParseSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.ParseSymbol(tail, ']'); !ok {
				return buf, false
			}
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
