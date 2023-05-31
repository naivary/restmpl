package models

import (
	"database/sql"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/pocketbase/dbx"
)

type Meta struct {
	ID       string          `jsonapi:"primary,meta"`
	DBDriver string          `jsonapi:"attr,driverName"`
	DBStats  sql.DBStats     `jsonapi:"attr,dbStats"`
	Svcs     []*service.Info `jsonapi:"attr,services"`
	Version  string          `jsonapi:"attr,version"`
	Name     string          `jsonapi:"attr,name"`
}

func NewMeta(k *koanf.Koanf, db *dbx.DB, env env.Env) Meta {
	m := Meta{}
	m.ID = env.ID()
	m.DBStats = db.DB().Stats()
	m.DBDriver = db.DriverName()
	m.Version = k.String("version")
	m.Name = k.String("name")
	return m
}
