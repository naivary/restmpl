//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/naivary/asive/internal/app/appointment"
	"github.com/naivary/asive/internal/app/client"
	"github.com/naivary/asive/internal/app/lesson"
	"github.com/naivary/asive/internal/app/school"
	"github.com/naivary/asive/internal/pkg/ctrl"
	"github.com/naivary/asive/internal/pkg/database"
	"github.com/naivary/asive/internal/pkg/models/dependency"
)

var dependencies = wire.NewSet(database.Connect)
var views = wire.NewSet(wire.Struct(new(client.Environment), "*"), wire.Struct(new(appointment.Environment), "*"), wire.Struct(new(school.Environment), "*"), wire.Struct(new(lesson.Environment), "*"), wire.Struct(new(ctrl.Views), "*"))
var depenciesStruct = wire.NewSet(wire.Struct(new(dependency.Default), "*"))

func StartApp() (*ctrl.App, error) {
	wire.Build(dependencies, views, ctrl.SetEndpoints, depenciesStruct, wire.Struct(new(ctrl.App), "*"))
	return &ctrl.App{}, nil
}
