package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/must"
	"github.com/pocketbase/dbx"
)

func buildDataDir(k *koanf.Koanf) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s_data", k.String("name"))
	path := filepath.Join(wd, name, k.String("version"))
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(filepath.Join(path, "backup"), os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil
}

func buildDsn(k *koanf.Koanf) string {
	dir, err := buildDataDir(k)
	must.Must(err)
	return fmt.Sprintf("file:%s/%s.db", dir, k.String("name"))
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
