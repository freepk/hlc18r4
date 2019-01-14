package database

import (
	"errors"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/freepk/dictionary"
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
	Brith         uint32
	Joined        uint32
	Fname         uint8
	Sname         uint16
	Country       uint8
	City          uint16
	PremiumStart  uint32
	PremiumFinish uint32
	Interests     []uint8
	LikesTo       []Like
	LikesFrom     []Like
}

type Database struct {
	fnames       *dictionary.Dictionary
	snames       *dictionary.Dictionary
	countries    *dictionary.Dictionary
	cities       *dictionary.Dictionary
	interests    *dictionary.Dictionary
	accounts     []Account
	lastInserted uint32
}

func NewDatabase(accountsNum int) (*Database, error) {
	fnames, _ := dictionary.NewDictionary(8)
	snames, _ := dictionary.NewDictionary(12)
	countries, _ := dictionary.NewDictionary(8)
	cities, _ := dictionary.NewDictionary(12)
	interests, _ := dictionary.NewDictionary(8)
	accounts := make([]Account, accountsNum*105/100)
	log.Println("New database, accountsNum", accountsNum, "allocated", accountsNum*105/100)
	return &Database{
		fnames:       fnames,
		snames:       snames,
		countries:    countries,
		cities:       cities,
		interests:    interests,
		accounts:     accounts,
		lastInserted: 0}, nil
}

func (db *Database) Ping() {
}

func (db *Database) updateLastInserted(id uint32) {
	last := atomic.LoadUint32(&db.lastInserted)
	if id > last {
		atomic.CompareAndSwapUint32(&db.lastInserted, last, id)
	}
}

func (db *Database) NewAccount(src *proto.Account) error {
	dst := &db.accounts[src.ID]
	// ID
	dst.Brith = uint32(src.Birth)
	dst.Joined = uint32(src.Joined)
	// Email
	if fname, err := db.fnames.GetKey(src.Fname); err == nil {
		dst.Fname = uint8(fname)
	}
	if sname, err := db.snames.GetKey(src.Sname); err == nil {
		dst.Sname = uint16(sname)
	}
	// Phone
	if country, err := db.countries.GetKey(src.Country); err == nil {
		dst.Country = uint8(country)
	}
	if city, err := db.cities.GetKey(src.City); err == nil {
		dst.City = uint16(city)
	}
	// Status
	dst.PremiumStart = uint32(src.Premium.Start)
	dst.PremiumFinish = uint32(src.Premium.Finish)
	dst.Interests = make([]uint8, len(src.Interests))
	for i := 0; i < len(src.Interests); i++ {
		if interest, err := db.interests.GetKey(src.Interests[i]); err == nil {
			dst.Interests[i] = uint8(interest)
		}
	}
	dst.LikesTo = make([]Like, len(src.Likes))
	for i := 0; i < len(src.Likes); i++ {
		dst.LikesTo[i].ID = uint32(src.Likes[i].ID)
		dst.LikesTo[i].TS = uint32(src.Likes[i].TS)
	}
	db.updateLastInserted(uint32(src.ID))
	return nil
}

type WalkCallback func(id uint32)

func (db *Database) WalkAccounts(first, last, step uint32, callback WalkCallback) {
	if first >= last {
		return
	}
	waitGroup := &sync.WaitGroup{}
	for i := uint32(first); i <= last; i += step {
		a := i
		b := i + step - 1
		if b > last {
			b = last
		}
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for id := a; id <= b; id++ {
				callback(id)
			}
		}()
	}
	waitGroup.Wait()
}

func (db *Database) buildLikesFrom() {
	log.Println("Counting LikesFrom")
	likesFromCount := make([]uint32, db.lastInserted+1)
	db.WalkAccounts(1, db.lastInserted, 10000, func(id uint32) {
		likesTo := db.accounts[id].LikesTo
		n := len(likesTo)
		for i := 0; i < n; i++ {
			atomic.AddUint32(&likesFromCount[likesTo[i].ID], 1)
		}
	})
	log.Println("Allocating LikesFrom")
	unused := uint32(0)
	allocated := uint32(0)
	unchanged := uint32(0)
	relocated := uint32(0)
	db.WalkAccounts(1, db.lastInserted, 10000, func(id uint32) {
		required := likesFromCount[id]
		if required == 0 {
			atomic.AddUint32(&unused, 1)
			return
		}
		capacity := cap(db.accounts[id].LikesFrom)
		if capacity == 0 {
			atomic.AddUint32(&allocated, 1)
			db.accounts[id].LikesFrom = make([]Like, 0, required)
			return
		}
		if uint32(capacity) == required {
			atomic.AddUint32(&unchanged, 1)
			return
		}
		atomic.AddUint32(&relocated, 1)
		buf := make([]Like, 0, required)
		db.accounts[id].LikesFrom = append(buf, db.accounts[id].LikesFrom...)
	})
	log.Println("Unused", unused, "allocated", allocated, "unchanged", unchanged, "relocated", relocated)
	log.Println("Copying LikesFrom")
	db.WalkAccounts(1, db.lastInserted, 10000, func(id uint32) {
		likesTo := db.accounts[id].LikesTo
		n := len(likesTo)
		for i := 0; i < n; i++ {
		}
	})
}

func (db *Database) BuildIndexes() {
	log.Println("Build indexes, lastInserted", db.lastInserted)
	db.buildLikesFrom()
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
			for {
				account.Reset()
				tail, ok = account.UnmarshalJSON(tail)
				if !ok {
					break
				}
				err = db.NewAccount(account)
				if err != nil {
					return err
				}
				tail, ok = parse.ParseSymbol(tail, ',')
			}
			tailSize = copy(buf, tail)
		} else if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}
