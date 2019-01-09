package main

type Premium struct {
	Start  int
	Finish int
}

type Like struct {
	ID int
	TS int
}

type Account struct {
}

func (a *Account) Reset() {
}

func (a *Account) SetID(id int) {
}

func (a *Account) SetBirth(birth int) {
}

func (a *Account) SetJoined(joined int) {
}

func (a *Account) SetEmail(email []byte) {
}

func (a *Account) SetFname(fname []byte) {
}

func (a *Account) SetSname(sname []byte) {
}

func (a *Account) SetPhone(phone []byte) {
}

func (a *Account) SetSex(sex []byte) {
}

func (a *Account) SetCountry(country []byte) {
	countryLookup.GetIndexOrSet(country)
}

func (a *Account) SetCity(city []byte) {
}

func (a *Account) SetStatus(status []byte) {
}

func (a *Account) SetPremium(premium Premium) {
}

func (a *Account) AddInterest(interest []byte) {
}

func (a *Account) AddLike(like Like) {
}

func (a *Account) UnmarshalJSON(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	if t, ok = parseSymbol(b, '{'); !ok {
		return b, false
	}
	for {
		t = parseSpaces(t)
		switch {
		case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
			var id int
			if id, t, ok = parseInt(t[accIdLen:]); !ok {
				return b, false
			}
			a.SetID(id)
		case len(t) > accBirthLen && string(t[:accBirthLen]) == accBirthKey:
			var birth int
			if birth, t, ok = parseInt(t[accBirthLen:]); !ok {
				return b, false
			}
			a.SetBirth(birth)
		case len(t) > accJoinedLen && string(t[:accJoinedLen]) == accJoinedKey:
			var joined int
			if joined, t, ok = parseInt(t[accJoinedLen:]); !ok {
				return b, false
			}
			a.SetJoined(joined)
		case len(t) > accEmailLen && string(t[:accEmailLen]) == accEmailKey:
			var email []byte
			if email, t, ok = parsePhrase(t[accEmailLen:], '"'); !ok {
				return b, false
			}
			a.SetEmail(email)
		case len(t) > accFnameLen && string(t[:accFnameLen]) == accFnameKey:
			var fname []byte
			if fname, t, ok = parsePhrase(t[accFnameLen:], '"'); !ok {
				return b, false
			}
			a.SetFname(fname)
		case len(t) > accSnameLen && string(t[:accSnameLen]) == accSnameKey:
			var sname []byte
			if sname, t, ok = parsePhrase(t[accSnameLen:], '"'); !ok {
				return b, false
			}
			a.SetSname(sname)
		case len(t) > accPhoneLen && string(t[:accPhoneLen]) == accPhoneKey:
			var phone []byte
			if phone, t, ok = parsePhrase(t[accPhoneLen:], '"'); !ok {
				return b, false
			}
			a.SetPhone(phone)
		case len(t) > accSexLen && string(t[:accSexLen]) == accSexKey:
			var sex []byte
			if sex, t, ok = parsePhrase(t[accSexLen:], '"'); !ok {
				return b, false
			}
			a.SetSex(sex)
		case len(t) > accCountryLen && string(t[:accCountryLen]) == accCountryKey:
			var country []byte
			if country, t, ok = parsePhrase(t[accCountryLen:], '"'); !ok {
				return b, false
			}
			a.SetCountry(country)
		case len(t) > accCityLen && string(t[:accCityLen]) == accCityKey:
			var city []byte
			if city, t, ok = parsePhrase(t[accCityLen:], '"'); !ok {
				return b, false
			}
			a.SetCity(city)
		case len(t) > accStatusLen && string(t[:accStatusLen]) == accStatusKey:
			var status []byte
			if status, t, ok = parsePhrase(t[accStatusLen:], '"'); !ok {
				return b, false
			}
			a.SetStatus(status)
		case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
			var premium Premium
			if t, ok = parseSymbol(t[accPremiumLen:], '{'); !ok {
				return b, false
			}
			for {
				t = parseSpaces(t)
				switch {
				case len(t) > premiumStartLen && string(t[:premiumStartLen]) == premiumStartKey:
					if premium.Start, t, ok = parseInt(t[premiumStartLen:]); !ok {
						return b, false
					}
				case len(t) > premiumFinishLen && string(t[:premiumFinishLen]) == premiumFinishKey:
					if premium.Finish, t, ok = parseInt(t[premiumFinishLen:]); !ok {
						return b, false
					}
				}
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, '}'); !ok {
				return b, false
			}
			a.SetPremium(premium)
		case len(t) > accInterestsLen && string(t[:accInterestsLen]) == accInterestsKey:
			if t, ok = parseSymbol(t[accInterestsLen:], '['); !ok {
				return b, false
			}
			for {
				var interest []byte
				if interest, t, ok = parsePhrase(t, '"'); !ok {
					return b, false
				}
				a.AddInterest(interest)
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, ']'); !ok {
				return b, false
			}
		case len(t) > accLikesLen && string(t[:accLikesLen]) == accLikesKey:
			if t, ok = parseSymbol(t[accLikesLen:], '['); !ok {
				return b, false
			}
			for {
				var like Like
				if t, ok = parseSymbol(t, '{'); !ok {
					return b, false
				}
				for {
					t = parseSpaces(t)
					switch {
					case len(t) > likesIdLen && string(t[:likesIdLen]) == likesIdKey:
						if like.ID, t, ok = parseInt(t[likesIdLen:]); !ok {
							return b, false
						}
					case len(t) > likesTsLen && string(t[:likesTsLen]) == likesTsKey:
						if like.TS, t, ok = parseInt(t[likesTsLen:]); !ok {
							return b, false
						}
					}
					if t, ok = parseSymbol(t, ','); !ok {
						break
					}
				}
				if t, ok = parseSymbol(t, '}'); !ok {
					return b, false
				}
				a.AddLike(like)
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, ']'); !ok {
				return b, false
			}
		}
		if t, ok = parseSymbol(t, ','); !ok {
			break
		}
	}
	if t, ok = parseSymbol(t, '}'); !ok {
		return b, false
	}
	return t, true
}
