package service

import (
	//"errors"
	"io"
	//"github.com/klauspost/compress/zip"
	//"gitlab.com/freepk/hlc18r4/parse"
	//"gitlab.com/freepk/hlc18r4/proto"
)

const (
	accountsPerFile = 10000
)

func RestoreAccountsService(path string) (*AccountsService, error) {
	return nil, nil
	/*
		arch, err := zip.OpenReader(path)
		if err != nil {
			return nil, err
		}
		defer arch.Close()
		size := len(arch.File) * accountsPerFile * 105 / 100
		svc := NewAccountsService(uint32(size))
		for _, file := range arch.File {
			src, err := file.Open()
			if err != nil {
				return nil, err
			}
			if !svc.ReadFrom(src)
			if err != nil {
				return nil, err
			}
			src.Close()
		}
		return svc, nil
	*/
}

func (svc *AccountsService) ReadFrom(src io.Reader) bool {
	return false
	/*
		buf := make([]byte, 8192)
		headerSize := 14
		num, err := src.Read(buf[:headerSize])
		if err != nil {
			return err
		}
		if num != headerSize || string(buf[:headerSize]) != `{"accounts": [` {
			return false
		}
		acc := &proto.Account{}
		tailSize := 0
		for {
			num, err = src.Read(buf[tailSize:])
			if num > 0 {
				num += tailSize
				tail, ok := parse.ParseSymbol(buf[:num], ',')
				for {
					acc.Reset()
					tail, ok = acc.UnmarshalJSON(tail)
					if !ok {
						break
					}
					if !svc.Create(acc) {
						return false
					}
					tail, ok = parse.ParseSymbol(tail, ',')
				}
				tailSize = copy(buf, tail)
			} else if err == io.EOF {
				break
			} else if err != nil {
				return false
			}
		}
		return true
	*/
}
