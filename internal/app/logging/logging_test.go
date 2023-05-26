package logging

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

	// check if host key is set to example.com
	buf := new(bytes.Buffer)
	buf.ReadFrom(w.Body)
	attrs := strings.Split(buf.String(), " ")
	for _, attr := range attrs {
		key, value, ok := strings.Cut(attr, "=")
		if !ok {
			t.Fatalf("cutting went wrong. Key:%s, Value:%s", key, value)
		}
		if key == "host" && value == "example.com" {
			return
		}
	}
}
