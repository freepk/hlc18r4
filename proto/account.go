package proto

import (
	"gitlab.com/freepk/hlc18r4/parse"
)

const EmailBytesMaxSize = 40

const PhoneBytesMaxSize = 12

const NumberBytesMaxSize = 10

const InterestsMaxSize = 10

type EmailBytes struct {
	Size  uint8
	Bytes [EmailBytesMaxSize]byte
}

type PhoneBytes [PhoneBytesMaxSize]byte

type NumberBytes [NumberBytesMaxSize]byte

type InterestsBytes [InterestsMaxSize]byte

type Like struct {
	ID uint32
	TS uint32
}

type Account struct {
	ID            NumberBytes
	Birth         NumberBytes
	Joined        NumberBytes
	Email         EmailBytes
	Fname         uint8
	Sname         uint16
	Phone         PhoneBytes
	Sex           SexEnum
	Country       uint8
	City          uint16
	Status        StatusEnum
	PremiumStart  NumberBytes
	PremiumFinish NumberBytes
	Interests     InterestsBytes
	LikesTo       []Like
}

func (a *Account) reset() {
	var number NumberBytes
	var email EmailBytes
	var phone PhoneBytes
	var interests InterestsBytes

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
			a.Email.Size = uint8(copy(a.Email.Bytes[:], temp))
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
			/*
				case len(tail) > PremiumLen && string(tail[:PremiumLen]) == PremiumKey:
					// ...
					if tail, ok = parse.ParseSymbol(tail[PremiumLen:], '{'); !ok {
						return buf, false
					}
					for {
						tail = parse.ParseSpaces(tail)
						switch {
						case len(tail) > StartLen && string(tail[:StartLen]) == StartKey:
							if tail, _, ok = parse.ParseInt(tail[StartLen:]); !ok {
								return buf, false
							}
						case len(tail) > FinishLen && string(tail[:FinishLen]) == FinishKey:
							if tail, _, ok = parse.ParseInt(tail[FinishLen:]); !ok {
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
					// ...
						case len(tail) > InterestsLen && string(tail[:InterestsLen]) == InterestsKey:
							if tail, ok = parse.ParseSymbol(tail[InterestsLen:], '['); !ok {
								return buf, false
							}
							for {
								var interest uint8
								if tail, interest, ok = ParseInterest(tail); !ok {
									return buf, false
								}
								a.Interests = append(a.Interests, interest)
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
								like := Like{}
								if tail, ok = parse.ParseSymbol(tail, '{'); !ok {
									return buf, false
								}
								for {
									tail = parse.ParseSpaces(tail)
									switch {
									case len(tail) > IdLen && string(tail[:IdLen]) == IdKey:
										if tail, like.ID, ok = parse.ParseUint32(tail[IdLen:]); !ok {
											return buf, false
										}
									case len(tail) > TsLen && string(tail[:TsLen]) == TsKey:
										if tail, like.TS, ok = parse.ParseUint32(tail[TsLen:]); !ok {
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
								a.LikesTo = append(a.LikesTo, like)
								if tail, ok = parse.ParseSymbol(tail, ','); !ok {
									break
								}
							}
							if tail, ok = parse.ParseSymbol(tail, ']'); !ok {
								return buf, false
							}
			*/
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
