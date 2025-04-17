package prompt

import "github.com/google/wire"

var ProviderSetMcpPrompt = wire.NewSet(
	NewUserInquirePrompt,
)
