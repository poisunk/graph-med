package conf

import (
	"os"

	"gopkg.in/yaml.v2"
)

const (
	DefaultConfigPath = "configs/config.yaml"
)

type AllConfig struct {
	Debug  bool   `yaml:"debug"`
	Server Server `yaml:"server"`
	Data   Data   `yaml:"data"`
	Redis  Redis  `yaml:"redis"`
	Email  Email  `yaml:"email"`
	Mcp    Mcp    `yaml:"mcp"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Data struct {
	Database Database `yaml:"database"`
	Neo4j    Neo4j    `yaml:"neo4j"`
}

type Database struct {
	Driver       string `yaml:"driver"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Addr         string `yaml:"addr"`
	Port         int    `yaml:"port"`
	Database     string `yaml:"database"`
	MaxOpenConns int    `yaml:"max_open_connections"`
	MaxIdleConns int    `yaml:"max_idle_connections"`
	MaxLifeSecs  int    `yaml:"max_life_seconds"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
	PoolSize int    `yaml:"pool_size"`
}

type Neo4j struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Email struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Mcp struct {
	Port string `yaml:"port"`
}

func LoadConfig(path string) (*AllConfig, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config AllConfig
	err = yaml.Unmarshal(yamlFile, &config)
	return &config, err
}
