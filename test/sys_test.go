package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/sys", nil)
	w := httptest.NewRecorder()
	hf := http.HandlerFunc(api.Services.Sys.Health)
	hf.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Errorf("statuscode was not 200. Got: %d", w.Result().StatusCode)
	}
}
