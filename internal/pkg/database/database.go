package database

import (
	"github.com/knadh/koanf/v2"
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

// inMem provides an in-memory
// sqlite database which is
// used for test purposes
func inMem() (*dbx.DB, error) {
	db, err := dbx.Open("sqlite", "file::memory:")
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
// This should be only considered for testing purposes.
func Connect(k *koanf.Koanf) (*dbx.DB, error) {
	// check if the connection is needed for testing purposes
	if k == nil {
		return inMem()
	}
	db, err := dbx.Open("sqlite", buildDsn(k))
	if err != nil {
		return nil, err
	}
	initOptions(db)
	err = initPragmas(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
