package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// newPreviewCmd 预览命令
func newPreviewCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "preview [path]",
		Short: "启动预览服务器",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]

			// 检查路径是否存在
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Printf("路径不存在: %s\n", path)
				os.Exit(1)
			}

			// 如果是文件，提取目录
			if !isDir(path) {
				path = filepath.Dir(path)
			}

			// 确保有 index.html
			indexPath := filepath.Join(path, "index.html")
			if _, err := os.Stat(indexPath); os.IsNotExist(err) {
				fmt.Printf("错误: %s 目录中找不到 index.html\n", path)
				os.Exit(1)
			}

			// 启动服务器
			fmt.Printf("预览服务器启动中...\n")
			fmt.Printf("路径: %s\n", path)
			fmt.Printf("地址: http://localhost:%d\n", port)
			fmt.Printf("\n按 Ctrl+C 停止服务器\n\n")

			fs := http.FileServer(http.Dir(path))
			http.Handle("/", fs)

			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "预览服务器端口")
	return cmd
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
