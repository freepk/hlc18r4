package repo

import (
	"testing"

	"gitlab.com/freepk/hlc18r4/proto"
)

func TestAccountsRepo(t *testing.T) {
	rep := NewAccountsRepo(1024)
	if rep.Get(-1) != nil {
		t.Fail()
	}
	if rep.Get(4096) != nil {
		t.Fail()
	}
	acc := &proto.Account{ID: 10, Email: []byte("test@mail.ru"), Interests: []uint8{1, 2, 3, 4, 5}}
	rep.Add(acc)
	acc.Email[0] = 'f'
	acc.Interests[0] = 10
	tmp := rep.Get(10)
	t.Log(string(tmp.Email), tmp.Interests)
	acc.Interests = acc.Interests[:0]
	acc.Interests = append(acc.Interests, 50)
	acc.Interests = append(acc.Interests, 60)
	t.Log(string(tmp.Email), tmp.Interests)
	tmp.Email = append([]byte{}, tmp.Email...)
	tmp.Interests = append([]uint8{}, tmp.Interests...)

	acc.Email[0] = 'r'
	acc.Interests[0] = 20
	t.Log(string(tmp.Email), tmp.Interests)

}
