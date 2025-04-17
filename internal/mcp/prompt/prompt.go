package prompt

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPPrompt interface {
	Name() string
	Definition() mcp.Prompt
	ToolHandlerFunc(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
}
