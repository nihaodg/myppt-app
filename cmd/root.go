package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewRootCmd 创建主命令
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oh-my-ppt",
		Short: "AI-Powered PPT Generator",
		Long: `本地优先的 AI 幻灯片生成工具
输入主题 -> AI 规划大纲 -> 生成风格 -> 逐页渲染 -> 预览 & 导出`,
		Run: func(cmd *cobra.Command, args []string) {
			printBanner()
			printHelp()
		},
	}

	// 添加子命令
	cmd.AddCommand(newGenerateCmd())
	cmd.AddCommand(newInteractiveCmd())
	cmd.AddCommand(newStylesCmd())
	cmd.AddCommand(newPreviewCmd())
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newInitCmd())

	// 全局标志
	cmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出")
	cmd.PersistentFlags().Bool("no-color", false, "禁用颜色输出")

	return cmd
}

// 打印横幅
func printBanner() {
	banner := `
===============================================
                                            
   ██████╗  ██████╗ ██████╗ ████████╗ ██████╗ ██████╗   
   ██╔══██╗██╔════╝██╔═══██╗╚══██╔══╝██╔═══██╗██╔══██╗  
   ██████╔╝██║     ██║   ██║   ██║   ██║   ██║██████╔╝  
   ██╔══██╗██║     ██║   ██║   ██║   ██║   ██║██╔══██╗  
   ██║  ██║╚██████╗╚██████╔╝   ██║   ╚██████╔╝██║  ██║  
   ╚═╝  ╚═╝ ╚═════╝ ╚═════╝    ╚═╝    ╚═════╝ ╚═╝  ╚═╝  
                                            
   AI-Powered PPT Generator in Go              
   Local-first - 30+ Styles - Fast Generation 
                                            
===============================================
`
	fmt.Println(banner)
}

// 打印帮助信息
func printHelp() {
	fmt.Println(`
使用指南:

  快速开始:
    oh-my-ppt generate --topic "AI发展趋势" --style minimal-white

  交互模式:
    oh-my-ppt interactive

  查看风格列表:
    oh-my-ppt styles

  预览 PPT:
    oh-my-ppt preview ./output/my-ppt

  导出为 PDF:
    oh-my-ppt export --format pdf ./output/my-ppt

  配置 API:
    oh-my-ppt config set openai.api_key YOUR_API_KEY
    oh-my-ppt config set openai.base_url https://api.openai.com/v1
    oh-my-ppt config set openai.model gpt-4o

  配置 Ollama (本地模型):
    oh-my-ppt config set ollama.base_url http://127.0.0.1:11434/v1
    oh-my-ppt config set ollama.model qwen2.5-coder:14b

支持风格:
    minimal-white     - 极简白
    cyber-neon       - 赛博霓虹
    bauhaus          - 包豪斯
    japanese-minimal  - 日式简约
    corporate-blue    - 企业蓝
    nature-green     - 自然绿
    dark-tech        - 暗黑科技
    retro-warm       - 复古暖色
    ... (更多风格请运行 oh-my-ppt styles)

文档:
    https://github.com/godjian/myppt-app
`)
}
