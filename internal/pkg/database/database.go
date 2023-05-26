package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/must"
	_ "modernc.org/sqlite"
)

func buildDataDir(k *koanf.Koanf) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s_data", k.String("name"))
	path := filepath.Join(wd, name)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}

func buildDsn(k *koanf.Koanf) string {
	dir, err := buildDataDir(k)
	must.Must(err)
	return fmt.Sprintf("file:%s/%s_%s.db", dir, k.String("name"), k.String("version"))
}

func initPragmas(db *sql.DB) error {
	query := `
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = 1;	
	`
	_, err := db.Exec(query, nil)
	if err != nil {
		return err
	}
	return nil
}

func initOptions(db *sql.DB) {
	db.SetMaxOpenConns(1)
}

// inMem provides an in-memory
// sqlite database which is
// used for test purposes
func inMem() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file::memory:")
	if err != nil {
		return nil, err
	}
	err = initPragmas(db)
	if err != nil {
		return nil, err
	}
	initOptions(db)
	return db, nil
}

// Connect creates a connectiont to the sqlite database.
// If `k` is nil, the database will be created in memory.
// This should be only considered for testing puposes.
func Connect(k *koanf.Koanf) (*sql.DB, error) {
	// check if the connection is needed for testing purposes
	if k == nil {
		return inMem()
	}

	db, err := sql.Open("sqlite", buildDsn(k))
	if err != nil {
		return nil, err
	}
	err = initPragmas(db)
	if err != nil {
		return nil, err
	}
	initOptions(db)
	return db, nil
}
