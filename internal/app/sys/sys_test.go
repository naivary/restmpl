package sys

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/naivary/instance/internal/pkg/config/configtest"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/must"
	"github.com/naivary/instance/internal/pkg/service"
	"github.com/naivary/instance/internal/pkg/testutil"
)

var (
	// sysTest is a test configured Sys struct
	// which is only used for test purposes
	sysTest = setupSys()

	ts = setupTs()
)

func setupSys() Sys {
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
	s.M = metadata.New(k, db)
	return s
}

func setupTs() *httptest.Server {
	api := env.NewAPI([]service.Service[chi.Router]{sysTest}, sysTest.K)
	root := api.Router()
	return httptest.NewServer(root)
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

	got := new(bytes.Buffer)
	expected := new(bytes.Buffer)

	_, err = got.ReadFrom(res.Body)
	if err != nil {
		t.Error(err)
	}

	file := must.Open("./testdata/health.json")
	_, err = expected.ReadFrom(file)
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	err = json.Compact(buf, expected.Bytes())
	if err != nil {
		t.Error(err)
	}

	if ok, err := testutil.AreEqualJSON(buf.String(), got.String()); !ok || err != nil {
		t.Log(buf.String(), got.String())
		t.Fatalf("Should be equal: %v", err)
	}
}
