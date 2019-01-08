package main

//import (
//	"fmt"
//)

type Like struct {
	ID int
	TS int
}

type Account struct {
	ID            int
	Birth         int
	Joined        int
	Email         []byte
	Fname         []byte
	Sname         []byte
	Phone         []byte
	Sex           []byte
	Country       []byte
	City          []byte
	Status        []byte
	PremiumStart  int
	PremiumFinish int
	Interests     [][]byte
	Likes         []Like
}

func (acc *Account) Reset() {
	acc.ID = 0
	acc.Birth = 0
	acc.Joined = 0
	acc.Email = acc.Email[:0]
	acc.Fname = acc.Fname[:0]
	acc.Sname = acc.Sname[:0]
	acc.Phone = acc.Phone[:0]
	acc.Sex = acc.Sex[:0]
	acc.Country = acc.Country[:0]
	acc.City = acc.City[:0]
	acc.Status = acc.Status[:0]
	acc.PremiumStart = 0
	acc.PremiumFinish = 0
	acc.Interests = acc.Interests[:0]
	acc.Likes = acc.Likes[:0]
}

func (acc *Account) Tokens() []int {
	return nil
}

func (acc *Account) Parse(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	acc.Reset()

	if t, ok = parseSymbol(b, '{'); !ok {
		return b, false
	}
	// println("Account {")
	for {
		t = parseSpaces(t)
		switch {

		case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
			if acc.ID, t, ok = parseInt(t[accIdLen:]); !ok {
				return b, false
			}
			// println("\tID", acc.ID)
		case len(t) > accBirthLen && string(t[:accBirthLen]) == accBirthKey:
			if acc.Birth, t, ok = parseInt(t[accBirthLen:]); !ok {
				return b, false
			}
			// println("\tBirth", acc.Birth)
		case len(t) > accJoinedLen && string(t[:accJoinedLen]) == accJoinedKey:
			if acc.Joined, t, ok = parseInt(t[accJoinedLen:]); !ok {
				return b, false
			}
			// println("\tJoined", acc.Joined)
		case len(t) > accEmailLen && string(t[:accEmailLen]) == accEmailKey:
			if acc.Email, t, ok = parsePhrase(t[accEmailLen:], '"'); !ok {
				return b, false
			}
			// println("\tEmail", string(acc.Email))
		case len(t) > accFnameLen && string(t[:accFnameLen]) == accFnameKey:
			if acc.Fname, t, ok = parsePhrase(t[accFnameLen:], '"'); !ok {
				return b, false
			}
			// println("\tFname", string(acc.Fname))
		case len(t) > accSnameLen && string(t[:accSnameLen]) == accSnameKey:
			if acc.Sname, t, ok = parsePhrase(t[accSnameLen:], '"'); !ok {
				return b, false
			}
			// println("\tSname", string(acc.Sname))
		case len(t) > accPhoneLen && string(t[:accPhoneLen]) == accPhoneKey:
			if acc.Phone, t, ok = parsePhrase(t[accPhoneLen:], '"'); !ok {
				return b, false
			}
			// println("\tPhone", string(acc.Phone))
		case len(t) > accSexLen && string(t[:accSexLen]) == accSexKey:
			if acc.Sex, t, ok = parsePhrase(t[accSexLen:], '"'); !ok {
				return b, false
			}
			// println("\tSex", string(acc.Sex))
		case len(t) > accCountryLen && string(t[:accCountryLen]) == accCountryKey:
			if acc.Country, t, ok = parsePhrase(t[accCountryLen:], '"'); !ok {
				return b, false
			}
			// println("\tCountry", string(acc.Country))
		case len(t) > accCityLen && string(t[:accCityLen]) == accCityKey:
			if acc.City, t, ok = parsePhrase(t[accCityLen:], '"'); !ok {
				return b, false
			}
			// println("\tCity", string(acc.City))
		case len(t) > accStatusLen && string(t[:accStatusLen]) == accStatusKey:
			if acc.Status, t, ok = parsePhrase(t[accStatusLen:], '"'); !ok {
				return b, false
			}
			// println("\tStatus", string(acc.Status))
		case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
			if t, ok = parseSymbol(t[accPremiumLen:], '{'); !ok {
				return b, false
			}
			// println("\tPremium {")
			for {
				t = parseSpaces(t)
				switch {
				case len(t) > premiumStartLen && string(t[:premiumStartLen]) == premiumStartKey:
					if acc.PremiumStart, t, ok = parseInt(t[premiumStartLen:]); !ok {
						return b, false
					}
					// println("\t\tStart", acc.PremiumStart)
				case len(t) > premiumFinishLen && string(t[:premiumFinishLen]) == premiumFinishKey:
					if acc.PremiumFinish, t, ok = parseInt(t[premiumFinishLen:]); !ok {
						return b, false
					}
					// println("\t\tFinish", acc.PremiumFinish)
				}
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, '}'); !ok {
				return b, false
			}
			// println("\t} // Premium")
		case len(t) > accInterestsLen && string(t[:accInterestsLen]) == accInterestsKey:
			if t, ok = parseSymbol(t[accInterestsLen:], '['); !ok {
				return b, false
			}
			// println("\tInterests [")
			for {
				var accInterest []byte
				if accInterest, t, ok = parsePhrase(t, '"'); !ok {
					return b, false
				}
				// println("\t\t", string(accInterest))
				acc.Interests = append(acc.Interests, accInterest)
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, ']'); !ok {
				return b, false
			}
			// println("\t] // Interests")
		case len(t) > accLikesLen && string(t[:accLikesLen]) == accLikesKey:
			if t, ok = parseSymbol(t[accLikesLen:], '['); !ok {
				return b, false
			}
			// println("\tLikes [")
			for {
				if t, ok = parseSymbol(t, '{'); !ok {
					return b, false
				}
				// println("\t\tLike {")
				accLikeId := 0
				accLikeTs := 0
				for {
					t = parseSpaces(t)
					switch {
					case len(t) > likesIdLen && string(t[:likesIdLen]) == likesIdKey:
						if accLikeId, t, ok = parseInt(t[likesIdLen:]); !ok {
							return b, false
						}
						// println("\t\t\tId", accLikeId)
					case len(t) > likesTsLen && string(t[:likesTsLen]) == likesTsKey:
						if accLikeTs, t, ok = parseInt(t[likesTsLen:]); !ok {
							return b, false
						}
						// println("\t\t\tTs", accLikeTs)
					}
					if t, ok = parseSymbol(t, ','); !ok {
						break
					}
				}
				if t, ok = parseSymbol(t, '}'); !ok {
					return b, false
				}
				acc.Likes = append(acc.Likes, Like{ID: accLikeId, TS: accLikeTs})
				// println("\t\t} //Like")
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, ']'); !ok {
				return b, false
			}
			// println("\t] // Likes")
		}
		if t, ok = parseSymbol(t, ','); !ok {
			break
		}
	}
	if t, ok = parseSymbol(t, '}'); !ok {
		return b, false
	}
	// println("} //Account")
	//fmt.Println(acc.ID, acc.Birth, acc.Joined)
	return t, true
}
