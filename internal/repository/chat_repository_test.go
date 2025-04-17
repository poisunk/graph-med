package repository

import (
	"fmt"
	"graph-med/internal/base/conf"
	"graph-med/internal/base/data"
	"graph-med/internal/model"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"xorm.io/xorm"
)

func TestChatRepository_CreateChatMessage(t *testing.T) {
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
	repo := NewChatRepository(data)

	msg := &model.ChatMessage{
		MessageID:       1,
		ParentMessageID: 0,
		UserID:          "123",
		ChatSessionID:   "123",
		Role:            "user",
		Content:         "hello",
	}

	err = repo.CreateChatMessage(msg)
	if err != nil {
		t.Errorf("CreateChatMessage error: %v", err)
		return
	}

	msgs, err := repo.GetChatMessageBySessionID("123", 10)
	if err != nil {
		t.Errorf("GetChatMessageBySessionID error: %v", err)
		return
	}

	for _, msg := range msgs {
		fmt.Printf("MessageID: %d, ParentMessageID: %d, UserID: %s, ChatSessionID: %s, Role: %s, Content: %s\n",
			msg.MessageID, msg.ParentMessageID, msg.UserID, msg.ChatSessionID, msg.Role, msg.Content)

		err := repo.DeleteChatMessage(msg.ID)
		if err != nil {
			t.Errorf("DeleteChatMessage error: %v", err)
			return
		}
	}
}
