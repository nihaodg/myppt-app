@echo off
echo ========================================
echo Godjian PPT - 本地幻灯片生成工具
echo ========================================
echo.

:: 检查 Python
python --version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未安装 Python
    echo 请从 https://www.python.org/downloads/ 下载安装
    pause
    exit /b 1
)

:: 运行应用
echo 正在启动应用...
python app.py

pause
