package service

import (
	"fmt"
	"graph-med/internal/base/conf"
	"graph-med/internal/repository"
	"graph-med/internal/schema"
	"os"
	"path/filepath"
)

type KGService struct {
	kgRepo       *repository.KGRepository
	nodeRepo     *repository.NodeRepository
	medicalAttrs map[string]string
}

func NewKGService(kgRepo *repository.KGRepository, nodeRepo *repository.NodeRepository) *KGService {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前路径失败:", err)
		panic(err)
	}

	medicalAttrs, err := conf.LoadMedicalAttrs(filepath.Join(currentDir, conf.DefaultMedicalAttrPath))
	if err != nil {
		fmt.Println("加载疾病属性失败:", err)
		panic(err)
	}

	return &KGService{
		kgRepo:       kgRepo,
		nodeRepo:     nodeRepo,
		medicalAttrs: medicalAttrs,
	}
}

// GetLabels 获取疾病库所有标签
func (s *KGService) GetLabels() ([]string, error) {
	return s.kgRepo.GetLabels()
}

// GetNodes 获取疾病节点
func (s *KGService) GetNodes(label, name string, limit int) ([]schema.KGNode, error) {
	nodes, err := s.kgRepo.SearchNodes(label, name, limit)
	if err != nil {
		return nil, err
	}

	var list []schema.KGNode
	for _, node := range nodes {
		list = append(list, schema.KGNode{
			ID:    node.ElementId,
			Label: node.Labels[0],
			Name:  node.Props["name"].(string),
		})
	}

	return list, nil
}

// GetNodeProperties 获取节点属性
func (s *KGService) GetNodeProperties(label, name string, attrs []string) (map[string]string, error) {
	nodeInfos, err := s.nodeRepo.GetNodeProperties(label, name, attrs)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]string)
	for _, nodeInfo := range nodeInfos {
		properties[nodeInfo.AttrName] = nodeInfo.AttrValue
	}
	return properties, nil
}

// GetNodeDetail 获取节点详情
func (s *KGService) GetNodeDetail(label, name string) (map[string]string, error) {
	nodeInfos, err := s.nodeRepo.GetNodeAllProperties(label, name)
	if err != nil {
		return nil, err
	}

	properties := make(map[string]string)
	for _, nodeInfo := range nodeInfos {
		properties[nodeInfo.AttrName] = nodeInfo.AttrValue
	}
	return properties, nil
}

// GetSubgraph 获取节点子图
func (s *KGService) GetSubgraph(label, name string, limit int) (*schema.KGGraph, error) {
	nodeList, relaList, err := s.kgRepo.GetNodeSubGraph(label, name, limit)
	if err != nil {
		return nil, err
	}

	subgraph := &schema.KGGraph{
		Nodes: make([]schema.KGNode, len(nodeList)),
		Edges: make([]schema.KGEdge, len(relaList)),
	}

	for i, node := range nodeList {
		subgraph.Nodes[i] = schema.KGNode{
			ID:       node.ElementId,
			Label:    node.Labels[0],
			Name:     node.Props["name"].(string),
			NodeType: node.Labels[0],
		}
	}

	for i, rela := range relaList {
		subgraph.Edges[i] = schema.KGEdge{
			ID:     rela.ElementId,
			Target: rela.StartElementId,
			Source: rela.EndElementId,
			Type:   rela.Type,
			Label:  rela.Type,
		}
	}

	return subgraph, nil
}

// GetMedicalAttrs 获取疾病属性
func (s *KGService) GetMedicalAttrs() map[string]string {
	return s.medicalAttrs
}
