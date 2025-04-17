package repository

import "github.com/google/wire"

var ProviderSetRepository = wire.NewSet(
	NewUserRepository,
	NewChatRepository,
	NewNodeRepository,
	NewKGRepository,
	NewMcpRepository,
)
