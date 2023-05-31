package log

import (
	"compress/gzip"
	"fmt"
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
	Shutdown()
}

var _ Manager = (*manager)(nil)

type manager struct {
	// size at which a rotate should be initiliazed
	maxSize uint64
	// maxAge to keep ol backups
	maxAge uint
	// maximum number of backups to create
	maxBackups uint
	// if the backups should be compressed
	compress bool

	k      *koanf.Koanf
	svc    service.Service
	file   *os.File
	gzipw  *gzip.Writer
	logger *slog.Logger
	stream chan builder.Recorder
}

func New(k *koanf.Koanf, svc service.Service) (Manager, error) {
	m := &manager{
		k:      k,
		svc:    svc,
		stream: make(chan builder.Recorder, 2),
	}
	if err := m.init(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m manager) Log(record builder.Recorder) {
	m.stream <- record
}

// shutdown
func (m manager) Shutdown() {
	m.file.Close()
	close(m.stream)
}

func (m manager) write() error {
	for record := range m.stream {
		ctx, rec := record.Data()
		m.logger.Handler().Handle(ctx, rec)
		err := m.gzipw.Flush()
		if err != nil {
			return err
		}
		if err := m.rotate(); err != nil {
			return err
		}
	}
	return nil
}

func (m manager) handle() {
	if err := m.write(); err != nil {
		m.Shutdown()
	}
}

func (m *manager) init() error {
	filename := fmt.Sprintf("%s_%s.gz", m.svc.Name(), m.svc.ID())
	p := filepath.Join(m.k.String("logsDir"), filename)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	m.gzipw = gzip.NewWriter(file)
	m.file = file
	m.logger = slog.New(slog.NewTextHandler(m.gzipw, nil)).With(m.commonAttrs()...)
	go m.handle()
	return nil
}

func (m manager) commonAttrs() []any {
	svc := slog.Group(
		"service",
		slog.String("id", m.svc.ID()),
		slog.String("name", m.svc.Name()),
	)
	api := slog.Group(
		"api",
		slog.String("version", m.k.String("version")),
	)
	return []any{api, svc}
}
