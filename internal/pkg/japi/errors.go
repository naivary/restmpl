package japi

import (
	"strconv"

	"github.com/google/jsonapi"
)

func Errors(errs ...*jsonapi.ErrorObject) []*jsonapi.ErrorObject {
	return errs
}

func NewError(err error, status int, reqID string) jsonapi.ErrorObject {
	return jsonapi.ErrorObject{
		Status: strconv.Itoa(status),
		Title:  "uknown",
		Detail: err.Error(),
		Meta: &map[string]interface{}{
			"reqID": reqID,
		},
	}
}
