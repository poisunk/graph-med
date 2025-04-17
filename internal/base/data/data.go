package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"graph-med/internal/base/conf"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type Data struct {
	DB          *xorm.Engine
	Neo4jDriver neo4j.DriverWithContext
}

func NewData(engine *xorm.Engine, neo4jDriver neo4j.DriverWithContext) *Data {
	return &Data{
		DB:          engine,
		Neo4jDriver: neo4jDriver,
	}
}

func NewDB(config *conf.AllConfig) (*xorm.Engine, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}

	debug := config.Debug
	dataConf := &config.Data.Database

	// 默认使用mysql
	if dataConf.Driver == "" {
		dataConf.Driver = "mysql"
	}

	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dataConf.Username,
		dataConf.Password,
		dataConf.Addr,
		dataConf.Port,
		dataConf.Database,
	)

	// 创建引擎
	engine, err := xorm.NewEngine(dataConf.Driver, connection)
	if err != nil {
		return nil, err
	}

	if debug {
		engine.ShowSQL(true)
	} else {
		engine.SetLogLevel(log.LOG_ERR)
	}

	// 设置连接池
	engine.SetMaxOpenConns(dataConf.MaxOpenConns)
	engine.SetMaxIdleConns(dataConf.MaxIdleConns)
	engine.SetConnMaxLifetime(time.Duration(dataConf.MaxLifeSecs) * time.Second)

	return engine, nil
}

func InitDB(config *conf.AllConfig) error {
	engine, err := NewDB(config)
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = engine.ImportFile(filepath.Join(dir, "/db/graph-med.sql"))
	if err != nil {
		return err
	}

	return nil
}

func NewNeo4jDriver(config *conf.AllConfig) (neo4j.DriverWithContext, error) {
	neo4jConfig := &config.Data.Neo4j
	driver, err := neo4j.NewDriverWithContext(neo4jConfig.Addr, neo4j.BasicAuth(neo4jConfig.Username, neo4jConfig.Password, ""))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}
	return driver, nil
}
