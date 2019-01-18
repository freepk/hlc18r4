package proto

import (
	"testing"
)

var (
	testAccount = []byte(`{"id":1300001,"status":"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b","sname":"\u041b\u0435\u0431\u0435\u0442\u0430\u0442\u0435\u0432",
		"joined":1327104000,"birth":569642191,"fname":"\u0424\u0451\u0434\u043e\u0440","phone":"8(944)2990268","sex":"m","country":"\u041c\u0430\u043b\u0438\u0437\u0438\u044f",
		"city":"\u0420\u043e\u0441\u043e\u0440\u0438\u0436","email":"osrortir@yahoo.com"}`)
)

func TestAccountUnmarshalJSON(t *testing.T) {
	a := &Account{}
	_, ok := a.UnmarshalJSON(testAccount)
	if !ok {
		t.Fail()
	}

	fname, _ := FnameDict.Value(uint64(a.Fname))
	sname, _ := SnameDict.Value(uint64(a.Sname))
	country, _ := CountryDict.Value(uint64(a.Country))
	city, _ := CityDict.Value(uint64(a.City))

	status := ""
	switch a.Status {
	case FreeStatus:
		status = "Свободны"
	case BusyStatus:
		status = "Заняты"
	case ComplicatedStatus:
		status = "Все сложно"
	}

	sex := ""
	switch a.Sex {
	case MaleSex:
		status = "m"
	case FemaleSex:
		status = "f"
	}

	t.Log(
		string(a.ID[:]),
		string(a.Joined[:]),
		string(a.Birth[:]),
		string(fname),
		string(sname),
		string(country),
		string(city),
		status,
		sex,
		string(a.Email.Bytes[:a.Email.Size]),
		string(a.Phone[:]),
		string(a.PremiumStart[:]),
		string(a.PremiumFinish[:]))
}

func BenchmarkAccountReset(b *testing.B) {
	a := &Account{}
	for i := 0; i < b.N; i++ {
		a.reset()
	}
}

func BenchmarkAccountUnmarshalJSON(b *testing.B) {
	a := &Account{}
	for i := 0; i < b.N; i++ {
		_, ok := a.UnmarshalJSON(testAccount)
		if !ok {
			b.Fatal()
		}
	}
}

func BenchmarkAccountMarshalJSON(b *testing.B) {
	a := &Account{}
	_, ok := a.UnmarshalJSON(testAccount)
	if !ok {
		b.Fatal("UnmarshalJSON error")
	}
	for i := 0; i < b.N; i++ {
		_, _ = FnameDict.Value(uint64(a.Fname))
		_, _ = SnameDict.Value(uint64(a.Sname))
		_, _ = CountryDict.Value(uint64(a.Country))
		_, _ = CityDict.Value(uint64(a.City))
		switch a.Status {
		case FreeStatus:
			_ = "Свободны"
		case BusyStatus:
			_ = "Заняты"
		case ComplicatedStatus:
			_ = "Все сложно"
		}
		switch a.Sex {
		case MaleSex:
			_ = "m"
		case FemaleSex:
			_ = "f"
		}
		_ = a.ID[:]
		_ = a.Joined[:]
		_ = a.Birth[:]
		_ = a.Email.Bytes[:a.Email.Size]
		_ = a.Phone[:]
		_ = a.PremiumStart[:]
		_ = a.PremiumFinish[:]
	}
}
