package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

type Monitor struct {
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

// create hash for paths in map
func (m *Monitor) Now() (int, error) {
	var counter int
	for path, lastHash := range m.Paths {
		newHash, err := DirHash(path)
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			// execute backup
			err := m.act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash
			counter++
		}
	}
	return counter, nil
}

func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := fmt.Sprintf(m.Archiver.DestFmt()(
		time.Now().UnixNano()))
	// directly return errors if happened
	// its because it is ok if error receiver can handle it
	// and if not, leave its resolution to upper class 
	return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}
