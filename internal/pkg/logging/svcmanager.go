package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/apitmpl/internal/pkg/logging/builder"
	"github.com/naivary/apitmpl/internal/pkg/service"
	"golang.org/x/exp/slog"
)

var _ Manager = (*svcManager)(nil)

type svcManager struct {
	// size at which a rotate should be initiliazed
	maxSize int64
	// maximum number of backups to create
	maxBackups int
	// if the backups should be compressed
	compress bool

	// created backups
	backups []*os.File
	k       *koanf.Koanf
	svc     service.Service
	// file to which logger writes. It does not hold the content of the
	// current log file. This should be
	file   *os.File
	logger *slog.Logger
	stream chan builder.Recorder
}

func NewSvcManager(k *koanf.Koanf, svc service.Service) (Manager, error) {
	m := &svcManager{
		k:          k,
		svc:        svc,
		stream:     make(chan builder.Recorder, 1),
		maxSize:    k.Int64("logs.maxSize"),
		maxBackups: k.Int("logs.maxBackups"),
		compress:   k.Bool("logs.compress"),
	}
	if err := m.init(); err != nil {
		return nil, err
	}
	m.backups = make([]*os.File, m.maxBackups)
	return m, nil
}

func (m svcManager) Log(r builder.Recorder) {
	m.stream <- r
}

func (m svcManager) Shutdown() {
	m.file.Close()
	close(m.stream)
}

func (m *svcManager) init() error {
	filename := fmt.Sprintf("%s_%s.log", m.svc.Name(), m.svc.ID())
	p := filepath.Join(m.k.String("logsDir"), filename)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	m.file = file
	m.logger = slog.New(slog.NewJSONHandler(m.file, nil)).With(m.commonAttrs()...)
	go m.handle()
	return nil
}

func (m *svcManager) handle() {
	for record := range m.stream {
		ctx, rec := record.Data()
		if err := m.logger.Handler().Handle(ctx, rec); err != nil {
			slog.Error("could not write", slog.String("err", err.Error()))
			return
		}
		if err := m.rotate(); err != nil {
			slog.Error("could not rotate", slog.String("err", err.Error()))
			return
		}
	}
}

func (m svcManager) commonAttrs() []any {
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
