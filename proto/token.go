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
