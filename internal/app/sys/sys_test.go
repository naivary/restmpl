package sys

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/naivary/instance/internal/pkg/config"
	"github.com/naivary/instance/internal/pkg/database"
	"github.com/naivary/instance/internal/pkg/models/metadata"
)

func setup() Sys {
	s := Sys{}
	db, err := database.InMemConnect()
	if err != nil {
		log.Fatal(err)
	}
	s.DB = db
	k, err := config.New("test_config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	s.K = k
	s.M = metadata.New(k)
	return s
}

var (
	// sysTest is a test configured Sys struct.
	// Only use for test purposes.
	sysTest = setup()
)

func TestHealth(t *testing.T) {
	ts := httptest.NewServer(sysTest.router())
	defer ts.Close()
	url, err := url.JoinPath(ts.URL, "/health")
	if err != nil {
		t.Error(err)
	}
	res, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	fmt.Println(buf)
	t.Log(res)
}
