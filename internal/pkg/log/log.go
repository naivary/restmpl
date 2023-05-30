package log

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"golang.org/x/exp/slog"
)

// Should all log to the same file or seperate?
// If so: Every service should provide its own file to which it logs
// in the schema <name>_<id>.log
type Manager interface {
	Log(context.Context, string, slog.Level, ...slog.Attr)
	AddLevel(slog.Leveler, *slog.Logger) error
	Init() error
}

var _ Manager = (*manager)(nil)

type manager struct {
	k        *koanf.Koanf
	loggers  map[slog.Leveler]*slog.Logger
	defaults map[slog.Leveler]string
}

func New(k *koanf.Koanf) Manager {
	return &manager{
		k: k,
		defaults: map[slog.Leveler]string{
			slog.LevelInfo:  "info.log",
			slog.LevelError: "error.log",
		},
	}
}

func (m manager) AddLevel(level slog.Leveler, l *slog.Logger) error {
	_, ok := m.loggers[level]
	if ok {
		return errors.New("level already exists")
	}
	m.loggers[level] = l
	return nil
}

func (m manager) Log(ctx context.Context, msg string, level slog.Level, attrs ...slog.Attr) {
	logger := m.loggers[level]
	logger.LogAttrs(ctx, level, msg, attrs...)
}

func (m manager) initLogger(level slog.Leveler, w io.Writer) {
	l := slog.New(slog.NewTextHandler(w, nil))
	m.loggers[level] = l
}

func (m manager) Init() error {
	for level, filename := range m.defaults {
		p := filepath.Join(m.k.String("logsDir"), filename)
		file, err := os.Create(p)
		if err != nil {
			return err
		}
		m.initLogger(level, file)
	}
	return nil
}
