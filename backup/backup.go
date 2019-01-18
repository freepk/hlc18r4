package backup

import (
	"errors"
	"io"
	"log"
	"sync"

	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/parse"
	"gitlab.com/freepk/hlc18r4/proto"
	"gitlab.com/freepk/hlc18r4/repo"
)

const (
	accountsPerFile = 10000
)

var (
	ReadError = errors.New("Read error")
)

func Restore(name string) (*repo.AccountsRepo, error) {
	arch, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer arch.Close()
	num := accountsPerFile * len(arch.File) * 105 / 100
	log.Println("New AccountsRepo", num)
	rep := repo.NewAccountsRepo(num)
	grp := &sync.WaitGroup{}
	for _, file := range arch.File {
		if src, err := file.Open(); err != nil {
			return nil, err
		} else {
			grp.Add(1)
			go func() {
				defer src.Close()
				defer grp.Done()
				if err := readFrom(rep, src); err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
	grp.Wait()
	return rep, nil
}

func readFrom(rep *repo.AccountsRepo, src io.Reader) error {
	buf := make([]byte, 8192)
	num, err := src.Read(buf[:14])
	if err != nil {
		return err
	}
	acc := &proto.Account{}
	pos := 0
	for {
		if num, err = src.Read(buf[pos:]); num > 0 {
			num += pos
			tail, ok := buf[:num], true
			for {
				tail, _ = parse.ParseSymbol(tail, ',')
				if tail, ok = acc.UnmarshalJSON(tail); !ok {
					break
				}
				rep.Add(acc.Clone())
			}
			if pos = copy(buf, tail); pos == len(buf) {
				return ReadError
			}
		} else if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}
