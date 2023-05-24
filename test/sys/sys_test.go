package sys

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naivary/instance/internal/pkg/ctrl"
)

func TestHealth(t *testing.T) {
	api, err := ctrl.New()
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/sys/health", nil)
	w := httptest.NewRecorder()
	api.Services.Sys.Health(w, r)

	res := w.Result()
	defer res.Body.Close()

	t.Log(res)

}
