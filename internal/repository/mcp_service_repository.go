package repository

import (
	"graph-med/internal/base/data"
	"graph-med/internal/model"
	"strings"
)

type McpRepository struct {
	data *data.Data
}

func NewMcpRepository(data *data.Data) *McpRepository {
	return &McpRepository{data: data}
}

// GetMcpServerByIds 获取MCP服务
func (m *McpRepository) GetMcpServerByIds(ids []string) ([]*model.McpService, error) {
	var mcpServers []*model.McpService
	err := m.data.DB.In("mcp_id", ids).Find(&mcpServers)
	return mcpServers, err
}

// GetMcpServerBySessionID 根据会话ID获取MCP ID
func (m *McpRepository) GetMcpServerBySessionID(sessionID string) ([]*model.McpService, error) {
	// 首先从chat_session表获取type_id
	var typeID struct {
		TypeID string `xorm:"type_id"`
	}
	exists, err := m.data.DB.Table("chat_session").Select("type_id").Where("session_id =?", sessionID).Get(&typeID)
	if err != nil || !exists {
		return nil, err
	}

	// 然后从chat_type表获取mcp_ids
	var mcpIDs struct {
		McpIDs string `xorm:"mcp_ids"`
	}
	exists, err = m.data.DB.Table("chat_type").Select("mcp_ids").Where("type_id =?", typeID.TypeID).Get(&mcpIDs)
	if err != nil || !exists {
		return nil, err
	}

	// 假设mcp_ids字段可能包含多个ID（逗号分隔）
	mcpIDList := strings.Split(mcpIDs.McpIDs, ",")

	// 通过mcp_id查询mcp_service表
	res := make([]*model.McpService, 0)
	err = m.data.DB.In("mcp_id", mcpIDList).Find(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
