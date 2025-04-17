package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"graph-med/internal/service"
)

type DiseaseSubgraphTool struct {
	kgService *service.KGService
}

func NewDiseaseSubgraphTool(kgService *service.KGService) *DiseaseSubgraphTool {
	return &DiseaseSubgraphTool{
		kgService: kgService,
	}
}

func (t *DiseaseSubgraphTool) Name() string {
	return "get_disease_subgraph"
}

func (t *DiseaseSubgraphTool) Definition() mcp.Tool {
	return mcp.NewTool(t.Name(),
		mcp.WithDescription("当你想查询指定疾病的相关知识图谱时，输入疾病名称，返回该疾病的知识图谱。"),
		mcp.WithString(
			"disease_name",
			mcp.Required(),
			mcp.Description("疾病名称"),
		),
	)
}

func (t *DiseaseSubgraphTool) ToolHandlerFunc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["disease_name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}

	subgraph, err := t.kgService.GetSubgraph("疾病", name, 100)
	if err != nil {
		return nil, err
	}

	if subgraph == nil {
		return mcp.NewToolResultText(fmt.Sprintf("没有找到疾病 %s 的相关知识图谱", name)), nil
	}

	result, err := json.Marshal(subgraph)
	if err != nil {
		return nil, err
	}

	//result := ""
	//idToName := make(map[string]string)
	//
	//for _, node := range subgraph.Nodes {
	//	idToName[node.ID] = node.Name
	//}
	//
	//for _, edge := range subgraph.Edges {
	//	result += fmt.Sprintf("%s ==%s==> %s\n\n", idToName[edge.Target], edge.Type, idToName[edge.Source])
	//}

	return mcp.NewToolResultText(string(result)), nil
}
