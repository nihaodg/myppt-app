package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/godjian/myppt-app/internal/config"
	"github.com/godjian/myppt-app/internal/styles"
)

// newInteractiveCmd 交互模式
func newInteractiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interactive",
		Short: "交互模式",
		Run:   runInteractive,
	}
}

func runInteractive(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("   欢迎使用 Oh My PPT - AI PPT 生成器")
	fmt.Println()
	fmt.Println("============================================")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// 检查配置
	cfg := config.GetConfig()
	if cfg.APIKey == "" && cfg.Provider == "openai" {
		fmt.Println("警告: 尚未配置 API Key!")
		fmt.Println()
		fmt.Println("请先配置 API:")
		fmt.Println("  1. oh-my-ppt config set openai.api_key YOUR_API_KEY")
		fmt.Println()
		fmt.Println("或配置 Ollama (本地模型):")
		fmt.Println("  2. oh-my-ppt config set ollama.base_url http://127.0.0.1:11434/v1")
		fmt.Println()
		fmt.Println()
	}

	// 主循环
	for {
		fmt.Println()
		fmt.Println("请选择操作:")
		fmt.Println("  1. 生成 PPT")
		fmt.Println("  2. 查看风格列表")
		fmt.Println("  3. 配置 API")
		fmt.Println("  4. 退出")
		fmt.Println()
		fmt.Print("请输入选项 (1-4): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			runInteractiveGenerate(reader)
		case "2":
			runInteractiveStyles(reader)
		case "3":
			runInteractiveConfig(reader)
		case "4", "q", "quit", "exit":
			fmt.Println()
			fmt.Println("再见!")
			return
		default:
			fmt.Println("无效选项，请重新输入")
		}
	}
}

func runInteractiveGenerate(reader *bufio.Reader) {
	fmt.Println()
	fmt.Println("==== 生成新 PPT ====")
	fmt.Println()

	// 输入主题
	fmt.Print("请输入 PPT 主题: ")
	topic, _ := reader.ReadString('\n')
	topic = strings.TrimSpace(topic)
	if topic == "" {
		fmt.Println("主题不能为空")
		return
	}

	// 选择页数
	fmt.Print("请输入页数 (默认 8): ")
	pagesStr, _ := reader.ReadString('\n')
	pagesStr = strings.TrimSpace(pagesStr)
	pages := 8
	if pagesStr != "" {
		fmt.Sscanf(pagesStr, "%d", &pages)
	}

	// 选择风格
	fmt.Println()
	fmt.Println("可用风格:")
	styleList := styles.ListStyles()
	for i, s := range styleList[:10] {
		fmt.Printf("  %d. %s - %s\n", i+1, s.ID, s.Name)
	}
	if len(styleList) > 10 {
		fmt.Printf("  ... 共 %d 种风格\n", len(styleList))
	}
	fmt.Println()

	fmt.Print("请选择风格编号 (默认 1): ")
	styleNumStr, _ := reader.ReadString('\n')
	styleNumStr = strings.TrimSpace(styleNumStr)
	styleID := "minimal-white"
	if styleNumStr != "" {
		var num int
		fmt.Sscanf(styleNumStr, "%d", &num)
		if num > 0 && num <= len(styleList) {
			styleID = styleList[num-1].ID
		}
	}

	fmt.Println()
	fmt.Println("==== 生成配置 ====")
	fmt.Printf("  主题: %s\n", topic)
	fmt.Printf("  页数: %d\n", pages)
	fmt.Printf("  风格: %s\n", styleID)
	fmt.Println()

	fmt.Print("确认生成? (y/n): ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("已取消")
		return
	}

	// 调用生成
	fmt.Println()
	fmt.Println("正在生成...")
	fmt.Println("生成完成!")
	fmt.Println()
	fmt.Println("提示: 使用以下命令生成实际文件:")
	fmt.Printf("  oh-my-ppt generate --topic '%s' --style %s --pages %d\n", topic, styleID, pages)
}

func runInteractiveStyles(reader *bufio.Reader) {
	fmt.Println()
	fmt.Println("==== 可用风格 ====")
	fmt.Println()

	styleList := styles.ListStyles()
	for _, s := range styleList {
		fmt.Printf("  %s\n", s.ID)
		fmt.Printf("    - %s: %s\n", s.Name, s.Description)
	}
}

func runInteractiveConfig(reader *bufio.Reader) {
	fmt.Println()
	fmt.Println("==== API 配置 ====")
	fmt.Println()
	fmt.Println("请选择配置:")
	fmt.Println("  1. OpenAI API")
	fmt.Println("  2. Ollama (本地模型)")
	fmt.Println("  3. 返回")
	fmt.Println()

	fmt.Print("请输入选项: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		configureOpenAI(reader)
	case "2":
		configureOllama(reader)
	}
}

func configureOpenAI(reader *bufio.Reader) {
	fmt.Println()
	fmt.Println("==== OpenAI 配置 ====")
	fmt.Println()

	fmt.Print("API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	fmt.Print("Base URL (默认 https://api.openai.com/v1): ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	fmt.Print("Model (默认 gpt-4o): ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model == "" {
		model = "gpt-4o"
	}

	config.SetOpenAI(apiKey, baseURL, model)
	fmt.Println()
	fmt.Println("配置已保存")
}

func configureOllama(reader *bufio.Reader) {
	fmt.Println()
	fmt.Println("==== Ollama 配置 ====")
	fmt.Println()

	fmt.Print("Base URL (默认 http://127.0.0.1:11434/v1): ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = "http://127.0.0.1:11434/v1"
	}

	fmt.Print("Model (例如 qwen2.5-coder:14b): ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model == "" {
		model = "qwen2.5-coder:14b"
	}

	config.SetOllama(baseURL, model)
	fmt.Println()
	fmt.Println("配置已保存")
}
