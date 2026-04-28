# Godjian PPT - Local PPT Generator

本地优先的幻灯片生成工具，基于 Go 语言开发。

## 功能特点

- 💬 **一句话生成** - 输入主题，自动规划大纲 + 配色 + 排版，直接出完整 PPT
- 🔒 **本地优先** - 全部跑在自己电脑上，不用注册、不用担心数据泄露
- 🎨 **30+ 内置风格** - 极简白、赛博霓虹、包豪斯、日式简约等
- 📄 **HTML 输出** - 生成 HTML 版 PPT，打开即预览、无需软件
- 🎬 **动画支持** - 基于内置 JavaScript 运行时的动画效果
- 🌐 **预览服务器** - 内置 HTTP 服务器预览 HTML PPT

## 快速开始

### 下载和安装

```bash
# 下载对应平台的 release 文件
# Windows: godjian-ppt.exe
# Linux/macOS: godjian-ppt
```

### 配置 API

#### OpenAI API
```bash
./godjian-ppt config set openai.api_key YOUR_API_KEY
./godjian-ppt config set openai.base_url https://api.openai.com/v1
./godjian-ppt config set openai.model gpt-4o
```

#### Ollama (本地模型)
```bash
./godjian-ppt config set ollama.base_url http://127.0.0.1:11434/v1
./godjian-ppt config set ollama.model qwen2.5-coder:14b
```

### 生成 PPT

```bash
# 基本用法
./godjian-ppt generate --topic "技术分享" --style minimal-white

# 指定页数
./godjian-ppt generate -t "产品介绍" -s cyber-neon -p 12

# 指定输出目录
./godjian-ppt generate -t "年度总结" -o ~/my-ppt
```

### 交互模式

```bash
./godjian-ppt interactive
```

### 查看风格列表

```bash
./godjian-ppt styles
```

### 预览 PPT

```bash
./godjian-ppt preview ./output/my-ppt
```

## 支持的风格

| 风格 ID | 名称 | 描述 |
|---------|------|------|
| minimal-white | 极简白 | 简洁干净的白色风格 |
| cyber-neon | 赛博霓虹 | 未来科技感的霓虹风格 |
| bauhaus | 包豪斯 | 德国包豪斯风格 |
| japanese-minimal | 日式简约 | 日式极简美学 |
| corporate-blue | 企业蓝 | 专业的企业蓝色调 |
| nature-green | 自然绿 | 清新自然的绿色主题 |
| dark-tech | 暗黑科技 | 深色科技风格 |
| retro-warm | 复古暖色 | 温暖的复古色调 |
| elegant-purple | 优雅紫 | 高贵的紫色主题 |
| ocean-blue | 海洋蓝 | 清新的海洋主题 |
| ... | ... | 共 30+ 种风格 |

## 输出文件

生成后会创建以下文件：

```
output-dir/
├── index.html      # 主页面 (PPT 导航)
├── page-1.html    # 第 1 页
├── page-2.html    # 第 2 页
├── ...
├── page-N.html    # 第 N 页
├── style.css      # 样式文件
└── assets/        # 资源文件
    └── ppt-runtime.js  # PPT 运行时
```

## 使用方法

1. **在浏览器中打开** `index.html`
2. **使用左右箭头键** 或点击导航切换页面
3. **按 F11** 进入全屏演示模式
4. **按 ESC** 退出全屏

## 导出为 PDF

1. 在浏览器中打开 `index.html`
2. 按 `Ctrl+P` 打开打印对话框
3. 选择"另存为 PDF"
4. 布局选择"横向"

## 从源码构建

### 环境要求

- Go 1.21+
- GCC (用于 cgo)

### 构建步骤

```bash
# 克隆项目
git clone https://github.com/godjian/myppt-app.git
cd myppt-app

# 下载依赖
go mod tidy

# 构建 Windows exe
GOOS=windows GOARCH=amd64 go build -o godjian-ppt.exe .

# 构建 Linux
GOOS=linux GOARCH=amd64 go build -o godjian-ppt .

# 构建 macOS
GOOS=darwin GOARCH=amd64 go build -o godjian-ppt .
```

## 项目结构

```
myppt-app/
├── main.go                 # 程序入口
├── cmd/                    # CLI 命令
│   ├── root.go            # 根命令
│   ├── generate.go        # 生成命令
│   ├── styles.go         # 风格列表命令
│   ├── preview.go        # 预览命令
│   ├── export.go        # 导出命令
│   ├── config.go        # 配置命令
│   └── interactive.go   # 交互模式
├── internal/
│   ├── ai/              # API 客户端
│   │   └── client.go    # OpenAI/Ollama API 调用
│   ├── config/          # 配置管理
│   │   └── config.go
│   ├── styles/          # 风格系统
│   │   └── builtin.go   # 30+ 内置风格
│   ├── html/            # HTML 生成器
│   │   └── generator.go
│   └── preview/         # 预览服务器
│       └── server.go
└── go.mod
```

## 技术栈

- **语言**: Go 1.21+
- **CLI 框架**: Cobra
- **集成**: OpenAI API / Ollama (OpenAI 兼容)
- **输出格式**: HTML + Tailwind CSS + JavaScript

## License

MIT License
