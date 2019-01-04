package main

func checkByte(b []byte, c byte) ([]byte, bool) {

	// Fast way
	//if len(b) == 0 {
	//	return b, false
	//}
	//if b[0] == c {
	//	return b[1:], true
	//}

	for i, x := range b {
		if x > 0x20 {
			if b[i] == c {
				return b[i+1:], true
			}
			return b, false
		}
	}

	return b, false

}

func parseAccount(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	if t, ok = checkByte(t, '{'); !ok {
		return b, false
	}

	switch {
	case len(t) > accBirthLen && string(t[:accBirthLen]) == accBirthKey:
	case len(t) > accCityLen && string(t[:accCityLen]) == accCityKey:
	case len(t) > accCountryLen && string(t[:accCountryLen]) == accCountryKey:
	case len(t) > accEmailLen && string(t[:accEmailLen]) == accEmailKey:
	case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
	case len(t) > accInterestsLen && string(t[:accInterestsLen]) == accInterestsKey:
	case len(t) > accPhoneLen && string(t[:accPhoneLen]) == accPhoneKey:
	case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
	}

	return b, false
}
