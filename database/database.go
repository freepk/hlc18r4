package database

import (
	"errors"
	"fmt"
	"io"
	"time"

	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

var (
	WrongFileHeaderError = errors.New("Wrong file header")
	NullAccountSexError  = errors.New("Account sex is null")
)

type likex struct {
	liker uint32
	likee uint32
	ts uint32
}

type like struct {
	id uint32
	ts uint32
}

type account struct {
	sex         string
	maleLikee   []uint32 // likes to account from other males, only liker id
	femaleLikee []uint32 // likes to account from other females, only liker id
	liker       []like   // from account to others accounts with same sex (m->m or f->f, likee id and ts)
}

type Database struct {
	currentTime time.Time
	accounts    []account
	likesTail   []likex
}

func NewDatabase(numAccounts int, currentTime time.Time) *Database {
	accounts := make([]account, numAccounts)
	likesTail := make([]likex, 0, 45000000)
	return &Database{currentTime: currentTime, accounts: accounts, likesTail: likesTail}
}

func (db *Database) PrintStatus() {
	maleLiker := 0
	femaleLiker := 0
	liker := 0
	for i := 0; i < len(db.accounts); i++ {
		maleLiker += len(db.accounts[i].maleLiker)
		femaleLiker += len(db.accounts[i].femaleLiker)
		liker += len(db.accounts[i].liker)
	}
	fmt.Println("male", maleLiker, "female", femaleLiker, "liker", liker)
}

func (db *Database) ShrinkLikesTail() {
	temp := make([]struct {
		mlikee uint32
		flikee uint32
		liker uint32
	}, 1400000 )
	n := len(db.likesTail)
	for i := 0; i < n; i++ {
		liker := db.likesTail[i].liker
		likee := db.likesTail[i].liker
		src := &db.accounts[liker]
		dst := &db.accounts[likee]
		if src.sex == "" { fmt.Println("src.sex empty!!!") }
		if dst.sex == "" { fmt.Println("dst.sex empty!!!") }
	}
	_ = temp
}

func (db *Database) NewAccount(src *proto.Account) error {
	self := &db.accounts[src.ID]
	if len(src.Sex) == 0 {
		return NullAccountSexError
	}
	self.sex = string(src.Sex)
	n := len(src.Likes)
	for i:=0;i<n;i++{
		likee := uint32(src.Likes[i].ID)
		ts := uint32(src.Likes[i].TS)
		db.likesTail = append(db.likesTail, likex{uint32(src.ID), likee, ts})
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
		return WrongFileHeaderError
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
