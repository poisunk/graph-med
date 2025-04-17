package repository

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"graph-med/internal/base/conf"
	"graph-med/internal/base/data"
	"os"
	"path/filepath"
	"testing"
	"xorm.io/xorm"
)

func TestKGRepository_GetNodeByName(t *testing.T) {
	wdPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Getwd error: %v", err)
		return
	}

	yamlFile, err := os.ReadFile(filepath.Join(wdPath, "../../configs/config.yaml"))
	if err != nil {
		t.Errorf("ReadFile error: %v", err)
		return
	}

	var config conf.AllConfig
	err = yaml.Unmarshal(yamlFile, &config)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Data.Database.Username,
		config.Data.Database.Password,
		config.Data.Database.Addr,
		config.Data.Database.Port,
		config.Data.Database.Database,
	)

	// 创建引擎
	engine, err := xorm.NewEngine(config.Data.Database.Driver, connection)
	if err != nil {
		t.Errorf("NewEngine error: %v", err)
		return
	}
	engine.ShowSQL(true)
	driver, err := data.NewNeo4jDriver(&config)
	if err != nil {
		t.Errorf("NewNeo4jDriver error: %v", err)
		return
	}

	data := data.NewData(engine, driver)

	repo := NewKGRepository(data)
	nodes, relations, err := repo.GetNodeSubGraph("疾病", "肺栓塞", 100)
	if err != nil {
		t.Errorf("GetNodeByName error: %v", err)
		return
	}
	fmt.Printf("node count: %d, relation count: %d\n", len(nodes), len(relations))
	fmt.Printf("nodes: %v\nrelations: %v\n", nodes, relations)
}
