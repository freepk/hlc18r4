package main

func parseAccount(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool
	var vint int
	var vbytes []byte

	if t, ok = parseSymbol(b, '{'); !ok {
		return b, false
	}
	println("Account {")
	for {
		t = parseSpaces(t)
		switch {
		case len(t) > accSexLen && string(t[:accSexLen]) == accSexKey:
			if vbytes, t, ok = parsePhrase(t[accSexLen:], '"'); !ok {
				return b, false
			}
			println("\tSex", string(vbytes))
		case len(t) > accEmailLen && string(t[:accEmailLen]) == accEmailKey:
			if vbytes, t, ok = parsePhrase(t[accEmailLen:], '"'); !ok {
				return b, false
			}
			println("\tEmail", string(vbytes))
		case len(t) > accSnameLen && string(t[:accSnameLen]) == accSnameKey:
			if vbytes, t, ok = parsePhrase(t[accSnameLen:], '"'); !ok {
				return b, false
			}
			println("\tSname", string(vbytes))
		case len(t) > accFnameLen && string(t[:accFnameLen]) == accFnameKey:
			if vbytes, t, ok = parsePhrase(t[accFnameLen:], '"'); !ok {
				return b, false
			}
			println("\tFname", string(vbytes))
		case len(t) > accStatusLen && string(t[:accStatusLen]) == accStatusKey:
			if vbytes, t, ok = parsePhrase(t[accStatusLen:], '"'); !ok {
				return b, false
			}
			println("\tStatus", string(vbytes))
		case len(t) > accCityLen && string(t[:accCityLen]) == accCityKey:
			if vbytes, t, ok = parsePhrase(t[accCityLen:], '"'); !ok {
				return b, false
			}
			println("\tCity", string(vbytes))
		case len(t) > accCountryLen && string(t[:accCountryLen]) == accCountryKey:
			if vbytes, t, ok = parsePhrase(t[accCountryLen:], '"'); !ok {
				return b, false
			}
			println("\tCountry", string(vbytes))
		case len(t) > accBirthLen && string(t[:accBirthLen]) == accBirthKey:
			if vint, t, ok = parseInt(t[accBirthLen:]); !ok {
				return b, false
			}
			println("\tBirth", vint)
		case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
			if vint, t, ok = parseInt(t[accIdLen:]); !ok {
				return b, false
			}
			println("\tID", vint)
		case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
			if t, ok = parseSymbol(t[accPremiumLen:], '{'); !ok {
				return b, false
			}
			println("\tPremium {")
			for {
				t = parseSpaces(t)
				switch {
				case len(t) > premiumStartLen && string(t[:premiumStartLen]) == premiumStartKey:
					if vint, t, ok = parseInt(t[premiumStartLen:]); !ok {
						return b, false
					}
					println("\t\tStart", vint)
				case len(t) > premiumFinishLen && string(t[:premiumFinishLen]) == premiumFinishKey:
					if vint, t, ok = parseInt(t[premiumFinishLen:]); !ok {
						return b, false
					}
					println("\t\tFinish", vint)
				}
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, '}'); !ok {
				return b, false
			}
			println("\t} // Premium")
		case len(t) > accJoinedLen && string(t[:accJoinedLen]) == accJoinedKey:
			if vint, t, ok = parseInt(t[accJoinedLen:]); !ok {
				return b, false
			}
			println("\tJoined", vint)
		case len(t) > accLikesLen && string(t[:accLikesLen]) == accLikesKey:
			if t, ok = parseSymbol(t[accLikesLen:], '['); !ok {
				return b, false
			}
			println("\tLikes [")
			for {
				if t, ok = parseSymbol(t, '{'); !ok {
					return b, false
				}
				println("\t\tLike {")
				for {
					t = parseSpaces(t)
					switch {
					case len(t) > likesIdLen && string(t[:likesIdLen]) == likesIdKey:
						if vint, t, ok = parseInt(t[likesIdLen:]); !ok {
							return b, false
						}
						println("\t\t\tId", vint)
					case len(t) > likesTsLen && string(t[:likesTsLen]) == likesTsKey:
						if vint, t, ok = parseInt(t[likesTsLen:]); !ok {
							return b, false
						}
						println("\t\t\tTs", vint)
					}
					if t, ok = parseSymbol(t, ','); !ok {
						break
					}
				}
				if t, ok = parseSymbol(t, '}'); !ok {
					return b, false
				}
				println("\t\t} //Like")
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, ']'); !ok {
				return b, false
			}
			println("\t] // Likes")
		case len(t) > accInterestsLen && string(t[:accInterestsLen]) == accInterestsKey:
			if t, ok = parseSymbol(t[accInterestsLen:], '['); !ok {
				return b, false
			}
			println("\tInterests [")
			for {
				if vbytes, t, ok = parsePhrase(t, '"'); !ok {
					return b, false
				}
				println("\t\t", string(vbytes))
				if t, ok = parseSymbol(t, ','); !ok {
					break
				}
			}
			if t, ok = parseSymbol(t, ']'); !ok {
				return b, false
			}
			println("\t] // Interests")
		}
		if t, ok = parseSymbol(t, ','); !ok {
			break
		}
	}
	if t, ok = parseSymbol(t, '}'); !ok {
		return b, false
	}
	println("} //Account")
	return b, false
}
