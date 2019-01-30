package proto

import (
	"testing"
)

var (
	testLikes = []byte(`{"likes":[{"liker":1169657,"likee":298566,"ts":1494898074},{"liker":715791,"likee":1030468,"ts":1541630468},
		{"liker":632431,"likee":814356,"ts":1534472543},{"liker":1169657,"likee":923086,"ts":1495577566},{"liker":477545,"likee":948866,"ts":1536331906},
		{"liker":1169657,"likee":568328,"ts":1522815727},{"liker":85426,"likee":352581,"ts":1497093003},{"liker":426021,"likee":1062200,"ts":1521930517},
		{"liker":146379,"likee":905858,"ts":1510495114},{"liker":632431,"likee":145414,"ts":1455444224},{"liker":426021,"likee":862102,"ts":1534193149},
		{"liker":117506,"likee":1138017,"ts":1460740893},{"liker":1169657,"likee":1073410,"ts":1526648993},{"liker":117506,"likee":1045641,"ts":1503898494},
		{"liker":168367,"likee":283698,"ts":1491565569},{"liker":477545,"likee":501334,"ts":1511658865},{"liker":477545,"likee":900614,"ts":1475907250},
		{"liker":1169657,"likee":419422,"ts":1511831152},{"liker":794233,"likee":302858,"ts":1532358124},{"liker":146379,"likee":893460,"ts":1505206990},
		{"liker":146379,"likee":873800,"ts":1502234861},{"liker":975510,"likee":1209833,"ts":1501239362},{"liker":426021,"likee":1119858,"ts":1456531830},
		{"liker":1169657,"likee":155672,"ts":1487579999},{"liker":117506,"likee":1136279,"ts":1487422114},{"liker":794233,"likee":516406,"ts":1518838603},
		{"liker":579319,"likee":848440,"ts":1533847847},{"liker":146379,"likee":685752,"ts":1531054281},{"liker":168367,"likee":658614,"ts":1460764360},
		{"liker":794233,"likee":736112,"ts":1522970376},{"liker":146379,"likee":251074,"ts":1479567029},{"liker":426021,"likee":842356,"ts":1524377337},
		{"liker":426021,"likee":1151862,"ts":1518264697},{"liker":715791,"likee":30300,"ts":1483052935},{"liker":477545,"likee":127802,"ts":1456662658},
		{"liker":632431,"likee":975424,"ts":1493933790}]}`)
)

func TestLikesUnmarshalJSON(t *testing.T) {
	likes := &Likes{}
	_, ok := likes.UnmarshalJSON(testLikes)
	if !ok {
		t.Fatal()
	}
}

func BenchmarkLikesUnmarshalJSON(b *testing.B) {
	likes := &Likes{}
	for i := 0; i < b.N; i++ {
		_, ok := likes.UnmarshalJSON(testLikes)
		if !ok {
			b.Fatal()
		}
	}
}
