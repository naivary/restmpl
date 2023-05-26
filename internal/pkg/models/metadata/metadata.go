package metadata

import (
	"database/sql"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/pocketbase/dbx"
)

func New(k *koanf.Koanf, d *dbx.DB) Metadata {
	return Metadata{
		// ID is static so it will be uniquely identified
		// on every deployment version.
		ID:         "26e4a9ae-67e4-430f-9263-de9a18d6160b",
		Version:    k.String("version"),
		DBStats:    d.DB().Stats(),
		DriverName: database.GetDriverName(d.DB()),
	}
}

type Metadata struct {
	ID         string      `jsonapi:"primary,metadata"`
	Version    string      `jsonapi:"attr,version"`
	DBStats    sql.DBStats `jsonapi:"attr,dbStats"`
	DriverName string      `jsonapi:"attr,driverName"`
}
