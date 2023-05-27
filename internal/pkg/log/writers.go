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

type Writers map[slog.Leveler]io.Writer

func NewWriters(k *koanf.Koanf) *Writers {
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
