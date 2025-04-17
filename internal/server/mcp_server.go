package server

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"graph-med/internal/base/conf"
	"graph-med/internal/mcp/prompt"
	"graph-med/internal/mcp/tools"
)

type McpServer struct {
	tools   []tools.McpTool
	prompts []prompt.MCPPrompt
	addr    string
	*server.MCPServer
	*server.SSEServer
}

// NewMCPServer 创建mcp服务
func NewMCPServer(config *conf.AllConfig, diseaseInfoTool *tools.DiseaseSubgraphTool) *McpServer {
	mcpServer := server.NewMCPServer(
		"knowledge_graph_assistant",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	sseServer := server.NewSSEServer(mcpServer)

	mcpTools := []tools.McpTool{
		diseaseInfoTool,
	}
	mcpPrompts := []prompt.MCPPrompt{}

	mcpConfig := config.Mcp
	addr := fmt.Sprintf(":%s", mcpConfig.Port)

	s := &McpServer{
		tools:     mcpTools,
		prompts:   mcpPrompts,
		addr:      addr,
		MCPServer: mcpServer,
		SSEServer: sseServer,
	}

	s.setupTools()
	return s
}

// Run 启动mcp服务
func (s *McpServer) Run() error {
	fmt.Println("Mcp Server is successfully started, tool count: ", len(s.tools))
	fmt.Println("Mcp Server serving on ", s.addr)
	if err := s.Start(s.addr); err != nil {
		return fmt.Errorf("mcp serve: %w", err)
	}
	return nil
}

// setupTools 设置mcp服务
func (s *McpServer) setupTools() {
	fmt.Println()
	for _, tool := range s.tools {
		fmt.Printf("Loading tool  ==>  %s\n", tool.Name())
		s.AddTool(tool.Definition(), tool.ToolHandlerFunc)
	}

	for _, prom := range s.prompts {
		fmt.Printf("Loading prompt  ==>  %s\n", prom.Name())
		s.AddPrompt(prom.Definition(), prom.ToolHandlerFunc)
	}
	fmt.Println()
}
