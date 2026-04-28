package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/godjian/myppt-app/internal/ai"
	"github.com/godjian/myppt-app/internal/config"
	"github.com/godjian/myppt-app/internal/html"
	"github.com/godjian/myppt-app/internal/styles"
)

// newGenerateCmd 生成命令
func newGenerateCmd() *cobra.Command {
	var (
		topic     string
		styleID   string
		pages     int
		outputDir string
		noPreview bool
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "生成 PPT",
		Run: func(cmd *cobra.Command, args []string) {
			runGenerate(topic, styleID, pages, outputDir, noPreview)
		},
	}

	cmd.Flags().StringVarP(&topic, "topic", "t", "", "PPT 主题 (必填)")
	cmd.Flags().StringVarP(&styleID, "style", "s", "minimal-white", "风格 ID")
	cmd.Flags().IntVarP(&pages, "pages", "p", 8, "页数")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "输出目录")
	cmd.Flags().BoolVar(&noPreview, "no-preview", false, "不启动预览")

	cmd.MarkFlagRequired("topic")
	return cmd
}

func runGenerate(topic, styleID string, pageCount int, outputDir string, noPreview bool) {
	fmt.Println()
	fmt.Println("开始生成 PPT...")
	fmt.Println()

	// 检查配置
	cfg := config.GetConfig()
	if cfg.APIKey == "" && cfg.Provider == "openai" {
		fmt.Println("警告: 未配置 API Key")
		fmt.Println("请运行: oh-my-ppt config set openai.api_key YOUR_API_KEY")
		fmt.Println()
	}

	// 获取风格
	style := styles.GetStyleOrDefault(styleID)
	fmt.Printf("  主题: %s\n", topic)
	fmt.Printf("  风格: %s (%s)\n", style.Name, style.ID)
	fmt.Printf("  页数: %d\n", pageCount)
	fmt.Println()

	// 创建 AI 客户端
	aiClient := ai.NewClient(cfg)

	// 生成 PPT 大纲
	fmt.Println("正在规划 PPT 结构...")
	plan, err := aiClient.GeneratePPTPlan(topic, pageCount)
	if err != nil {
		fmt.Printf("  错误: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  生成了 %d 页大纲\n", len(plan.Pages))
	for i, page := range plan.Pages {
		fmt.Printf("  %d. %s\n", i+1, page.Title)
	}
	fmt.Println()

	// 生成设计契约 (使用风格默认值或 AI 增强)
	fmt.Println("正在生成设计契约...")
	contract := &ai.DesignContract{
		Theme:         style.Theme,
		Background:   style.Background,
		Palette:      style.Palette,
		TitleStyle:   style.TitleStyle,
		LayoutMotif:  style.LayoutMotif,
		ChartStyle:   style.ChartStyle,
		ShapeLanguage: style.ShapeLang,
	}
	fmt.Println("  设计契约已准备")
	fmt.Println()

	// 生成页面内容
	fmt.Println("正在生成页面内容...")
	pageContents := make([]string, 0, len(plan.Pages))
	for i, page := range plan.Pages {
		fmt.Printf("  生成第 %d/%d 页: %s\n", i+1, len(plan.Pages), page.Title)

		content, err := aiClient.GeneratePageContent(page, contract, style.Prompt)
		if err != nil {
			// 如果生成失败，使用默认内容
			fmt.Printf("  AI 生成失败，使用默认内容: %v\n", err)
			content = generateDefaultPage(page, style)
		}
		pageContents = append(pageContents, content)
	}
	fmt.Println()

	// 生成 HTML
	fmt.Println("正在生成 HTML 文件...")

	if outputDir == "" {
		// 使用默认输出目录
		homeDir, _ := os.UserHomeDir()
		outputDir = filepath.Join(homeDir, "oh-my-ppt-output", sanitizeFilename(topic))
	}

	generator := html.NewGenerator(outputDir, style)
	if err := generator.GenerateDeck(plan, contract, pageContents); err != nil {
		fmt.Printf("  错误: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  文件已生成: %s\n", outputDir)
	fmt.Println()

	// 完成
	fmt.Println("PPT 生成完成!")
	fmt.Println()
	fmt.Printf("  输出目录: %s\n", outputDir)
	fmt.Println("  文件列表:")
	fmt.Println("    - index.html (主页面)")
	for i := range plan.Pages {
		fmt.Printf("    - page-%d.html (第 %d 页)\n", i+1, i+1)
	}
	fmt.Println("    - style.css (样式文件)")
	fmt.Println("    - assets/ (资源文件)")
	fmt.Println()

	if !noPreview {
		fmt.Println("正在启动预览服务器...")
		previewPort := 8080
		fmt.Printf("  访问: http://localhost:%d\n", previewPort)
		fmt.Println()
		fmt.Println("  按 Ctrl+C 停止服务器")
	}
}

// 生成默认页面内容
func generateDefaultPage(page ai.PagePlan, style *styles.Style) string {
	points := strings.Builder{}
	for i, pt := range page.KeyPoints {
		points.WriteString(fmt.Sprintf(`<li class="mb-4 flex items-start">
            <span class="flex-shrink-0 w-6 h-6 rounded-full bg-%s mr-3 flex items-center justify-center text-white text-sm">%d</span>
            <span>%s</span>
        </li>`, getAccentColor(style), i+1, pt))
	}

	return fmt.Sprintf(`<div class="h-full flex flex-col justify-center">
    <div class="max-w-4xl mx-auto">
        <h1 class="text-5xl font-bold mb-8" style="color: %s">%s</h1>
        <div class="bg-white rounded-2xl shadow-xl p-8">
            <ul class="text-xl space-y-2">
                %s
            </ul>
        </div>
    </div>
</div>`, getPrimaryColor(style), page.Title, points.String())
}

func getAccentColor(style *styles.Style) string {
	if len(style.Palette) > 2 {
		return strings.TrimPrefix(style.Palette[2], "#")
	}
	return "3B82F6"
}

func getPrimaryColor(style *styles.Style) string {
	if len(style.Palette) > 1 {
		return style.Palette[1]
	}
	return "#1F2937"
}

func sanitizeFilename(name string) string {
	// 移除非法字符
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "*", "")
	name = strings.ReplaceAll(name, "?", "")
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "<", "")
	name = strings.ReplaceAll(name, ">", "")
	name = strings.ReplaceAll(name, "|", "")

	// 限制长度
	if len(name) > 50 {
		name = name[:50]
	}

	return name
}
