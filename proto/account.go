package proto

import (
	"sync"

	"github.com/freepk/parse"
	"gitlab.com/freepk/hlc18r4/tokens"
)

type buffer struct {
	B []byte
}

var bufferPool = &sync.Pool{New: func() interface{} { return new(buffer) }}

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
	BirthTS       uint32
	Joined        NumberBuf
	JoinedTS      uint32
	Email         EmailBuf
	Fname         uint8
	Sname         uint16
	Phone         PhoneBuf
	Sex           uint8
	Country       uint8
	City          uint16
	Status        uint8
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
	a.BirthTS = 0
	a.Joined = number
	a.JoinedTS = 0
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

type extendedWriter interface {
	Write(p []byte) (int, error)
	WriteString(s string) (int, error)
}

func (a *Account) WriteJSON(fields int, w extendedWriter) {
	w.WriteString(`{"id":`)
	w.Write(trim(a.ID[:]))
	if (fields & BirthField) == BirthField {
		w.WriteString(`,"birth":`)
		w.Write(trim(a.Birth[:]))
	}
	if (fields & JoinedField) == JoinedField {
		w.WriteString(`,"joined":`)
		w.Write(trim(a.Joined[:]))
	}
	if (fields & EmailField) == EmailField {
		w.WriteString(`,"email":"`)
		w.Write(a.Email.Buf[:a.Email.Len])
		w.WriteString(`"`)
	}
	if (fields&FnameField) == FnameField && a.Fname > 0 {
		fname, _ := tokens.FnameVal(int(a.Fname))
		w.WriteString(`,"fname":"`)
		w.Write(fname)
		w.WriteString(`"`)
	}
	if (fields&SnameField) == SnameField && a.Sname > 0 {
		sname, _ := tokens.SnameVal(int(a.Sname))
		w.WriteString(`,"sname":"`)
		w.Write(sname)
		w.WriteString(`"`)
	}
	if (fields&PhoneField) == PhoneField && a.Phone[0] > 0 {
		w.WriteString(`,"phone":"`)
		w.Write(trim(a.Phone[:]))
		w.WriteString(`"`)
	}
	if (fields & SexField) == SexField {
		switch a.Sex {
		case tokens.MaleSex:
			w.WriteString(`,"sex":"m"`)
		case tokens.FemaleSex:
			w.WriteString(`,"sex":"f"`)
		}
	}
	if (fields&CountryField) == CountryField && a.Country > 0 {
		country, _ := tokens.CountryVal(int(a.Country))
		w.WriteString(`,"country":"`)
		w.Write(country)
		w.WriteString(`"`)
	}
	if (fields&CityField) == CityField && a.City > 0 {
		city, _ := tokens.CityVal(int(a.City))
		w.WriteString(`,"city":"`)
		w.Write(city)
		w.WriteString(`"`)
	}
	if (fields & StatusField) == StatusField {
		switch a.Status {
		case tokens.SingleStatus:
			w.WriteString(`,"status":"свободны"`)
		case tokens.InRelStatus:
			w.WriteString(`,"status":"заняты"`)
		case tokens.ComplStatus:
			w.WriteString(`,"status":"всё сложно"`)
		}
	}
	if (fields&PremiumField) == PremiumField && a.PremiumFinish[0] > 0 {
		w.WriteString(`,"premium":{"finish":`)
		w.Write(trim(a.PremiumFinish[:]))
		w.WriteString(`,"start":`)
		w.Write(trim(a.PremiumStart[:]))
		w.WriteString(`}`)
	}
	w.WriteString(`}`)
}

func (a *Account) UnmarshalJSON(buf []byte) ([]byte, bool) {
	var tail []byte
	var temp []byte
	var ok bool

	a.reset()

	enc := bufferPool.Get().(*buffer)
	defer bufferPool.Put(enc)

	if tail, ok = parse.SkipSymbol(buf, '{'); !ok {
		return buf, false
	}
	for {
		tail = parse.SkipSpaces(tail)
		switch {
		case len(tail) > 5 && string(tail[:5]) == `"id":`:
			if tail, temp, ok = parse.ParseNumber(tail[5:]); !ok {
				return buf, false
			}
			copy(a.ID[:], temp)
		case len(tail) > 8 && string(tail[:8]) == `"birth":`:
			if tail, temp, ok = parse.ParseNumber(tail[8:]); !ok {
				return buf, false
			}
			copy(a.Birth[:], temp)
			_, a.BirthTS, ok = parse.ParseUint32(temp)
		case len(tail) > 9 && string(tail[:9]) == `"joined":`:
			if tail, temp, ok = parse.ParseNumber(tail[9:]); !ok {
				return buf, false
			}
			copy(a.Joined[:], temp)
			_, a.JoinedTS, ok = parse.ParseUint32(temp)
		case len(tail) > 8 && string(tail[:8]) == `"email":`:
			if tail, temp, ok = parse.ParseQuoted(tail[8:]); !ok {
				return buf, false
			}
			a.Email.Len = uint8(copy(a.Email.Buf[:], temp))
		case len(tail) > 8 && string(tail[:8]) == `"fname":`:
			if tail, a.Fname, ok = parseFname(tail[8:], enc); !ok {
				return buf, false
			}
		case len(tail) > 8 && string(tail[:8]) == `"sname":`:
			if tail, a.Sname, ok = parseSname(tail[8:], enc); !ok {
				return buf, false
			}
		case len(tail) > 8 && string(tail[:8]) == `"phone":`:
			if tail, temp, ok = parse.ParseQuoted(tail[8:]); !ok {
				return buf, false
			}
			copy(a.Phone[:], temp)
		case len(tail) > 6 && string(tail[:6]) == `"sex":`:
			if tail, a.Sex, ok = parseSex(tail[6:]); !ok {
				return buf, false
			}
		case len(tail) > 10 && string(tail[:10]) == `"country":`:
			if tail, a.Country, ok = parseCountry(tail[10:], enc); !ok {
				return buf, false
			}
		case len(tail) > 7 && string(tail[:7]) == `"city":`:
			if tail, a.City, ok = parseCity(tail[7:], enc); !ok {
				return buf, false
			}
		case len(tail) > 9 && string(tail[:9]) == `"status":`:
			if tail, a.Status, ok = parseStatus(tail[9:]); !ok {
				return buf, false
			}
		case len(tail) > 10 && string(tail[:10]) == `"premium":`:
			if tail, ok = parse.SkipSymbol(tail[10:], '{'); !ok {
				return buf, false
			}
			for {
				tail = parse.SkipSpaces(tail)
				switch {
				case len(tail) > 8 && string(tail[:8]) == `"start":`:
					if tail, temp, ok = parse.ParseNumber(tail[8:]); !ok {
						return buf, false
					}
					copy(a.PremiumStart[:], temp)
				case len(tail) > 9 && string(tail[:9]) == `"finish":`:
					if tail, temp, ok = parse.ParseNumber(tail[9:]); !ok {
						return buf, false
					}
					copy(a.PremiumFinish[:], temp)
				}
				if tail, ok = parse.SkipSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.SkipSymbol(tail, '}'); !ok {
				return buf, false
			}
		case len(tail) > 12 && string(tail[:12]) == `"interests":`:
			if tail, ok = parse.SkipSymbol(tail[12:], '['); !ok {
				return buf, false
			}
			var i uint8
			var interest uint8
			for {
				if tail, interest, ok = parseInterest(tail, enc); !ok {
					return buf, false
				}
				a.Interests[i] = interest
				i++
				if tail, ok = parse.SkipSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.SkipSymbol(tail, ']'); !ok {
				return buf, false
			}
		case len(tail) > 8 && string(tail[:8]) == `"likes":`:
			if tail, ok = parse.SkipSymbol(tail[8:], '['); !ok {
				return buf, false
			}
			for {
				id := 0
				ts := 0
				if tail, ok = parse.SkipSymbol(tail, '{'); !ok {
					return buf, false
				}
				for {
					tail = parse.SkipSpaces(tail)
					switch {
					case len(tail) > 5 && string(tail[:5]) == `"id":`:
						if tail, id, ok = parse.ParseInt(tail[5:]); !ok {
							return buf, false
						}
					case len(tail) > 5 && string(tail[:5]) == `"ts":`:
						if tail, ts, ok = parse.ParseInt(tail[5:]); !ok {
							return buf, false
						}
					}
					if tail, ok = parse.SkipSymbol(tail, ','); !ok {
						break
					}
				}
				if tail, ok = parse.SkipSymbol(tail, '}'); !ok {
					return buf, false
				}
				a.LikesTo = append(a.LikesTo, Like{ID: uint32(id), TS: uint32(ts)})
				if tail, ok = parse.SkipSymbol(tail, ','); !ok {
					break
				}
			}
			if tail, ok = parse.SkipSymbol(tail, ']'); !ok {
				return buf, false
			}
		}
		if tail, ok = parse.SkipSymbol(tail, ','); !ok {
			break
		}
	}
	if tail, ok = parse.SkipSymbol(tail, '}'); !ok {
		return buf, false
	}
	return tail, true
}
