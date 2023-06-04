package env

import (
	"github.com/naivary/apitmpl/internal/pkg/service"
)

type meta struct {
	ID       string          `json:"id"`
	DBDriver string          `json:"dbDriverName"`
	Svcs     []*service.Info `json:"services"`
	Version  string          `json:"version"`
	Name     string          `json:"name"`
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
