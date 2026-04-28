# Godjian PPT - 本地幻灯片生成工具

## 快速开始

### 方法一：直接打开（推荐）

1. 双击 `start.bat` 启动应用
2. 或双击 `index.html` 在浏览器中打开（部分功能受限）

### 方法二：Python 环境

1. 确保已安装 Python 3.8+
2. 运行 `start.bat`

## 打包为 exe（Windows）

在已安装 Python 的电脑上运行：

```batch
pip install pyinstaller
pyinstaller --onefile --windowed --name godjian-ppt app.py
```

打包后的 exe 文件在 `dist/godjian-ppt.exe`

## 功能特点

- 20+ 内置风格（极简白、赛博霓虹、包豪斯、日式简约等）
- 本地优先，数据不上传
- 生成 HTML 格式幻灯片，浏览器直接预览
- 支持键盘导航和全屏演示
