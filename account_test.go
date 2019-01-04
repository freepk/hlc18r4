package main

import (
	"encoding/json"
	"github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"testing"
)

var testAccount = []byte(`{"likes":[{"id":484053,"ts":1476605353},{"id":1055765,"ts":1467406524},{"id":793185,"ts":1534580707},{"id":955917,"ts":1493975640},
	{"id":375379,"ts":1539985529},{"id":836773,"ts":1481693076},{"id":691655,"ts":1533981341},{"id":1299873,"ts":1524997142},{"id":179795,"ts":1522779391},
	{"id":1126729,"ts":1514445485},{"id":206461,"ts":1482808241},{"id":991685,"ts":1530470598},{"id":1020737,"ts":1515144827},{"id":1173581,"ts":1515302882},
	{"id":1261119,"ts":1463651782},{"id":885573,"ts":1485624715},{"id":919053,"ts":1470545613},{"id":1271139,"ts":1506959707},{"id":702175,"ts":1454866547},
	{"id":561969,"ts":1490090604},{"id":553607,"ts":1476342851},{"id":1219173,"ts":1511964110},{"id":1184959,"ts":1516403946},{"id":1047217,"ts":1520897805},
	{"id":372333,"ts":1499389210},{"id":22407,"ts":1482159817},{"id":115595,"ts":1470196720},{"id":1237357,"ts":1482467303},{"id":332051,"ts":1530385062},
	{"id":147429,"ts":1502689778},{"id":702575,"ts":1532061307},{"id":1018911,"ts":1512525241},{"id":898517,"ts":1473651435},{"id":1188903,"ts":1502156143},
	{"id":153049,"ts":1517466016},{"id":1162483,"ts":1499023984},{"id":15269,"ts":1502187927},{"id":996347,"ts":1528218998},{"id":147613,"ts":1468601563},
	{"id":392399,"ts":1510199233},{"id":736801,"ts":1463709303},{"id":1130559,"ts":1469541346},{"id":624079,"ts":1460211737},{"id":264859,"ts":1536663165},
	{"id":64879,"ts":1514674199},{"id":353639,"ts":1531604504},{"id":1119055,"ts":1473334604},{"id":1089201,"ts":1530561229},{"id":279193,"ts":1503074606},
	{"id":888793,"ts":1466196633},{"id":934971,"ts":1462351420},{"id":1094279,"ts":1531759443},{"id":861239,"ts":1522612147},{"id":920775,"ts":1453119149},
	{"id":1125623,"ts":1517244991},{"id":280299,"ts":1513154136},{"id":991349,"ts":1466842610},{"id":627749,"ts":1481164079},{"id":673817,"ts":1482148311},
	{"id":315959,"ts":1488954321},{"id":1149819,"ts":1504710243},{"id":154063,"ts":1514595923},{"id":1234339,"ts":1499566684},{"id":255459,"ts":1488352776},
	{"id":662745,"ts":1481216775},{"id":213813,"ts":1532623824},{"id":220947,"ts":1539456021},{"id":343395,"ts":1466892947},{"id":426101,"ts":1508163208},
	{"id":727005,"ts":1511899829},{"id":1075081,"ts":1535117975},{"id":328859,"ts":1492901718},{"id":141525,"ts":1491766730},{"id":1281457,"ts":1508215106},
	{"id":185571,"ts":1492654210},{"id":853249,"ts":1501262537},{"id":627375,"ts":1501865401},{"id":567469,"ts":1523708266},{"id":179143,"ts":1538726433},
	{"id":30491,"ts":1474338607},{"id":451121,"ts":1465475793},{"id":1110529,"ts":1505407844},{"id":33841,"ts":1520879410},{"id":1116037,"ts":1478847461}],
	"id":1300026,"joined":1332115200,"birth":721360390}`)

func TestAccountParse(t *testing.T) {
	parseAccount(testAccount)
}

func TestAccountParseFastJson(t *testing.T) {
	var p fastjson.Parser
	v, err := p.ParseBytes(testAccount)
	t.Log(v, err)
}

func TestAccountParseStandart(t *testing.T) {
	var account struct {
		ID      int `json:"id"`
		Joined  int `json:"joined"`
		Birth   int `json:"birth"`
		Premium struct {
			Finish int `json:"finish"`
			Start  int `json:"start"`
		} `json:"premium"`
		Likes []struct {
			ID int `json:"id"`
			Ts int `json:"ts"`
		} `json:"likes"`
	}
	err := json.Unmarshal(testAccount, &account)
	t.Log(account, err)
}

func TestAccountParseJsoniter(t *testing.T) {
	var account struct {
		ID      int `json:"id"`
		Joined  int `json:"joined"`
		Birth   int `json:"birth"`
		Premium struct {
			Finish int `json:"finish"`
			Start  int `json:"start"`
		} `json:"premium"`
		Likes []struct {
			ID int `json:"id"`
			Ts int `json:"ts"`
		} `json:"likes"`
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(testAccount, &account)
	t.Log(account, err)
}

func TestAccountParseEasyJson(t *testing.T) {
	account := &EasyAccount{}
	err := account.UnmarshalJSON(testAccount)
	t.Log(account, err)
}

func BenchmarkAccountParseInHouse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseAccountX(testAccount)
	}
}

func BenchmarkAccountParseFastjson(b *testing.B) {
	var p fastjson.Parser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.ParseBytes(testAccount)
	}
}

func BenchmarkAccountParseJsoniter(b *testing.B) {
	var account struct {
		ID      int `json:"id"`
		Joined  int `json:"joined"`
		Birth   int `json:"birth"`
		Premium struct {
			Finish int `json:"finish"`
			Start  int `json:"start"`
		} `json:"premium"`
		Likes []struct {
			ID int `json:"id"`
			Ts int `json:"ts"`
		} `json:"likes"`
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		json.Unmarshal(testAccount, &account)
	}
}

func BenchmarkAccountParseStandart(b *testing.B) {
	var account struct {
		ID      int `json:"id"`
		Joined  int `json:"joined"`
		Birth   int `json:"birth"`
		Premium struct {
			Finish int `json:"finish"`
			Start  int `json:"start"`
		} `json:"premium"`
		Likes []struct {
			ID int `json:"id"`
			Ts int `json:"ts"`
		} `json:"likes"`
	}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(testAccount, &account)
	}
}

func BenchmarkAccountParseEasyJson(b *testing.B) {
	account := &EasyAccount{}
	for i := 0; i < b.N; i++ {
		account.UnmarshalJSON(testAccount)
	}
}
