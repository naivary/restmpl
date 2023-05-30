package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/service"
	"golang.org/x/exp/slog"
)

// Should all log to the same file or seperate?
// If so: Every service should provide its own file to which it logs
// in the schema <name>_<id>.log
type Manager interface {
	Log(context.Context, string, slog.Level, ...slog.Attr)
	Init() error
}

var _ Manager = (*manager)(nil)

type manager struct {
	isInited bool

	k      *koanf.Koanf
	svc    service.Service
	w      io.Writer
	logger *slog.Logger
}

func New(k *koanf.Koanf, svc service.Service) Manager {
	return &manager{
		k: k,
	}
}

// TODO(naivary): make slog.Attr to group and private so it has to use a builder
func (m manager) Log(ctx context.Context, msg string, level slog.Level, attrs ...slog.Attr) {
	m.logger.LogAttrs(ctx, level, msg, attrs...)
}

func (m *manager) Init() error {
	if m.isInited {
		return nil
	}
	filename := fmt.Sprintf("%s_%s.log", m.svc.Name(), m.svc.ID())
	p := filepath.Join(m.k.String("logsDir"), filename)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	m.logger = slog.New(slog.NewTextHandler(file, nil))
	m.w = file
	return nil
}
