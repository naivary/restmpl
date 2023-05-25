package sys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
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

func mustOpen(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func areEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
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

	file := mustOpen("./testdata/health.json")
	_, err = expected.ReadFrom(file)
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	err = json.Compact(buf, expected.Bytes())
	if err != nil {
		t.Error(err)
	}

	if ok, err := areEqualJSON(buf.String(), got.String()); !ok || err != nil {
		t.Fatalf("Should be equal: %v", err)
	}
}
