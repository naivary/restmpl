package env

import (
	"github.com/naivary/apitmpl/internal/pkg/service"
)

type meta struct {
	ID       string          `jsonapi:"primary,meta"`
	DBDriver string          `jsonapi:"attr,driverName"`
	Svcs     []*service.Info `jsonapi:"attr,services"`
	Version  string          `jsonapi:"attr,version"`
	Name     string          `jsonapi:"attr,name"`
}

func (a API) newMeta() meta {
	m := meta{
		ID: a.ID(),
	}
	m.DBDriver = a.db.DriverName()
	m.Version = a.k.String("version")
	m.Name = a.k.String("name")
	return m
}
