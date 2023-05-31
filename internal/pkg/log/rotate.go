package log

import (
	"fmt"
	"time"
)

// shouldRotate answers the question if
// a rotate should be applied.
func (m manager) shouldRotate() (bool, error) {
	info, err := m.file.Stat()
	if err != nil {
		return false, err
	}
	if info.Size() >= int64(m.maxSize) {
		return true, nil
	}
	return false, nil
}

// createBackup creates a log backup
// in the format <name>_<timestamp>.log
func (m manager) createBackup() error {
	filename := fmt.Sprintf("%s_%d.log", m.svc.Name(), time.Now().Unix())
	fmt.Println(filename)
	return nil
}
func (m manager) compressFile() {}

func (m manager) rotate() error {
	ok, err := m.shouldRotate()
	if err != nil {
		return err
	}
	if !ok {
		// no rotate needed
		return nil
	}
	if err := m.createBackup(); err != nil {
		return err
	}
	return nil
}
