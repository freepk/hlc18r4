package proto

func GetFnameToken(b []byte) (int, bool) {
	return fnameDict.Token(b)
}

func GetSnameToken(b []byte) (int, bool) {
	return snameDict.Token(b)
}

func GetCountryToken(b []byte) (int, bool) {
	return countryDict.Token(b)
}

func GetCityToken(b []byte) (int, bool) {
	return cityDict.Token(b)
}

func GetInterestToken(b []byte) (int, bool) {
	return interestDict.Token(b)
}

func GetFname(t int) ([]byte, bool) {
	return fnameDict.Value(t)
}

func GetSname(t int) ([]byte, bool) {
	return snameDict.Value(t)
}

func GetCountry(t int) ([]byte, bool) {
	return countryDict.Value(t)
}

func GetCity(t int) ([]byte, bool) {
	return cityDict.Value(t)
}

func GetInterest(t int) ([]byte, bool) {
	return interestDict.Value(t)
}
