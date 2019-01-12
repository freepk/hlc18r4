package database

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
)

var (
	FileCorruptedError = errors.New("File corrupted")
	WrongStateError    = errors.New("Database in wrong state")
	SexValidationError = errors.New("Sex is not valid")
)

type DatabaseState int

const (
	OnlineState = DatabaseState(iota)
	RecoveryState
)

const (
	likesPerAccount = 32
)

type account struct {
	sex              uint8
	maleLikesCount   uint8
	femaleLikesCount uint8
	likesCount       uint8
}

type likex struct {
	liker uint32
	likee uint32
	ts    uint32
}

type Database struct {
	sync.RWMutex
	state     DatabaseState
	accounts  []account
	likesTemp []likex
}

func NewDatabase(accountsNum int) (*Database, error) {
	accounts := make([]account, accountsNum)
	fmt.Println("New database, accountsNum", accountsNum)
	return &Database{state: OnlineState, accounts: accounts}, nil
}

func (db *Database) State() DatabaseState {
	db.RLock()
	defer db.RUnlock()
	return db.state
}

func (db *Database) addLike(liker, likee, ts uint32) {
	db.RLock()
	if db.state == RecoveryState {
		db.likesTemp = append(db.likesTemp, likex{liker, likee, ts})
		{
			likerAcc := &db.accounts[liker]
			likerAcc.likesCount++
			likeeAcc := &db.accounts[likee]
			switch likerAcc.sex {
			case 'm':
				likeeAcc.maleLikesCount++
			case 'f':
				likeeAcc.femaleLikesCount++
			}
		}
	}
	db.RUnlock()
}

func (db *Database) estimateLikesSize() int {
	n := len(db.accounts)
	size := 0
	for i := 0; i < n; i++ {
		self := &db.accounts[i]
		size += int(self.likesCount) / 2 * 4
		size += int(self.maleLikesCount+self.femaleLikesCount) * 8
		fmt.Println("sex", self.sex, "likesCount", self.likesCount, "femaleLikesCount", self.femaleLikesCount, "maleLikesCount", self.maleLikesCount)
	}
	return size
}

func (db *Database) NewAccount(src *proto.Account) error {
	dst := &db.accounts[src.ID]
	dst.sex = src.Sex[0]
	n := len(src.Likes)
	for i := 0; i < n; i++ {
		db.addLike(uint32(src.ID), uint32(src.Likes[i].ID), uint32(src.Likes[i].TS))
	}
	return nil
}

func (db *Database) StartRecovery() error {
	db.Lock()
	defer db.Unlock()
	if db.state != OnlineState {
		return WrongStateError
	}
	numLikesInTemp := len(db.accounts) * likesPerAccount
	db.likesTemp = make([]likex, 0, numLikesInTemp)
	db.state = RecoveryState
	fmt.Println("Start revovery, numLikesInTemp", numLikesInTemp)
	return nil
}

func (db *Database) FinishRecovery() error {
	db.Lock()
	defer db.Unlock()
	if db.state != RecoveryState {
		return WrongStateError
	}
	db.likesTemp = nil
	db.state = OnlineState
	fmt.Println("Finish revovery, estimateLikesSize", db.estimateLikesSize())
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
