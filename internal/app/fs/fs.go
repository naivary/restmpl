package fs

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/filestore"
	"github.com/naivary/apitmpl/internal/pkg/logging"
	"github.com/naivary/apitmpl/internal/pkg/metrics"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/afero"
)

var _ service.Service = (*Fs)(nil)

type Fs struct {
	K *koanf.Koanf

	store filestore.Store[afero.File]
	l     logging.Manager
	m     metrics.Managee
}

func (f Fs) Health() (*service.Info, error) {
	if f.K == nil {
		return nil, errors.New("missing config manager")
	}

	if f.store == nil {
		return nil, errors.New("missing filestore")
	}

	if err := f.store.Health(); err != nil {
		return nil, errors.New("filestore unhealthy")
	}
	return &service.Info{
		ID:   f.ID(),
		Name: f.Name(),
		Desc: f.Description(),
	}, nil
}

func (f Fs) HTTP() chi.Router {
	r := chi.NewRouter()
	for _, mw := range f.Middlewares() {
		r.Use(mw)
	}
	r.Post("/", f.create)
	r.Delete("/remove", f.remove)
	r.Get("/read", f.read)
	return r
}

func (f Fs) ID() string {
	return uuid.NewString()
}

func (f Fs) Name() string {
	return "filestore"
}

func (f Fs) Pattern() string {
	return "/fs"
}

func (f *Fs) Init() error {
	fstore, err := filestore.New(f.K)
	if err != nil {
		return err
	}
	f.store = fstore

	mngr, err := logging.NewSvcManager(f.K, f)
	if err != nil {
		return err
	}
	f.l = mngr
	f.m = metrics.NewLocaler()
	return nil
}

func (f Fs) Shutdown() error {
	f.l.Shutdown()
	return nil
}

func (f Fs) Metrics() []prometheus.Collector {
	f.m.AddCounter("req", metrics.IncomingHTTPRequest(&f))
	f.m.AddCounter("err", metrics.NumberOfErrors(&f))
	return f.m.All()
}

func (f Fs) Description() string {
	return "simple file storage based upon sqlite"
}
