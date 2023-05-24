package test

import (
	"log"

	"github.com/naivary/instance/internal/pkg/ctrl"
)

var api *ctrl.API

func init() {
	app, err := ctrl.New()
	if err != nil {
		log.Fatal(err)
	}
	api = app
}
