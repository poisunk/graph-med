package repository

import (
	"context"
	"fmt"
	"graph-med/internal/base/data"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type KGRepository struct {
	data *data.Data
}

func NewKGRepository(data *data.Data) *KGRepository {
	return &KGRepository{
		data: data,
	}
}

// GetLabels 获取疾病库所有标签
func (r *KGRepository) GetLabels() ([]string, error) {
	cypher := `CALL db.labels() YIELD label RETURN label`

	var labels []string
	ctx := context.Background()
	session := r.data.Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{})
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			if value, ok := record.Get("label"); ok {
				label := value.(string)
				labels = append(labels, label)
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}

		return labels, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return labels, err
}

// SearchNodes 搜索节点
func (r *KGRepository) SearchNodes(label, name string, limit int) ([]neo4j.Node, error) {
	cypher := "MATCH (p:`%s` {name:'%s') RETURN p LIMIT %d"

	cypher = fmt.Sprintf(cypher, label, name, limit)

	var list []neo4j.Node
	ctx := context.Background()
	session := r.data.Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{})
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			if value, ok := record.Get("p"); ok {
				node := value.(neo4j.Node)
				list = append(list, node)
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}

		return list, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return list, err
}

// GetNodeSubGraph 获取节点子图
func (r *KGRepository) GetNodeSubGraph(label, name string, limit int) ([]neo4j.Node, []neo4j.Relationship, error) {
	cypher := "MATCH path=(s:`%s` {name: '%s'})-[*..1]-(connected) " +
		"RETURN nodes(path) AS nodes, relationships(path) AS relationships LIMIT %d"

	cypher = fmt.Sprintf(cypher, label, name, limit)

	nodeIdentitySet := make(map[string]struct{})
	var nodeList []neo4j.Node
	var relaList []neo4j.Relationship
	ctx := context.Background()
	session := r.data.Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{})
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			if value, ok := record.Get("nodes"); ok {
				list := value.([]interface{})
				for _, ele := range list {
					node := ele.(neo4j.Node)
					if _, ok := nodeIdentitySet[node.ElementId]; !ok {
						nodeIdentitySet[node.ElementId] = struct{}{}
						nodeList = append(nodeList, node)
					}
				}
			}

			if value, ok := record.Get("relationships"); ok {
				list := value.([]interface{})
				for _, ele := range list {
					relationship := ele.(neo4j.Relationship)
					if _, ok := nodeIdentitySet[relationship.ElementId]; !ok {
						nodeIdentitySet[relationship.ElementId] = struct{}{}
						relaList = append(relaList, relationship)
					}
				}
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Println("Read error:", err)
	}
	return nodeList, relaList, err
}
