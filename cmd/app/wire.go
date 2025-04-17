//go:generate wire
//go:build wireinject
// +build wireinject

package main

import (
	"graph-med/internal/base/casbin"
	"graph-med/internal/base/conf"
	"graph-med/internal/base/data"
	"graph-med/internal/controller"
	"graph-med/internal/repository"
	"graph-med/internal/router"
	"graph-med/internal/server"
	"graph-med/internal/service"

	"github.com/google/wire"
)

func initializeServer(config *conf.AllConfig) (*server.Server, error) {
	panic(wire.Build(
		data.NewDB,
		data.NewData,
		data.NewNeo4jDriver,
		casbin.NewEnforcer,
		repository.ProviderSetRepository,
		service.ProviderSetService,
		controller.ProviderSetController,
		router.ProviderSetRouter,
		server.ProviderSetServer,
	))
}
