package filestore

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
)

var (
	ErrWrongNaming = jsonapi.ErrorObject{
		ID:     "50493051-a153-4430-bbf7-b9b40bece960",
		Title:  "wrong name convention",
		Detail: "file name must follow the following regex pattern: [a-z._0-9-]+",
		Status: strconv.Itoa(http.StatusBadRequest),
		Code:   "bc138545",
	}
)
