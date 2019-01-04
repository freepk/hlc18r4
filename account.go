package main

func parseAccount(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	if t, ok = parseSymbol(b, '{'); !ok {
		return b, false
	}
	println("Account {")
	for {
		t = parseSpaces(t)
		switch {
		case len(t) > accBirthLen && string(t[:accBirthLen]) == accBirthKey:
			t = t[accBirthLen:]
			v := 0
			if v, t, ok = parseInt(t); !ok {
				return b, false
			}
			_ = v
			println("\tBirth", v)
		case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
			t = t[accIdLen:]
			v := 0
			if v, t, ok = parseInt(t); !ok {
				return b, false
			}
			_ = v
			println("\tID", v)
		case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
			t = t[accPremiumLen:]
			if t, ok = parseSymbol(t, '{'); !ok {
				return b, false
			}
			println("\tPremium {")
			for {
				t = parseSpaces(t)
				switch {
				case len(t) > premiumStartLen && string(t[:premiumStartLen]) == premiumStartKey:
					t = t[premiumStartLen:]
					v := 0
					if v, t, ok = parseInt(t); !ok {
						return b, false
					}
					_ = v
					println("\t\tStart", v)
				case len(t) > premiumFinishLen && string(t[:premiumFinishLen]) == premiumFinishKey:
					t = t[premiumFinishLen:]
					v := 0
					if v, t, ok = parseInt(t); !ok {
						return b, false
					}
					_ = v
					println("\t\tFinish", v)
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
			t = t[accJoinedLen:]
			v := 0
			if v, t, ok = parseInt(t); !ok {
				return b, false
			}
			_ = v
			println("\tJoined", v)
		case len(t) > accLikesLen && string(t[:accLikesLen]) == accLikesKey:
			t = t[accLikesLen:]
			if t, ok = parseSymbol(t, '['); !ok {
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
						t = t[likesIdLen:]
						v := 0
						if v, t, ok = parseInt(t); !ok {
							return b, false
						}
						_ = v
						println("\t\t\tId", v)
					case len(t) > likesTsLen && string(t[:likesTsLen]) == likesTsKey:
						t = t[likesTsLen:]
						v := 0
						if v, t, ok = parseInt(t); !ok {
							return b, false
						}
						_ = v
						println("\t\t\tTs", v)
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
			//case len(t) > accCityLen && string(t[:accCityLen]) == accCityKey:
			//case len(t) > accCountryLen && string(t[:accCountryLen]) == accCountryKey:
			//case len(t) > accEmailLen && string(t[:accEmailLen]) == accEmailKey:
			//case len(t) > accInterestsLen && string(t[:accInterestsLen]) == accInterestsKey:
			//case len(t) > accPhoneLen && string(t[:accPhoneLen]) == accPhoneKey:
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

func parseAccountX(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	if t, ok = parseSymbol(b, '{'); !ok {
		return b, false
	}
	for {
		t = parseSpaces(t)
		switch {
		case len(t) > accBirthLen && string(t[:accBirthLen]) == accBirthKey:
			t = t[accBirthLen:]
			if _, t, ok = parseInt(t); !ok {
				return b, false
			}
		case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
			t = t[accIdLen:]
			if _, t, ok = parseInt(t); !ok {
				return b, false
			}
		case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
			t = t[accPremiumLen:]
			if t, ok = parseSymbol(t, '{'); !ok {
				return b, false
			}
			for {
				t = parseSpaces(t)
				switch {
				case len(t) > premiumStartLen && string(t[:premiumStartLen]) == premiumStartKey:
					t = t[premiumStartLen:]
					if _, t, ok = parseInt(t); !ok {
						return b, false
					}
				case len(t) > premiumFinishLen && string(t[:premiumFinishLen]) == premiumFinishKey:
					t = t[premiumFinishLen:]
					if _, t, ok = parseInt(t); !ok {
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
		case len(t) > accJoinedLen && string(t[:accJoinedLen]) == accJoinedKey:
			t = t[accJoinedLen:]
			if _, t, ok = parseInt(t); !ok {
				return b, false
			}
		case len(t) > accLikesLen && string(t[:accLikesLen]) == accLikesKey:
			t = t[accLikesLen:]
			if t, ok = parseSymbol(t, '['); !ok {
				return b, false
			}
			for {
				if t, ok = parseSymbol(t, '{'); !ok {
					return b, false
				}
				for {
					t = parseSpaces(t)
					switch {
					case len(t) > likesIdLen && string(t[:likesIdLen]) == likesIdKey:
						t = t[likesIdLen:]
						if _, t, ok = parseInt(t); !ok {
							return b, false
						}
					case len(t) > likesTsLen && string(t[:likesTsLen]) == likesTsKey:
						t = t[likesTsLen:]
						if _, t, ok = parseInt(t); !ok {
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
	return b, false
}
