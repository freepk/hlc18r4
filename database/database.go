package database

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync/atomic"
	"unsafe"

	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

var (
	FileCorruptedError = errors.New("File corrupted")
)

const (
	likesPerAccount = 35
)

type Like struct {
	ID uint32
	TS uint32
}

type Account struct {
	Sex        uint8
	LikesCount uint8
	LikesPos   uint32
}

type Database struct {
	accounts   []Account
	arenaPos   uint32
	arenaBytes []byte
}

func NewDatabase(accountsNum int) (*Database, error) {
	accounts := make([]Account, accountsNum)
	likesNum := accountsNum * likesPerAccount
	arenaSize := likesNum * 8
	arenaBytes := make([]byte, arenaSize)
	fmt.Println("New database, accountsNum", accountsNum, "likesNum", likesNum, "arenaBytes", len(arenaBytes))
	return &Database{accounts: accounts, arenaPos: 0, arenaBytes: arenaBytes}, nil
}

func (db *Database) Ping() {
}

func (db *Database) PrintStats() {
}

func (db *Database) NewAccount(src *proto.Account) error {
	dst := &db.accounts[src.ID]
	dst.Sex = src.Sex[0]
	n := len(src.Likes)
	if n > 0 {
		likesCount := uint8(n)
		likesSize := uint32(likesCount * 8)
		likesPos := atomic.AddUint32(&db.arenaPos, likesSize) - likesSize
		arenaSlice := (*reflect.SliceHeader)(unsafe.Pointer(&db.arenaBytes))
		likesSlice := reflect.SliceHeader{
			Data: arenaSlice.Data + uintptr(likesPos),
			Len:  int(likesCount),
			Cap:  int(likesCount),
		}
		likes := *(*[]Like)(unsafe.Pointer(&likesSlice))
		for i := 0; i < n; i++ {
			likes[i].ID = uint32(src.Likes[i].ID)
			likes[i].TS = uint32(src.Likes[i].TS)
		}
		dst.LikesCount = likesCount
		dst.LikesPos = likesPos
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
