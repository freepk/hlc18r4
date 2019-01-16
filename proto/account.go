package proto

import (
	"gitlab.com/freepk/hlc18r4/parse"
)

type Account struct {
	ID        int
	Birth     int
	Joined    int
	Email     []byte
	Fname     []byte
	Sname     []byte
	Phone     []byte
	Sex       []byte
	Country   []byte
	City      []byte
	Status    []byte
	Interests [][]byte
	Premium   struct {
		Start  int
		Finish int
	}
	Likes []struct {
		ID int
		TS int
	}
}

func (a *Account) Reset() {
	a.ID = 0
	a.Birth = 0
	a.Joined = 0
	a.Email = a.Email[:0]
	a.Fname = a.Fname[:0]
	a.Sname = a.Sname[:0]
	a.Phone = a.Phone[:0]
	a.Sex = a.Sex[:0]
	a.Country = a.Country[:0]
	a.City = a.City[:0]
	a.Status = a.Status[:0]
	a.Premium.Start = 0
	a.Premium.Finish = 0
	a.Interests = a.Interests[:0]
	a.Likes = a.Likes[:0]
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
			if tail, a.ID, ok = parse.ParseInt(tail[IdLen:]); !ok {
				return buf, false
			}
		case len(tail) > BirthLen && string(tail[:BirthLen]) == BirthKey:
			if tail, a.Birth, ok = parse.ParseInt(tail[BirthLen:]); !ok {
				return buf, false
			}
		case len(tail) > JoinedLen && string(tail[:JoinedLen]) == JoinedKey:
			if tail, a.Joined, ok = parse.ParseInt(tail[JoinedLen:]); !ok {
				return buf, false
			}
		case len(tail) > EmailLen && string(tail[:EmailLen]) == EmailKey:
			if tail, a.Email, ok = parse.ParseQuoted(tail[EmailLen:]); !ok {
				return buf, false
			}
		case len(tail) > FnameLen && string(tail[:FnameLen]) == FnameKey:
			if tail, a.Fname, ok = parse.ParseQuoted(tail[FnameLen:]); !ok {
				return buf, false
			}
		case len(tail) > SnameLen && string(tail[:SnameLen]) == SnameKey:
			if tail, a.Sname, ok = parse.ParseQuoted(tail[SnameLen:]); !ok {
				return buf, false
			}
		case len(tail) > PhoneLen && string(tail[:PhoneLen]) == PhoneKey:
			if tail, a.Phone, ok = parse.ParseQuoted(tail[PhoneLen:]); !ok {
				return buf, false
			}
		case len(tail) > SexLen && string(tail[:SexLen]) == SexKey:
			if tail, a.Sex, ok = parse.ParseQuoted(tail[SexLen:]); !ok {
				return buf, false
			}
		case len(tail) > CountryLen && string(tail[:CountryLen]) == CountryKey:
			if tail, a.Country, ok = parse.ParseQuoted(tail[CountryLen:]); !ok {
				return buf, false
			}
		case len(tail) > CityLen && string(tail[:CityLen]) == CityKey:
			if tail, a.City, ok = parse.ParseQuoted(tail[CityLen:]); !ok {
				return buf, false
			}
		case len(tail) > StatusLen && string(tail[:StatusLen]) == StatusKey:
			if tail, a.Status, ok = parse.ParseQuoted(tail[StatusLen:]); !ok {
				return buf, false
			}
		case len(tail) > PremiumLen && string(tail[:PremiumLen]) == PremiumKey:
			var premium struct {
				Start  int
				Finish int
			}
			if tail, ok = parse.ParseSymbol(tail[PremiumLen:], '{'); !ok {
				return buf, false
			}
			for {
				tail = parse.ParseSpaces(tail)
				switch {
				case len(tail) > StartLen && string(tail[:StartLen]) == StartKey:
					if tail, premium.Start, ok = parse.ParseInt(tail[StartLen:]); !ok {
						return buf, false
					}
				case len(tail) > FinishLen && string(tail[:FinishLen]) == FinishKey:
					if tail, premium.Finish, ok = parse.ParseInt(tail[FinishLen:]); !ok {
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
			a.Premium = premium
		case len(tail) > InterestsLen && string(tail[:InterestsLen]) == InterestsKey:
			if tail, ok = parse.ParseSymbol(tail[InterestsLen:], '['); !ok {
				return buf, false
			}
			for {
				var interest []byte
				if tail, interest, ok = parse.ParseQuoted(tail); !ok {
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
				var like struct {
					ID int
					TS int
				}
				if tail, ok = parse.ParseSymbol(tail, '{'); !ok {
					return buf, false
				}
				for {
					tail = parse.ParseSpaces(tail)
					switch {
					case len(tail) > IdLen && string(tail[:IdLen]) == IdKey:
						if tail, like.ID, ok = parse.ParseInt(tail[IdLen:]); !ok {
							return buf, false
						}
					case len(tail) > TsLen && string(tail[:TsLen]) == TsKey:
						if tail, like.TS, ok = parse.ParseInt(tail[TsLen:]); !ok {
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
				a.Likes = append(a.Likes, like)
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
