package fs

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/naivary/instance/internal/pkg/config/configtest"
	"github.com/naivary/instance/internal/pkg/dependency"
	"github.com/naivary/instance/internal/pkg/env"
	"github.com/naivary/instance/internal/pkg/filestore/filestoretest"
	"github.com/naivary/instance/internal/pkg/must"
	"github.com/naivary/instance/internal/pkg/service"
)

var (
	fsTest, ts = setup()
)

func setup() (Fs, *httptest.Server) {
	f := Fs{}
	k, err := configtest.New()
	if err != nil {
		log.Fatal(err)
	}
	f.K = k

	st, err := filestoretest.New(k)
	if err != nil {
		log.Fatal(err)
	}
	f.Store = st
	kDep := dependency.New(nil, k)
	stDep := dependency.New(nil, st)
	e := env.NewAPI([]service.Service{f}, k, nil, []dependency.Pinger{kDep, stDep})
	return f, httptest.NewServer(e.HTTP())
}

func upload(client *http.Client, url string, params map[string]string, formKey string, path string) (*http.Response, error) {
	file := must.Open(path)
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	form := new(bytes.Buffer)
	w := multipart.NewWriter(form)

	pdfFile, err := w.CreateFormFile(formKey, info.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(pdfFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		w.WriteField(key, val)
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return client.Post(url, w.FormDataContentType(), form)
}

func TestCreate(t *testing.T) {
	c := ts.Client()
	u, err := url.JoinPath(ts.URL, "fs")
	if err != nil {
		t.Error(err)
	}

	params := map[string]string{
		"filepath": "pdf/",
	}

	res, err := upload(c, u, params, fsTest.K.String("fs.formKey"), "./testdata/dummy.pdf")
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code to be %d. Got: %d while sending request to: %s", http.StatusCreated, res.StatusCode, u)
	}

	read(t)
	remove(t)

}

func remove(t *testing.T) {
	c := ts.Client()
	u, err := url.JoinPath(ts.URL, "fs", "remove")
	if err != nil {
		t.Error(err)
	}
	val := url.Values{}
	val.Add("filepath", "pdf/dummy.pdf")

	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		t.Error(err)
	}

	req.URL.RawQuery = val.Encode()

	res, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code to be %d. Got: %d while sending request to: %s", http.StatusNoContent, res.StatusCode, u)
	}
}

func read(t *testing.T) {
	c := ts.Client()
	file := must.Open("./testdata/dummy.pdf")
	u, err := url.JoinPath(ts.URL, "fs", "read")
	if err != nil {
		t.Error(err)
	}
	val := url.Values{}
	val.Add("filepath", "pdf/dummy.pdf")

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		t.Error(err)
	}
	req.URL.RawQuery = val.Encode()
	res, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code to be %d. Got: %d while sending request to: %s", http.StatusOK, res.StatusCode, u)
	}

	got := new(bytes.Buffer)
	expected := new(bytes.Buffer)

	_, err = expected.ReadFrom(file)
	if err != nil {
		t.Error(err)
	}

	_, err = got.ReadFrom(res.Body)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(got.Bytes(), expected.Bytes()) {
		t.Fatalf("Expected too be equal")
	}

}
