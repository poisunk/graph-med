package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"graph-med/internal/base/conf"
	"log"
	"os"
	"path/filepath"
)

func init() {
	for _, cmd := range []*cobra.Command{runCmd} {
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
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// RunApp 启动mcp服务
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
	mcpServer, err := initializeMcpServer(config)
	if err != nil {
		fmt.Println("初始化mcp服务失败:", err)
		return
	}

	if err := mcpServer.Run(); err != nil {
		fmt.Println("启动mcp服务失败:", err)
		return
	}
}
