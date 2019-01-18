package proto

import (
	"testing"
)

var (
	testAccount = []byte(`{"id":1300001,"status":"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b","sname":"\u041b\u0435\u0431\u0435\u0442\u0430\u0442\u0435\u0432",
		"joined":1327104000,"birth":569642191,"fname":"\u0424\u0451\u0434\u043e\u0440","phone":"8(944)2990268","sex":"m","country":"\u041c\u0430\u043b\u0438\u0437\u0438\u044f",
		"city":"\u0420\u043e\u0441\u043e\u0440\u0438\u0436","email":"osrortir@yahoo.com","interests":["\u0421\u043e\u043b\u043d\u0446\u0435",
		"\u041a\u0440\u0430\u0441\u043d\u043e\u0435 \u0432\u0438\u043d\u043e","\u041d\u043e\u0432\u044b\u0435 \u043c\u0435\u0441\u0442\u0430",
		"\u0415\u0434\u0430 \u0438 \u041d\u0430\u043f\u0438\u0442\u043a\u0438"],"likes":[{"id":607754,"ts":1475172998},{"id":912716,"ts":1523114133},
		{"id":712878,"ts":1485002960},{"id":1183626,"ts":1476080693},{"id":1043570,"ts":1505703466},{"id":295370,"ts":1513827086},{"id":279466,"ts":1452645536},
		{"id":122860,"ts":1519377612},{"id":1281780,"ts":1503744640},{"id":545550,"ts":1505010457},{"id":275498,"ts":1458432926},{"id":990978,"ts":1513836768},
		{"id":27514,"ts":1454427988},{"id":396724,"ts":1531842729},{"id":1238476,"ts":1539260965},{"id":334410,"ts":1491916827},{"id":710130,"ts":1455574816},
		{"id":1120234,"ts":1478873755},{"id":1240592,"ts":1516373787},{"id":854966,"ts":1523072145}]}`)
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
		string(a.Email.Buf[:a.Email.Len]),
		string(a.Phone[:]),
		string(a.PremiumStart[:]),
		string(a.PremiumFinish[:]))

	for _, v := range a.Interests[:] {
		if v == 0 {
			break
		}
		interest, _ := InterestDict.Value(uint64(v))
		t.Log(string(interest))
	}

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
		_ = a.Email.Buf[:a.Email.Len]
		_ = a.Phone[:]
		_ = a.PremiumStart[:]
		_ = a.PremiumFinish[:]

		//for _, v := range a.Interests[:] {
		//	if v == 0 {
		//		break
		//	}
		//	_, _ = InterestDict.Value(uint64(v))
		//}

	}
}
