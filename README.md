# Godjian PPT - 本地幻灯片生成工具

一个简洁的本地 PPT 生成工具，基于 Go 语言开发。

## 功能特点

- 20+ 内置风格（极简白、赛博霓虹、包豪斯、日式简约等）
- 本地优先，数据不上传
- 生成 HTML 格式幻灯片，浏览器直接预览
- 支持键盘导航和全屏演示

## 快速开始

### 下载 exe

前往 [Releases](https://github.com/nihaodg/myppt-app/releases) 下载最新版本

### 或使用源代码

1. 克隆项目
2. 使用 Wails 构建：
```bash
wails build
```

## 文件说明

- `index.html` - 直接在浏览器中打开
- `app.py` - Python 版本（需要 Python 环境）
- `start.bat` - Windows 启动脚本

## 打包为 exe

### Python 版本

```bash
pip install pyinstaller
pyinstaller --onefile --windowed --name godjian-ppt app.py
```

## License

MIT
