package preview

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Server 预览服务器
type Server struct {
	port     int
	dir      string
	server   *http.Server
}

// NewServer 创建预览服务器
func NewServer(port int, dir string) *Server {
	return &Server{
		port: port,
		dir:  dir,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 检查目录
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		return fmt.Errorf("目录不存在: %s", s.dir)
	}

	// 检查 index.html
	indexPath := filepath.Join(s.dir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return fmt.Errorf("目录中没有 index.html: %s", s.dir)
	}

	// 创建 mux
	mux := http.NewServeMux()
	
	// 设置文件服务器
	fs := http.FileServer(http.Dir(s.dir))
	mux.Handle("/", fs)

	// 创建服务器
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	fmt.Printf("🌐 预览服务器启动中...\n")
	fmt.Printf("📂 路径: %s\n", s.dir)
	fmt.Printf("🌐 地址: http://localhost:%d\n", s.port)
	fmt.Printf("\n按 Ctrl+C 停止服务器\n\n")

	// 启动
	log.Fatal(s.server.ListenAndServe())
	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}
