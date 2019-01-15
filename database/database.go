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
	sync.RWMutex
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
	accounts := make([]Account, (accountsNum * 105 / 100))
	log.Println("New database, accountsNum", accountsNum, "allocated", (accountsNum * 105 / 100))
	return &Database{
		fnames:    fnames,
		snames:    snames,
		countries: countries,
		cities:    cities,
		interests: interests,
		accounts:  accounts}, nil
}

func (db *Database) Ping() {
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
	dst.Interests = make([]uint8, 0, len(src.Interests))
	for _, interest := range src.Interests {
		interest, _ := db.interests.GetKey(interest)
		dst.Interests = append(dst.Interests, uint8(interest))
	}
	dst.LikesTo = make([]Like, 0, len(src.Likes))
	for _, like := range src.Likes {
		dst.LikesTo = append(dst.LikesTo, Like{uint32(like.ID), uint32(like.TS)})
	}
	inserted := uint32(src.ID)
	db.updateLastInserted(inserted)
	return nil
}

func (db *Database) buildLikesFrom() {
	log.Println("Counting LikesFrom")
	counters := make([]int, db.lastInserted + 1)
	for _, liker := range db.accounts {
		for _, like := range liker.LikesTo {
			counters[like.ID]++
		}
	}
	log.Println("Allocating LikesFrom")
	for id, count := range counters {
		likee := &db.accounts[id]
		if count == len(likee.LikesFrom) {
			counters[id] = 0
		} else {
			likee.LikesFrom = make([]Like, 0, count)
		}
	}
	log.Println("Copying LikesFrom")
	for id, liker := range db.accounts {
		for _, like := range liker.LikesTo {
			if counters[like.ID] > 0 {
				likee := &db.accounts[like.ID]
				likee.LikesFrom = append(likee.LikesFrom, Like{uint32(id), uint32(like.TS)})
			}
		}
	}
}

func (db *Database) BuildIndexes() {
	log.Println("Build Indexes, lastInserted", db.lastInserted)
	db.buildLikesFrom()
	db.buildLikesFrom()
	db.buildLikesFrom()
	db.buildLikesFrom()
	db.buildLikesFrom()
	db.buildLikesFrom()
	db.buildLikesFrom()
	db.buildLikesFrom()
}

func (db *Database) updateLastInserted(inserted uint32) {
	lastInserted := atomic.LoadUint32(&db.lastInserted)
	if inserted > lastInserted {
		db.Lock()
		if inserted > db.lastInserted {
			db.lastInserted = inserted
		}
		db.Unlock()
	}
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
