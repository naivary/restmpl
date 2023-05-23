package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	_ "modernc.org/sqlite"
)

func Connect(k *koanf.Koanf) (*sql.DB, error) {
	name := fmt.Sprintf("%s_%s.db", k.String("name"), k.String("version"))
	return sql.Open("sqlite", filepath.Join(k.String("db"), name))
}
