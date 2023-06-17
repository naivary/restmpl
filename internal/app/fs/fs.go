package fs

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/objst"
	"github.com/naivary/restmpl/internal/pkg/jwtauth"
	"github.com/naivary/restmpl/internal/pkg/logging"
	"github.com/naivary/restmpl/internal/pkg/metrics"
	"github.com/naivary/restmpl/internal/pkg/service"
	"github.com/pocketbase/dbx"
	"github.com/prometheus/client_golang/prometheus"
)

var _ service.Service = (*Fs)(nil)

type Fs struct {
	K  *koanf.Koanf
	DB *dbx.DB

	l       logging.Manager
	m       metrics.Managee
	b       *objst.Bucket
	maxSize int64
	formKey string
}

func (f *Fs) Init() error {
	f.l = logging.NewSvcManager(f)
	f.m = metrics.NewManagee()
	opts := badger.DefaultOptions("/tmp/badger/store")
	b, err := objst.NewBucket(&opts)
	if err != nil {
		return err
	}
	f.b = b
	f.maxSize = f.K.Int64("fs.maxSize")
	f.formKey = f.K.String("fs.formKey")
	return nil
}

func (f Fs) Health() (*service.Info, error) {
	if f.l == nil {
		return nil, errors.New("missing loggin manager")
	}
	if f.m == nil {
		return nil, errors.New("missing metrics managee")
	}
	if f.b == nil {
		return nil, errors.New("bucket is missing")
	}
	if err := f.b.Health(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (f Fs) HTTP() chi.Router {
	r := chi.NewRouter()
	r.Use(jwtauth.Verify)
	r.Get("/", f.read)
	r.Post("/", f.create)
	r.Delete("/{id}", f.remove)
	return r
}

func (f Fs) Name() string {
	return "filestore"
}

func (f Fs) ID() string {
	return uuid.NewString()
}

func (f Fs) Description() string {
	return "simple embedded filestore to store small to medium size files"
}

func (f Fs) Metrics() []prometheus.Collector {
	return f.m.All()
}

func (f Fs) Pattern() string {
	return "/fs"
}

func (f Fs) Shutdown() error {
	return f.b.Shutdown()
}
