package logging

import (
	"os"

	"github.com/naivary/restmpl/internal/pkg/logging/builder"
	"golang.org/x/exp/slog"
)

type Manager interface {
	Log(builder.Recorder) error
	AddCommonAttrs(attrs ...any)
}

var _ Manager = (*manager)(nil)

type manager struct {
	logger *slog.Logger
}

func NewManager() *manager {
	return &manager{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (m manager) Log(record builder.Recorder) error {
	return m.logger.Handler().Handle(record.Data())
}

func (m *manager) AddCommonAttrs(attrs ...any) {
	m.logger = m.logger.With(attrs...)
}
