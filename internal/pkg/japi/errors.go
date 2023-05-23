package japi

import "github.com/google/jsonapi"

func Errors(errs ...*jsonapi.ErrorObject) []*jsonapi.ErrorObject {
	return errs
}
