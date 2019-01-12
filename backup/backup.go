package backup

import (
	"sync"
	"time"

	"github.com/klauspost/compress/zip"
	"gitlab.com/freepk/hlc18r4/database"
)

func Restore(path string) (*database.Database, error) {
	numAccounts := 1400000
	currentTime := time.Now()
	db := database.NewDatabase(numAccounts, currentTime)
	arch, err := zip.OpenReader(path + "data.zip")
	if err != nil {
		return nil, err
	}
	defer arch.Close()
	n := len(arch.File)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(n)
	errChan := make(chan error, n)
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
