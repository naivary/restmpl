package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/sys/health", nil)
	w := httptest.NewRecorder()
	api.Router.ServeHTTP(w, r)
	fmt.Println(w.Body)
	if w.Result().StatusCode != 200 {
		t.Errorf("statuscode was not 200. Got: %d", w.Result().StatusCode)
	}
}
