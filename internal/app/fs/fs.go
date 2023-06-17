package fs

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/logging"
	"github.com/naivary/apitmpl/internal/pkg/metrics"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"github.com/prometheus/client_golang/prometheus"
)

var _ service.Service = (*Fs)(nil)

type Fs struct {
	K *koanf.Koanf

	l logging.Manager
	m metrics.Managee
}

func (f Fs) Health() (*service.Info, error) {
	if f.K == nil {
		return nil, errors.New("missing config manager")
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
	lMngr, err := logging.NewSvcManager(f.K, f)
	if err != nil {
		return err
	}
	f.l = lMngr
	f.m = metrics.NewManagee()
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
