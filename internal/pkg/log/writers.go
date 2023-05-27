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
// slog.Leveler should be writing. So the
// common pattern would be that a Writers
// would be injected into every Service which
// needs logging. Based upon the slog.Leveler
// which is needed the right io.Writer will be
// delivered. Writers implementes the LogWriter
// interface.
type Writers map[slog.Leveler]io.Writer

func NewWriters(k *koanf.Koanf) LogWriter {
	w := &Writers{}
	for level, filename := range defaults {
		p := filepath.Join(k.String("logsDir"), filename)
		file, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			slog.Error("cant open the log file", err.Error())
			return nil
		}
		w.Add(level, file)
	}
	return w
}

func (w Writers) Add(level slog.Leveler, wr io.Writer) {
	w[level] = wr
}

func (w Writers) Get(level slog.Leveler) (io.Writer, bool) {
	wr, ok := w[level]
	return wr, ok
}
