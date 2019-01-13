package backup

import (
	"sync"

	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/database"
)

const accountsPerFile = 10000

func Restore(path string) (*database.Database, error) {
	arch, err := zip.OpenReader(path + "data.zip")
	if err != nil {
		return nil, err
	}
	defer arch.Close()
	n := len(arch.File)
	accountsNum := n * accountsPerFile
	db, err := database.NewDatabase(accountsNum)
	if err != nil {
		return nil, err
	}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(n)
	errChan := make(chan error)
	for i := 0; i < n; i++ {
		file := arch.File[i]
		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		go func() {
			defer src.Close()
			defer waitGroup.Done()
			err := db.ReadFrom(src)
			if err != nil {
				errChan <- err
			}
		}()
	}
	waitGroup.Wait()
	close(errChan)
	for err := range errChan {
		return nil, err
	}
	return db, nil
}
