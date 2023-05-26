package logging

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/naivary/instance/internal/pkg/config"
	"golang.org/x/exp/slog"
)

const (
	cfgFile = "instance.yaml"
)

func testLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, nil))
}

func setupLogging() Logging {
	l := Logging{}
	k, err := config.New(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	l.K = k
	return l
}

func TestDefaultLogger(t *testing.T) {
	l := setupLogging()
	r := httptest.NewRequest(http.MethodGet, "/sys/health", nil)
	w := httptest.NewRecorder()
	l.Info = testLogger(w)
	empty := func(w http.ResponseWriter, r *http.Request) {}
	hf := l.Logger(http.HandlerFunc(empty))
	id := middleware.RequestID(hf)
	id.ServeHTTP(w, r)
	t.Log(w.Body)
}
