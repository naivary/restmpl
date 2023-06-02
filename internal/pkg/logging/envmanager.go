package logging

import (
	"io"

	"github.com/naivary/apitmpl/internal/pkg/logging/builder"
	"golang.org/x/exp/slog"
)

var _ Manager = (*envManager)(nil)

type envManager struct {
	logger *slog.Logger
	w      io.Writer
}

func NewEnvManager(w io.Writer) *envManager {
	e := &envManager{
		w:      w,
		logger: slog.New(slog.NewTextHandler(w, nil)),
	}
	return e
}

func (e envManager) Shutdown() {}

func (e envManager) Log(r builder.Recorder) {
	ctx, rec := r.Data()
	if err := e.logger.Handler().Handle(ctx, rec); err != nil {
		slog.Error("could not write", slog.String("err", err.Error()))
		return
	}
}
