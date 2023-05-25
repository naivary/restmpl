package sys

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/models/metadata"
	"github.com/naivary/instance/internal/pkg/routes"
)

const (
	cfgFile = "test_config.yaml"
)

var (
	// sysTest is a test configured Sys struct
	// which is only used for test purposes
	sysTest = setupSys()

	ts = setupTestServer()
)

func setupSys() Sys {
	s := Sys{}
	k, err := config.New(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	s.K = k

	db, err := database.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}
	s.DB = db

	s.M = metadata.New(k)

	return s
}

func setupTestServer() *httptest.Server {
	root := routes.NewTestRouter()
	sysTest.Register(root)
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
}
