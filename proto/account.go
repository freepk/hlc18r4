package proto

import (
	"gitlab.com/freepk/hlc18r4/parse"
)

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

type SexEnum byte

const (
	_          = iota
	FreeStatus = StatusEnum(iota)
	BusyStatus
	ComplicatedStatus
)

type StatusEnum uint8

const (
	_       = iota
	MaleSex = SexEnum(iota)
	FemaleSex
)

const (
	EmailBufMaxLen  = 40
	PhoneBufMaxLen  = 15
	NumberBufMaxLen = 10
	InterestsMaxLen = 10
)

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

func trim(b []byte) []byte {
	n := len(b)
	for n > 0 {
		if b[n-1] > 0 {
			break
		}
		n--
	}
	return b[:n]
}

func (a *Account) MarshalToJSON(fields int, buf []byte) []byte {
	buf = append(buf, `{"id":`...)
	buf = append(buf, trim(a.ID[:])...)
	if (fields & BirthField) == BirthField {
		buf = append(buf, `,"birth":`...)
		buf = append(buf, trim(a.Birth[:])...)
	}
	if (fields & JoinedField) == JoinedField {
		buf = append(buf, `,"joined":`...)
		buf = append(buf, trim(a.Joined[:])...)
	}
	if (fields & EmailField) == EmailField {
		buf = append(buf, `,"email":`...)
		buf = append(buf, a.Email.Buf[:a.Email.Len]...)
	}
	if (fields&FnameField) == FnameField && a.Fname > 0 {
		fname, _ := FnameDict.Value(uint64(a.Fname))
		buf = append(buf, `,"fname":`...)
		buf = append(buf, fname...)
	}
	if (fields&SnameField) == SnameField && a.Sname > 0 {
		sname, _ := SnameDict.Value(uint64(a.Sname))
		buf = append(buf, `,"sname":`...)
		buf = append(buf, sname...)
	}
	if (fields&PhoneField) == PhoneField && a.Phone[0] > 0 {
		buf = append(buf, `,"phone":`...)
		buf = append(buf, trim(a.Phone[:])...)
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
		buf = append(buf, `,"country":`...)
		buf = append(buf, country...)
	}
	if (fields&CityField) == CityField && a.City > 0 {
		city, _ := CityDict.Value(uint64(a.City))
		buf = append(buf, `,"city":`...)
		buf = append(buf, city...)
	}
	if (fields & StatusField) == StatusField {
		switch a.Status {
		case FreeStatus:
			buf = append(buf, `,"status":\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b"`...)
		case BusyStatus:
			buf = append(buf, `,"status":"\u0437\u0430\u043d\u044f\u0442\u044b"`...)
		case ComplicatedStatus:
			buf = append(buf, `,"status":"\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e"`...)
		}
	}
	if (fields&PremiumField) == PremiumField && a.PremiumFinish[0] > 0 {
		buf = append(buf, `',"premium":{"finish":`...)
		buf = append(buf, trim(a.PremiumFinish[:])...)
		buf = append(buf, `,"start":`...)
		buf = append(buf, trim(a.PremiumStart[:])...)
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
		case len(tail) > 5 && string(tail[:5]) == `"id":`:
			if tail, temp, ok = parse.ParseNumbers(tail[5:]); !ok {
				return buf, false
			}
			copy(a.ID[:], temp)
		case len(tail) > 8 && string(tail[:8]) == `"birth":`:
			if tail, temp, ok = parse.ParseNumbers(tail[8:]); !ok {
				return buf, false
			}
			copy(a.Birth[:], temp)
		case len(tail) > 9 && string(tail[:9]) == `"joined":`:
			if tail, temp, ok = parse.ParseNumbers(tail[9:]); !ok {
				return buf, false
			}
			copy(a.Joined[:], temp)
		case len(tail) > 8 && string(tail[:8]) == `"email":`:
			if tail, temp, ok = parse.ParseQuoted(tail[8:]); !ok {
				return buf, false
			}
			a.Email.Len = uint8(copy(a.Email.Buf[:], temp))
		case len(tail) > 8 && string(tail[:8]) == `"fname":`:
			if tail, a.Fname, ok = ParseFname(tail[8:]); !ok {
				return buf, false
			}
		case len(tail) > 8 && string(tail[:8]) == `"sname":`:
			if tail, a.Sname, ok = ParseSname(tail[8:]); !ok {
				return buf, false
			}
		case len(tail) > 8 && string(tail[:8]) == `"phone":`:
			if tail, temp, ok = parse.ParseQuoted(tail[8:]); !ok {
				return buf, false
			}
			copy(a.Phone[:], temp)
		case len(tail) > 6 && string(tail[:6]) == `"sex":`:
			if tail, a.Sex, ok = ParseSex(tail[6:]); !ok {
				return buf, false
			}
		case len(tail) > 10 && string(tail[:10]) == `"country":`:
			if tail, a.Country, ok = ParseCountry(tail[10:]); !ok {
				return buf, false
			}
		case len(tail) > 7 && string(tail[:7]) == `"city":`:
			if tail, a.City, ok = ParseCity(tail[7:]); !ok {
				return buf, false
			}
		case len(tail) > 9 && string(tail[:9]) == `"status":`:
			if tail, a.Status, ok = ParseStatus(tail[9:]); !ok {
				return buf, false
			}
		case len(tail) > 10 && string(tail[:10]) == `"premium":`:
			if tail, ok = parse.ParseSymbol(tail[10:], '{'); !ok {
				return buf, false
			}
			for {
				tail = parse.ParseSpaces(tail)
				switch {
				case len(tail) > 8 && string(tail[:8]) == `"start":`:
					if tail, temp, ok = parse.ParseNumbers(tail[8:]); !ok {
						return buf, false
					}
					copy(a.PremiumStart[:], temp)
				case len(tail) > 9 && string(tail[:9]) == `"finish":`:
					if tail, temp, ok = parse.ParseNumbers(tail[9:]); !ok {
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
		case len(tail) > 12 && string(tail[:12]) == `"interests":`:
			if tail, ok = parse.ParseSymbol(tail[12:], '['); !ok {
				return buf, false
			}
			var i uint8
			var interest uint8
			for {
				if tail, interest, ok = ParseInterest(tail); !ok {
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
		case len(tail) > 8 && string(tail[:8]) == `"likes":`:
			if tail, ok = parse.ParseSymbol(tail[8:], '['); !ok {
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
					case len(tail) > 5 && string(tail[:5]) == `"id":`:
						if tail, temp, ok = parse.ParseNumbers(tail[5:]); !ok {
							return buf, false
						}
						ID = parse.AtoiNocheck(temp)
					case len(tail) > 5 && string(tail[:5]) == `"ts":`:
						if tail, temp, ok = parse.ParseNumbers(tail[5:]); !ok {
							return buf, false
						}
						TS = parse.AtoiNocheck(temp)
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
