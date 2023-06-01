package builder

import (
	"context"

	"golang.org/x/exp/slog"
)

type Recorder interface {
	Data() (context.Context, slog.Record)
}
