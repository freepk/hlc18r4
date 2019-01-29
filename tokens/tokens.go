package tokens

import (
	"github.com/freepk/dictionary"
)

var (
	fnameDict    = dictionary.NewDictionary(4)
	snameDict    = dictionary.NewDictionary(4)
	countryDict  = dictionary.NewDictionary(4)
	cityDict     = dictionary.NewDictionary(4)
	interestDict = dictionary.NewDictionary(4)
)

func AddFname(b []byte) int {
	k, _ := fnameDict.AddKey(b)
	return k
}

func AddSname(b []byte) int {
	k, _ := snameDict.AddKey(b)
	return k
}

func AddCountry(b []byte) int {
	k, _ := countryDict.AddKey(b)
	return k
}

func AddCity(b []byte) int {
	k, _ := cityDict.AddKey(b)
	return k
}

func AddInterest(b []byte) int {
	k, _ := interestDict.AddKey(b)
	return k
}

func FnameVal(k int) ([]byte, bool) {
	return fnameDict.Val(k)
}

func SnameVal(k int) ([]byte, bool) {
	return snameDict.Val(k)
}

func CountryVal(k int) ([]byte, bool) {
	return countryDict.Val(k)
}

func CityVal(k int) ([]byte, bool) {
	return cityDict.Val(k)
}

func InterestVal(k int) ([]byte, bool) {
	return interestDict.Val(k)
}
