package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"graph-med/internal/base/captcha"
	"graph-med/internal/base/conf"
	"graph-med/internal/base/data"
	"graph-med/internal/base/logger"
	"graph-med/internal/base/redis"
	"log"
	"os"
	"path/filepath"
)

func init() {
	for _, cmd := range []*cobra.Command{runCmd, initCmd} {
		rootCmd.AddCommand(cmd)
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "app",
		Short: "app应用",
		Long:  "app应用",
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "启动命令",
		Long:  "启动命令",
		Run: func(cmd *cobra.Command, args []string) {
			RunApp()
		},
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "初始化命令",
		Long:  "初始化命令",
		Run: func(cmd *cobra.Command, args []string) {
			InitDB()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// RunApp 启动应用
func RunApp() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前路径失败:", err)
		return
	}
	config, err := conf.LoadConfig(filepath.Join(currentDir, conf.DefaultConfigPath))
	if err != nil {
		fmt.Println("加载配置失败:", err)
		return
	}

	logger.Initialize(config)

	// 初始化Redis
	err = redis.Setup(config)
	if err != nil {
		fmt.Println("redis 连接失败:", err)
		return
	}

	// 初始化 email 验证码
	captcha.SetupEmailCaptcha(config.Email.Username, config.Email.Password, config.Email.Host, config.Email.Port)

	server, err := initializeServer(config)
	if err != nil {
		fmt.Println("初始化服务失败:", err)
		return
	}

	if err := server.Run(); err != nil {
		fmt.Println("启动服务失败:", err)
		return
	}
}

// InitDB 初始化数据库
func InitDB() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前路径失败:", err)
		return
	}
	config, err := conf.LoadConfig(filepath.Join(currentDir, conf.DefaultConfigPath))
	if err != nil {
		fmt.Println("加载配置失败:", err)
		return
	}

	err = data.InitDB(config)
	if err != nil {
		fmt.Println("初始化数据库失败:", err)
		return
	}
}
