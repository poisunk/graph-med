package tools

import "github.com/google/wire"

var ProviderSetMcpTool = wire.NewSet(
	NewDiseaseSubgraphTool,
)
