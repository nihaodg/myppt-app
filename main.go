package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/godjian/myppt-app/internal/config"
	"github.com/godjian/myppt-app/cmd"
)

func main() {
	// 设置工作目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	
	// 创建配置目录
	configDir := filepath.Join(homeDir, ".oh-my-ppt")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("创建配置目录失败: %v\n", err)
	}

	// 初始化配置
	if err := config.InitConfig(homeDir); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
	}

	// 执行命令
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		if !strings.Contains(err.Error(), "unknown command") &&
		   !strings.Contains(err.Error(), "required flag") {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	}
}
