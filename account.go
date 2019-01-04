package main

func checkByte(b []byte, c byte) ([]byte, bool) {

	for i, x := range b {
		if x > 0x20 {
			if b[i] == c {
				return b[i+1:], true
			}
			return b, false
		}
	}

	//if len(b) == 0 {
	//	return b, false
	//}
	//if b[0] == c {
	//	return b[1:], true
	//}

	return b, false

}

const (
	accIdKey        = `"id"`
	accIdLen        = len(accIdKey)
	accInterestsKey = `"interests"`
	accInterestsLen = len(accInterestsKey)
	accPhoneKey     = `"phone"`
	accPhoneLen     = len(accPhoneKey)
	accPremiumKey   = `"premium"`
	accPremiumLen   = len(accPremiumKey)
)

func parseAccount(b []byte) ([]byte, bool) {
	var t []byte
	var ok bool

	t, ok = checkByte(b, '{')
	if !ok {
		return b, false
	}

	// Validate key
	//      i       id interests
	//      e       email
	//      f       finish fname
	//      p       phone premium
	//      b       birth
	//      c       city country
	//      j       joined
	//      s       sex sname status
	//      l       likes

	// Minimal value len == len(`"_":_`)
	if len(t) < 5 {
		return b, false
	}

	switch {
	case len(t) > accIdLen && string(t[:accIdLen]) == accIdKey:
		//println("\n\nValid id key")
	case len(t) > accInterestsLen && string(t[:accInterestsLen]) == accInterestsKey:
		//println("\n\nValid interests key")
	case len(t) > accPhoneLen && string(t[:accPhoneLen]) == accPhoneKey:
		//println("\n\nValid phone key")
	case len(t) > accPremiumLen && string(t[:accPremiumLen]) == accPremiumKey:
		//println("\n\nValid premium key")
	}

	return b, false
}
