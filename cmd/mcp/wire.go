//go:generate wire
//go:build wireinject
// +build wireinject

package main

import (
	"graph-med/internal/base/conf"
	"graph-med/internal/base/data"
	"graph-med/internal/mcp/tools"
	"graph-med/internal/repository"
	"graph-med/internal/server"
	"graph-med/internal/service"

	"github.com/google/wire"
)

func initializeMcpServer(config *conf.AllConfig) (*server.McpServer, error) {
	panic(wire.Build(
		data.NewDB,
		data.NewData,
		data.NewNeo4jDriver,
		server.NewMCPServer,
		tools.ProviderSetMcpTool,
		service.ProviderSetService,
		repository.ProviderSetRepository,
	))
}
