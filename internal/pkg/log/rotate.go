package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/naivary/instance/internal/pkg/random"
)

func (m *manager) rotate() error {
	if ok, err := m.shouldRotate(); !ok || err != nil {
		return err
	}
	oldFile, err := m.changeCurrentLogFile()
	if err != nil {
		return err
	}

	backup, err := m.backup(oldFile)
	if err != nil {
		return err
	}
	err = m.removeOldFile(oldFile.Name())
	if err != nil {
		return err
	}
	m.addBackup(backup)
	return nil
}

// shouldRotate answers the question if
// a rotate should be applied.
func (m manager) shouldRotate() (bool, error) {
	info, err := m.file.Stat()
	if err != nil {
		return false, err
	}
	return info.Size() >= m.maxSize, nil
}

// changeCurrentLogFile creates a new logfile to which logger
// will log further incoming logs and returns the old log file
// which containing all the content.
func (m *manager) changeCurrentLogFile() (*os.File, error) {
	oldFilename := m.file.Name()
	id := random.ID(2)
	filename := fmt.Sprintf("%s_%s_%s.log", m.svc.Name(), m.svc.ID(), id)
	p := filepath.Join(m.k.String("logsDir"), filename)
	newFile, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	*m.file = *newFile
	oldFile, err := os.Open(oldFilename)
	if err != nil {
		return nil, err
	}
	return oldFile, nil
}

func (m *manager) compressBackupFile(src io.Reader) (*os.File, error) {
	filename := fmt.Sprintf("%s_%d.gz", m.svc.Name(), time.Now().Unix())
	p := filepath.Join(m.k.String("logsDir"), filename)
	backup, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	wr := gzip.NewWriter(backup)
	if _, err := io.Copy(wr, src); err != nil {
		return nil, err
	}
	if err := wr.Flush(); err != nil {
		return nil, err
	}
	return backup, nil
}

func (m *manager) rawBackupFile(src io.Reader) (*os.File, error) {
	filename := fmt.Sprintf("%s_%d.log", m.svc.Name(), time.Now().Unix())
	p := filepath.Join(m.k.String("logsDir"), filename)
	backup, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(backup, src); err != nil {
		return nil, err
	}
	return backup, nil
}

func (m *manager) backup(src io.Reader) (*os.File, error) {
	fmt.Println(len(m.backups), m.maxBackups)
	if len(m.backups) == m.maxBackups {
		m.deleteOldBackups()
	}
	if m.compress {
		return m.compressBackupFile(src)
	}
	return m.rawBackupFile(src)
}

func (m manager) removeOldFile(name string) error {
	return os.Remove(name)
}

func (m *manager) deleteOldBackups() error {
	for _, backup := range m.backups {
		if err := os.Remove(backup.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) addBackup(backup *os.File) {
	m.backups = append(m.backups, backup)
}
