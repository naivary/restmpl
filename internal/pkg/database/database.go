package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	_ "modernc.org/sqlite"
)

func Connect(k *koanf.Koanf) (*sql.DB, error) {
	name := fmt.Sprintf("file:%s_%s.db?mode=", k.String("name"), k.String("version"))
	return sql.Open("sqlite", filepath.Join(k.String("db"), name))
}

// InMemConnect creates a sqlite database connection
// which is stored in memory. It is used for test purposes
func InMemConnect() (*sql.DB, error) {
	return sql.Open("sqlite", "file::memory:?mode=rwc&_journal_mode=WAL&_fk=true")
}
