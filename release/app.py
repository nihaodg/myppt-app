import os
import sys
import json
import webbrowser
import http.server
import socketserver
import threading
import tempfile
import shutil
from pathlib import Path

# 风格数据
STYLES = [
    {"id": "minimal-white", "name": "极简白", "description": "简洁干净的白色风格", "palette": ["#FFFFFF", "#F9FAFB", "#111827", "#3B82F6", "#10B981"]},
    {"id": "cyber-neon", "name": "赛博霓虹", "description": "未来科技感的霓虹风格", "palette": ["#0F172A", "#1E293B", "#06B6D4", "#8B5CF6", "#EC4899"]},
    {"id": "bauhaus", "name": "包豪斯", "description": "德国包豪斯风格", "palette": ["#FFFFFF", "#F5F5DC", "#E53935", "#1E88E5", "#FDD835"]},
    {"id": "japanese-minimal", "name": "日式简约", "description": "日式极简美学", "palette": ["#FBF7F4", "#2C2C2C", "#C4A484", "#8B7355", "#D4C4B0"]},
    {"id": "corporate-blue", "name": "企业蓝", "description": "专业的企业蓝色调", "palette": ["#FFFFFF", "#EFF6FF", "#1E40AF", "#3B82F6", "#60A5FA"]},
    {"id": "nature-green", "name": "自然绿", "description": "清新自然的绿色主题", "palette": ["#FFFFFF", "#F0FDF4", "#166534", "#22C55E", "#86EFAC"]},
    {"id": "dark-tech", "name": "暗黑科技", "description": "深色科技风格", "palette": ["#0A0A0A", "#171717", "#404040", "#22D3EE", "#A855F7"]},
    {"id": "retro-warm", "name": "复古暖色", "description": "温暖的复古色调", "palette": ["#FFF8F0", "#2D2A26", "#D97706", "#DC2626", "#059669"]},
    {"id": "elegant-purple", "name": "优雅紫", "description": "高贵的紫色主题", "palette": ["#FFFFFF", "#FAF5FF", "#7C3AED", "#A855F7", "#C084FC"]},
    {"id": "ocean-blue", "name": "海洋蓝", "description": "清新的海洋主题", "palette": ["#FFFFFF", "#F0F9FF", "#0369A1", "#0EA5E9", "#38BDF8"]},
    {"id": "sunset-orange", "name": "日落橙", "description": "温暖的日落色调", "palette": ["#FFFFFF", "#FFFBEB", "#C2410C", "#F97316", "#FDBA74"]},
    {"id": "mint-fresh", "name": "薄荷清新", "description": "清新的薄荷绿", "palette": ["#FFFFFF", "#F0FDFA", "#0D9488", "#14B8A6", "#5EEAD4"]},
    {"id": "rose-pink", "name": "玫瑰粉", "description": "温柔的玫瑰粉色", "palette": ["#FFFFFF", "#FFF1F2", "#BE185D", "#EC4899", "#F9A8D4"]},
    {"id": "nordic-light", "name": "北欧光", "description": "明亮的北欧风格", "palette": ["#FFFFFF", "#FAFAFA", "#1F2937", "#4B5563", "#9CA3AF"]},
    {"id": "gold-luxury", "name": "金色奢华", "description": "高贵的金色主题", "palette": ["#000000", "#1C1917", "#D97706", "#F59E0B", "#FCD34D"]},
    {"id": "cherry-blossom", "name": "樱花浪漫", "description": "温柔的樱花色调", "palette": ["#FFF0F5", "#FDF2F8", "#BE185D", "#EC4899", "#F9A8D4"]},
    {"id": "monochrome", "name": "黑白单色", "description": "极简的黑白色调", "palette": ["#FFFFFF", "#F9FAFB", "#111827", "#374151", "#6B7280"]},
    {"id": "academic-blue", "name": "学术蓝", "description": "严谨的学术风格", "palette": ["#FFFFFF", "#EFF6FF", "#1E3A8A", "#2563EB", "#60A5FA"]},
    {"id": "startup-modern", "name": "创业现代", "description": "适合创业公司的现代风格", "palette": ["#FFFFFF", "#F8FAFC", "#6366F1", "#8B5CF6", "#A855F7"]},
    {"id": "forest-earth", "name": "森林大地", "description": "沉稳的大地色系", "palette": ["#FFFDF7", "#FEF3C7", "#365314", "#65A30D", "#84CC16"]},
]

# 尝试导入 tkinter (内置)
try:
    import tkinter as tk
    from tkinter import ttk
    TKINTER_AVAILABLE = True
except ImportError:
    TKINTER_AVAILABLE = False


def generate_html_index(title, page_count, style):
    thumbs = "".join(f'<div class="slide-thumb{" active" if i == 1 else ""}" onclick="goToSlide({i})">{i}</div>' for i in range(1, page_count + 1))
    return f'''<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{title} - PPT</title>
    <style>
        * {{ margin: 0; padding: 0; box-sizing: border-box; }}
        body {{ font-family: system-ui, sans-serif; background: #1a1a2e; color: #fff; overflow: hidden; }}
        .container {{ display: flex; height: 100vh; }}
        .sidebar {{ width: 200px; background: #16213e; padding: 20px; overflow-y: auto; }}
        .sidebar h3 {{ margin-bottom: 16px; color: #9ca3af; font-size: 14px; }}
        .slide-thumb {{ width: 100%; aspect-ratio: 16/9; background: #0f3460; border-radius: 4px; margin-bottom: 8px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 18px; font-weight: bold; transition: all 0.2s; }}
        .slide-thumb:hover {{ background: #1e5f74; transform: scale(1.05); }}
        .slide-thumb.active {{ background: #3b82f6; }}
        .main {{ flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 20px; }}
        .frame {{ width: 100%; max-width: 1200px; aspect-ratio: 16/9; background: #000; border-radius: 8px; overflow: hidden; box-shadow: 0 20px 60px rgba(0,0,0,0.5); }}
        .frame iframe {{ width: 100%; height: 100%; border: none; }}
        .controls {{ display: flex; gap: 16px; margin-top: 20px; align-items: center; }}
        .controls button {{ padding: 12px 24px; background: #3b82f6; color: white; border: none; border-radius: 8px; cursor: pointer; font-size: 16px; }}
        .controls button:hover {{ background: #2563eb; }}
        .page-indicator {{ color: #9ca3af; font-size: 16px; }}
    </style>
</head>
<body>
    <div class="container">
        <div class="sidebar"><h3>幻灯片</h3>{thumbs}</div>
        <div class="main">
            <div class="frame"><iframe id="frame" src="page-1.html"></iframe></div>
            <div class="controls">
                <button onclick="prevSlide()">上一页</button>
                <span class="page-indicator"><span id="current">1</span> / {page_count}</span>
                <button onclick="nextSlide()">下一页</button>
                <button onclick="toggleFullscreen()">全屏</button>
            </div>
        </div>
    </div>
    <script>
        let current = 1, total = {page_count};
        function update() {{ document.getElementById('frame').src = 'page-' + current + '.html'; document.getElementById('current').textContent = current; document.querySelectorAll('.slide-thumb').forEach((t,i) => t.classList.toggle('active', i+1 === current)); }}
        function goToSlide(n) {{ current = n; update(); }}
        function nextSlide() {{ if (current < total) {{ current++; update(); }} }}
        function prevSlide() {{ if (current > 1) {{ current--; update(); }} }}
        function toggleFullscreen() {{ if (!document.fullscreenElement) {{ document.documentElement.requestFullscreen(); }} else {{ document.exitFullscreen(); }} }}
        document.addEventListener('keydown', e => {{ if (e.key === 'ArrowRight' || e.key === ' ') nextSlide(); if (e.key === 'ArrowLeft') prevSlide(); if (e.key === 'Escape') document.exitFullscreen(); }});
    </script>
</body>
</html>'''


def generate_page_html(page, style, page_num, total_pages):
    bg_color = style["palette"][0]
    primary_color = style["palette"][1] if len(style["palette"]) > 1 else "#1F2937"
    points = "".join(
        f'<li style="margin-bottom:16px;display:flex;align-items:center;"><span style="display:inline-flex;width:32px;height:32px;background:{primary_color};color:white;border-radius:50%;align-items:center;justify-content:center;margin-right:16px;font-weight:bold;font-size:14px;">{i+1}</span><span style="font-size:20px;">{pt}</span></li>'
        for i, pt in enumerate(page["keyPoints"])
    )
    return f'''<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{page['title']}</title>
    <style>
        * {{ margin: 0; padding: 0; box-sizing: border-box; }}
        body {{ font-family: system-ui, sans-serif; background: {bg_color}; min-height: 100vh; display: flex; align-items: center; justify-content: center; }}
        .page {{ width: 100%; max-width: 1200px; padding: 60px; }}
        h1 {{ font-size: 48px; font-weight: 700; color: {primary_color}; margin-bottom: 40px; }}
        ul {{ list-style: none; font-size: 24px; line-height: 1.8; }}
    </style>
</head>
<body>
    <div class="page">
        <h1>{page['title']}</h1>
        <ul>{points}</ul>
    </div>
</body>
</html>'''


def sanitize_filename(name):
    invalid = '/\\:*?"<>|'
    for c in invalid:
        name = name.replace(c, '-')
    return name[:50]


def generate_ppt(topic, style_id, page_count):
    style = next((s for s in STYLES if s["id"] == style_id), STYLES[0])

    # 创建输出目录
    output_dir = os.path.join(os.path.expanduser("~"), "godjian-ppt-output", sanitize_filename(topic))
    os.makedirs(output_dir, exist_ok=True)

    # 生成页面结构
    pages = [{"title": topic, "keyPoints": ["演讲者", "日期", "版本号"]}]
    if page_count > 2:
        pages.insert(1, {"title": "目录", "keyPoints": ["内容概览", "重点章节"]})

    remaining = page_count - len(pages)
    for i in range(remaining):
        if len(pages) < page_count:
            pages.append({"title": f"第 {len(pages)} 部分", "keyPoints": ["要点 1", "要点 2", "要点 3"]})

    if len(pages) < page_count:
        pages.append({"title": "总结", "keyPoints": ["核心观点", "下一步"]})

    # 生成 index.html
    index_html = generate_html_index(topic, len(pages), style)
    with open(os.path.join(output_dir, "index.html"), "w", encoding="utf-8") as f:
        f.write(index_html)

    # 生成每一页
    for i, page in enumerate(pages):
        page_html = generate_page_html(page, style, i + 1, len(pages))
        with open(os.path.join(output_dir, f"page-{i+1}.html"), "w", encoding="utf-8") as f:
            f.write(page_html)

    return output_dir


def main():
    # 简单的 GUI
    root = tk.Tk()
    root.title("Godjian PPT - 本地幻灯片生成工具")
    root.geometry("800x600")

    # 主题输入
    tk.Label(root, text="PPT 主题:").pack(pady=5)
    topic_entry = tk.Entry(root, width=50)
    topic_entry.pack(pady=5)

    # 风格选择
    tk.Label(root, text="选择风格:").pack(pady=5)
    style_var = tk.StringVar(value="minimal-white")
    style_combo = ttk.Combobox(root, textvariable=style_var, width=30)
    style_combo["values"] = [f"{s['id']} - {s['name']}" for s in STYLES]
    style_combo.pack(pady=5)

    # 页数
    tk.Label(root, text="页数:").pack(pady=5)
    pages_entry = tk.Entry(root, width=10)
    pages_entry.insert(0, "8")
    pages_entry.pack(pady=5)

    # 输出目录
    output_label = tk.Label(root, text="")
    output_label.pack(pady=10)

    def on_generate():
        topic = topic_entry.get().strip()
        if not topic:
            tk.messagebox.showerror("错误", "请输入 PPT 主题")
            return

        style_id = style_var.get().split(" - ")[0]
        pages = int(pages_entry.get() or 8)

        try:
            output_dir = generate_ppt(topic, style_id, pages)
            output_label.config(text=f"已生成: {output_dir}")
            tk.messagebox.showinfo("成功", f"PPT 已生成！\n\n文件位置:\n{output_dir}\n\n请在浏览器中打开 index.html 预览")
        except Exception as e:
            tk.messagebox.showerror("错误", str(e))

    tk.Button(root, text="生成 PPT", command=on_generate, bg="#3B82F6", fg="white", padx=20, pady=10).pack(pady=20)

    root.mainloop()


if __name__ == "__main__":
    if TKINTER_AVAILABLE:
        main()
    else:
        print("Tkinter 不可用，请安装 Python Tkinter")
        print("在 Ubuntu/Debian: sudo apt-get install python3-tk")
