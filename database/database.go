package database

import (
	"errors"
	"fmt"
	"io"

	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

var (
	FileCorruptedError = errors.New("File corrupted")
)

const (
	likesPerAccount = 34
)

type like struct {
	id uint32
	ts uint32
}

type account struct {
	sex    uint8
	likees []like
}

type Database struct {
	accounts []account
	likees   []like
}

func NewDatabase(accountsNum int) (*Database, error) {
	accounts := make([]account, accountsNum)
	likesNum := accountsNum * likesPerAccount
	likees := make([]like, likesNum)
	fmt.Println("New database, accountsNum", accountsNum, "likesNum", likesNum)
	return &Database{accounts: accounts, likees: likees}, nil
}

func (db *Database) Ping() {
}

func (db *Database) NewAccount(src *proto.Account) error {
	dst := &db.accounts[src.ID]
	dst.sex = src.Sex[0]
	n := len(src.Likes)
	dst.likees = db.likees[:n]
	db.likees = db.likees[n:]
	for i := 0; i < n; i++ {
		dst.likees[i].id = uint32(src.Likes[i].ID)
		dst.likees[i].ts = uint32(src.Likes[i].TS)
	}
	return nil
}

func (db *Database) ReadFrom(r io.Reader) error {
	buf := make([]byte, 8192)
	headerSize := 14
	n, err := r.Read(buf[:headerSize])
	if err != nil {
		return err
	}
	if n != headerSize || string(buf[:headerSize]) != `{"accounts": [` {
		return FileCorruptedError
	}
	account := &proto.Account{}
	tailSize := 0
	for {
		n, err = r.Read(buf[tailSize:])
		if n > 0 {
			n += tailSize
			tail, ok := parse.ParseSymbol(buf[:n], ',')
			account.Reset()
			tail, ok = account.UnmarshalJSON(tail)
			if !ok {
				break
			}
			err = db.NewAccount(account)
			if err != nil {
				return err
			}
			tailSize = copy(buf, tail)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
