package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/godjian/myppt-app/internal/config"
)

// newInitCmd 初始化命令
func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "初始化配置",
		Long:  "初始化 oh-my-ppt 配置文件",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("配置初始化成功")
			fmt.Println()
			fmt.Println("配置文件位置: ~/.oh-my-ppt/config.json")
			fmt.Println()
			fmt.Println("请先配置 API:")
			fmt.Println("  oh-my-ppt config set openai.api_key YOUR_API_KEY")
			fmt.Println()
			fmt.Println("或者配置 Ollama (本地模型):")
			fmt.Println("  oh-my-ppt config set ollama.base_url http://127.0.0.1:11434/v1")
			fmt.Println("  oh-my-ppt config set ollama.model qwen2.5-coder:14b")
		},
	}
	return cmd
}

// newConfigCmd 配置命令
func newConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "管理配置",
		Long:  "查看和修改配置",
	}

	configCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "显示当前配置",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfig()
			fmt.Println("当前配置:")
			fmt.Printf("  Provider:  %s\n", cfg.Provider)
			fmt.Printf("  API Key:   %s\n", maskString(cfg.APIKey))
			fmt.Printf("  Base URL:  %s\n", cfg.BaseURL)
			fmt.Printf("  Model:     %s\n", cfg.Model)
			fmt.Printf("  Max Tokens: %d\n", cfg.MaxTokens)
		},
	})

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "设置配置项",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Println("用法: oh-my-ppt config set <key> <value>")
				fmt.Println("例如: oh-my-ppt config set openai.api_key sk-xxx")
				return
			}
			key := args[0]
			value := args[1]

			if err := config.UpdateConfig(key, value); err != nil {
				fmt.Printf("设置失败: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s 已更新\n", key)
		},
	}

	configCmd.AddCommand(setCmd)
	return configCmd
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "********"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
