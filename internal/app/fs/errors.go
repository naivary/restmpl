package fs

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
)

var (
	errEmptyFilepath = &jsonapi.ErrorObject{
		ID:     "457be418-e7ee-445c-bf30-e098abe5573a",
		Title:  "missing filepath",
		Detail: "query parameter filepath must be set",
		Status: strconv.Itoa(http.StatusBadRequest),
		Code:   "91abb73e",
	}
)
