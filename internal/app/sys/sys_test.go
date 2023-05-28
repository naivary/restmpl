package sys

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"github.com/naivary/instance/internal/pkg/config/configtest"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/service"
)

var (
	ts = setup()
)

func setup() *httptest.Server {
	s := Sys{}
	k, err := configtest.New()
	if err != nil {
		log.Fatal(err)
	}
	s.K = k

	db, err := database.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}
	s.DB = db

	e := env.NewAPI([]service.Service[chi.Router]{s}, k)
	s.M = metadata.New[*koanf.Koanf, chi.Router](k, db, &e)
	return httptest.NewServer(e.Router())
}

func TestHealth(t *testing.T) {
	c := ts.Client()
	url, err := url.JoinPath(ts.URL, "sys", "health")
	if err != nil {
		t.Error(err)
	}
	res, err := c.Get(url)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code to be %d. Got: %d while sending request to: %s", http.StatusOK, res.StatusCode, url)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		t.Error(err)
	}

	if buf.Len() <= 0 {
		t.Fatalf("no response body found in the response. Got: %v", buf.String())
	}
}
