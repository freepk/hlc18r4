package main

import (
	//"errors"
	//"io"

	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/repo"
	//"gitlab.com/freepk/hlc18r4/parse"
	//"gitlab.com/freepk/hlc18r4/proto"
)

func restore(path string) (*repo.AccountsRepo, error) {
	arch, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer arch.Close()
	for _, file := range arch.File {
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		src.Close()
	}
	return nil, nil
}

/*
func (svc *AccountsService) ReadFrom(src io.Reader) error {
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
				acc.Reset()
				if tail, ok = acc.UnmarshalJSON(tail); !ok {
					break
				}
				if !svc.Create(acc) {
					return AccountsServiceReadFromError
				}
			}
			pos = copy(buf, tail)
		} else if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}
*/
