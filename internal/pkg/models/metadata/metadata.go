package metadata

import (
	"database/sql"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/pocketbase/dbx"
)

func New[T any, R any](k *koanf.Koanf, d *dbx.DB, e env.Env[T, R]) Metadata {
	return Metadata{
		// ID is static so it will be uniquely identified
		// on every deployment version.
		ID:         "26e4a9ae-67e4-430f-9263-de9a18d6160b",
		Version:    k.String("version"),
		DBStats:    d.DB().Stats(),
		DriverName: d.DriverName(),
		Env: envStats{
			ID: e.ID(),
			Services: func() map[string]string {
				m := make(map[string]string, len(e.Services()))
				for _, svc := range e.Services() {
					m[svc.ID()] = svc.Name()
				}
				return m
			}(),
		},
	}
}

type envStats struct {
	ID       string            `jsonapi:"attr,id"`
	Services map[string]string `jsonapi:"attr,services"`
}

type Metadata struct {
	ID         string      `jsonapi:"primary,metadata"`
	Version    string      `jsonapi:"attr,version"`
	DBStats    sql.DBStats `jsonapi:"attr,dbStats"`
	DriverName string      `jsonapi:"attr,driverName"`
	Env        envStats    `jsonapi:"attr,env"`
}
