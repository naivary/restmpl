package filestore

import (
	"errors"
)

var (
	ErrWrongNaming = errors.New("file name must follow the followin regex pattern: [a-z._0-9-]+")
)
