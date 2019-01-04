package main

func parseAccount(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	if t, ok = parseSymbol(b, '{'); !ok {
		return b, false
	}

	for {

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

		break
	}

	if t, ok = parseSymbol(t, '}'); !ok {
		return b, false
	}

	return b, false
}
