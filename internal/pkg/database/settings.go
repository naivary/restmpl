package database

import (
	"fmt"

	"github.com/knadh/koanf/v2"
	"github.com/pocketbase/dbx"
)

func buildDsn(k *koanf.Koanf) string {
	return fmt.Sprintf("file:%s/%s.db", k.String("versionDir"), k.String("name"))
}

func initPragmas(db *dbx.DB) error {
	query := `
		PRAGMA busy_timeout       = 10000;
		PRAGMA journal_mode       = WAL;
		PRAGMA journal_size_limit = 200000000;
		PRAGMA synchronous        = NORMAL;
		PRAGMA foreign_keys       = 1;	
	`
	if _, err := db.NewQuery(query).Execute(); err != nil {
		return err
	}
	return nil
}

func initOptions(db *dbx.DB) {
	db.DB().SetMaxOpenConns(1)
}
