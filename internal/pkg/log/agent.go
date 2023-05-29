package log

import (
	"io"

	"github.com/naivary/instance/internal/pkg/service"
	"golang.org/x/exp/slog"
)

type Agent interface {
}

var _ Agent = (*agent)(nil)

type agent struct {
	svcs []service.Service
	// locations to write to
	locs map[slog.Leveler]io.Writer
}

func NewAgent() Agent {
	return agent{}
}
