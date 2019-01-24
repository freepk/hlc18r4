package proto

func FnameToken(b []byte) (int, bool) {
	return fnameDict.Token(b)
}

func SnameToken(b []byte) (int, bool) {
	return snameDict.Token(b)
}

func CountryToken(b []byte) (int, bool) {
	return countryDict.Token(b)
}

func CityToken(b []byte) (int, bool) {
	return cityDict.Token(b)
}

func InterestToken(b []byte) (int, bool) {
	return interestDict.Token(b)
}
