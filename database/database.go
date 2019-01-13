package database

import (
	"errors"
	"fmt"
	"io"

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
	LikesFrom     []uint32
}

type Database struct {
	fnames    *dictionary.Dictionary
	snames    *dictionary.Dictionary
	countries *dictionary.Dictionary
	cities    *dictionary.Dictionary
	interests *dictionary.Dictionary
	accounts  []Account
}

func NewDatabase(accountsNum int) (*Database, error) {
	accounts := make([]Account, accountsNum)
	fnames, _ := dictionary.NewDictionary(8)
	snames, _ := dictionary.NewDictionary(12)
	countries, _ := dictionary.NewDictionary(8)
	cities, _ := dictionary.NewDictionary(12)
	interests, _ := dictionary.NewDictionary(8)
	fmt.Println("New database, accountsNum", accountsNum)
	return &Database{
		accounts:  accounts,
		fnames:    fnames,
		snames:    snames,
		countries: countries,
		cities:    cities,
		interests: interests}, nil
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
