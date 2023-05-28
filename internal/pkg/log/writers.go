package log

import (
	"io"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/v2"
	"golang.org/x/exp/slog"
)

var (
	defaults = map[slog.Leveler]string{
		slog.LevelInfo:  "info.log",
		slog.LevelError: "error.log",
	}
)

type LogWriter interface {
	Get(slog.Leveler) (io.Writer, bool)
	Add(slog.Leveler, io.Writer)
}

// Writers persist the location to which a
// slog.Leveler should be writing.
type writers map[slog.Leveler]io.Writer

func New(k *koanf.Koanf) LogWriter {
	w := &writers{}
	for level, filename := range defaults {
		p := filepath.Join(k.String("logsDir"), filename)
		file, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			slog.Error("cant open log file", "err", err.Error(), "filename", filename)
			return nil
		}
		w.Add(level, file)
	}
	return w
}

func (w writers) Add(level slog.Leveler, wr io.Writer) {
	w[level] = wr
}

func (w writers) Get(level slog.Leveler) (io.Writer, bool) {
	wr, ok := w[level]
	return wr, ok
}
