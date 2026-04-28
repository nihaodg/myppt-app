# Oh My PPT

AI驱动的PPT生成器 - 用一句话描述需求，AI自动生成专业幻灯片。

![Preview](docs/preview.png)

## 功能特性

- **AI智能生成** - 输入主题或详细描述，AI自动规划大纲、配色和排版
- **36种风格技能** - 内置丰富的视觉风格，从极简白到赛博霓虹
- **本地优先** - 数据仅在本地处理，保护隐私安全
- **跨平台桌面应用** - Windows原生体验，双击即可运行
- **HTML幻灯片** - 生成的PPT为HTML格式，浏览器直接预览
- **OpenAI兼容** - 支持OpenAI API及本地Ollama模型

## 风格预览

| 风格 | 描述 |
|------|------|
| 极简白 | 克制高级，强文字层级 |
| 赛博霓虹 | 纯黑+霓虹粉青黄发光 |
| 融资路演 | YC风白底蓝紫渐变 |
| 日式极简 | 象牙白+朱红accent |
| 北欧 | 清冷蓝白，冷静理性 |
| 磨砂玻璃 | 毛玻璃+多色光斑 |
| ... | [更多风格](docs/styles.md) |

## 快速开始

### 下载安装

1. 从 [Releases](https://github.com/nihaodg/myppt-app/releases) 下载最新版本
2. 解压到任意目录
3. 双击 `oh-my-ppt.exe` 运行

### 配置API

首次使用需要配置AI API：

1. 点击右上角设置图标
2. 填写配置信息：
   - **Base URL**: `https://api.openai.com/v1` (或你的Ollama地址)
   - **Model**: `gpt-4o` (或本地模型如 `qwen2.5:14b`)
   - **API Key**: 你的密钥

### 生成PPT

1. 选择喜欢的风格主题
2. 在输入框描述你的PPT需求
3. 点击「生成PPT」按钮
4. 生成后在浏览器中打开预览

## 系统要求

- Windows 10/11 (64位)
- WebView2 Runtime (Windows 11自带，Windows 10[点击下载](https://go.microsoft.com/fwlink/p/?LinkId=2124703))

## 技术栈

- **后端**: Go + Wails
- **前端**: React + TypeScript + TailwindCSS
- **AI**: OpenAI API / Ollama

## 项目结构

```
oh-my-ppt/
├── app.go              # Go后端逻辑
├── main.go             # Wails入口
├── frontend/           # React前端
│   └── src/
│       ├── App.tsx     # 主界面组件
│       └── style.css   # 样式
├── resources/
│   └── styles.json     # 36种风格定义
└── build/              # 构建配置
```

## 开发

```bash
# 安装依赖
cd frontend && npm install

# 开发模式
wails dev

# 构建
wails build
```

## License

MIT License
