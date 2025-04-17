package repository

import (
	"graph-med/internal/base/data"
	"graph-med/internal/model"
)

type NodeRepository struct {
	data *data.Data
}

func NewNodeRepository(data *data.Data) *NodeRepository {
	return &NodeRepository{
		data: data,
	}
}

// GetNodeProperties 获取节点属性
func (n *NodeRepository) GetNodeProperties(label, name string, attrs []string) ([]*model.NodeInfo, error) {
	var list []*model.NodeInfo
	err := n.data.DB.Where("label = ? AND name = ?", label, name).In("attr_name", attrs).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// GetNodeAllProperties 获取节点所有属性
func (n *NodeRepository) GetNodeAllProperties(label, name string) ([]*model.NodeInfo, error) {
	var list []*model.NodeInfo
	err := n.data.DB.Where("label = ? AND name = ?", label, name).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
