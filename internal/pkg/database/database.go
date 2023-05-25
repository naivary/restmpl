package database

import (
	"database/sql"
	"fmt"

	"github.com/knadh/koanf/v2"
	_ "modernc.org/sqlite"
)

func buildDsn(k *koanf.Koanf) string {
	return fmt.Sprintf("file:%s_%s.db", k.String("name"), k.String("version"))
}

func initPragmas(db *sql.DB) error {
	query := `
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = TRUE;	
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
