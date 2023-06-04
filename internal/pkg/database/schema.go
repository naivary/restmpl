package database

import "github.com/pocketbase/dbx"

func initSchema(db *dbx.DB) error {
	tables := []func(*dbx.DB) error{
		usersTable,
	}
	for _, initTable := range tables {
		if err := initTable(db); err != nil {
			return err
		}
	}
	return nil
}

func usersTable(db *dbx.DB) error {
	_, err := db.CreateTable("users", map[string]string{
		"id":         "TEXT PRIMRAY KEY",
		"username":   "TEXT",
		"password":   "TEXT",
		"email":      "TEXT UNIQUE",
		"created_at": "INTEGER",
	}).Execute()
	return err
}
