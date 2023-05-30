package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/log/builder"
	"github.com/naivary/instance/internal/pkg/service"
	"golang.org/x/exp/slog"
)

// Should all log to the same file or seperate?
// If so: Every service should provide its own file to which it logs
// in the schema <name>_<id>.log
type Manager interface {
	Log(builder.Recorder)
	// Shutdown() error
}

var _ Manager = (*manager)(nil)

type manager struct {
	isInited   bool
	maxSize    uint
	maxAge     uint
	maxBackups uint
	filename   string
	compress   bool

	k      *koanf.Koanf
	svc    service.Service
	w      io.Writer
	logger *slog.Logger
	ch     chan builder.Recorder
}

func New(k *koanf.Koanf, svc service.Service) (Manager, error) {
	m := &manager{
		k:   k,
		svc: svc,
		ch:  make(chan builder.Recorder, 2),
	}
	if err := m.init(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m manager) Log(record builder.Recorder) {
	m.ch <- record
}

func (m manager) handle() {
	for record := range m.ch {
		ctx, rec := record.Data()
		m.logger.Handler().Handle(ctx, rec)
	}
}

func (m *manager) init() error {
	filename := fmt.Sprintf("%s_%s.log", m.svc.Name(), m.svc.ID())
	p := filepath.Join(m.k.String("logsDir"), filename)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	m.logger = slog.New(slog.NewTextHandler(file, nil))
	m.w = file
	// IDK if this is good or not
	go m.handle()
	return nil
}
