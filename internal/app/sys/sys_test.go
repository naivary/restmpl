package sys

import (
	"log"
	"net/http"
	"net/http/httptest"
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

// sysTest is a test configured Sys struct.
// Only use for test purposes.
var sysTest = setup()

func TestHealth(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/hi", nil)
	w := httptest.NewRecorder()

	sysTest.Health(w, r)

	t.Log(w.Body)
}
