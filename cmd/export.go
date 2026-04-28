package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// newExportCmd 导出命令
func newExportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export [path]",
		Short: "导出 PPT 为 PDF",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]

			// 检查路径
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("路径不存在: %s\n", path)
				os.Exit(1)
			}

			// 确保是目录
			if !isDir(path) {
				path = filepath.Dir(path)
			}

			fmt.Println("PDF 导出需要浏览器手动操作")
			fmt.Println()
			fmt.Println("请使用浏览器手动导出:")
			fmt.Println("  1. 在浏览器中打开:", filepath.Join(path, "index.html"))
			fmt.Println("  2. Ctrl+P 打印")
			fmt.Println("  3. 选择 '另存为 PDF'")
			fmt.Println("  4. 布局选择 '横向'")
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "pdf", "导出格式 (pdf/png/pptx)")
	return cmd
}
