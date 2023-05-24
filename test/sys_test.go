package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naivary/instance/internal/pkg/models/metadata"
)

func TestHealth(t *testing.T) {
	m := metadata.Metadata{}
	r := httptest.NewRequest(http.MethodGet, "/sys/health", nil)
	w := httptest.NewRecorder()
	api.Services.Sys.Health(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fail()
	}

	err := json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		t.Error(err)
	}
}
