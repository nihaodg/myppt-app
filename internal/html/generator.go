package html

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/godjian/myppt-app/internal/ai"
	"github.com/godjian/myppt-app/internal/styles"
)

// Generator HTML 生成器
type Generator struct {
	outputDir string
	style     *styles.Style
}

// NewGenerator 创建生成器
func NewGenerator(outputDir string, style *styles.Style) *Generator {
	return &Generator{
		outputDir: outputDir,
		style:     style,
	}
}

// GenerateDeck 生成整份 PPT
func (g *Generator) GenerateDeck(plan *ai.PPTPlan, contract *ai.DesignContract, pageContents []string) error {
	// 创建输出目录
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 生成 index.html
	indexHTML := g.generateIndexHTML(plan, len(pageContents))
	if err := os.WriteFile(filepath.Join(g.outputDir, "index.html"), []byte(indexHTML), 0644); err != nil {
		return fmt.Errorf("写入 index.html 失败: %w", err)
	}

	// 生成每一页
	for i, content := range pageContents {
		pageFilename := fmt.Sprintf("page-%d.html", i+1)
		pagePath := filepath.Join(g.outputDir, pageFilename)
		
		pageHTML := g.wrapPageHTML(content, i+1, len(pageContents))
		if err := os.WriteFile(pagePath, []byte(pageHTML), 0644); err != nil {
			return fmt.Errorf("写入 %s 失败: %w", pageFilename, err)
		}
	}

	// 生成 CSS 文件
	css := g.generateCSS(contract)
	if err := os.WriteFile(filepath.Join(g.outputDir, "style.css"), []byte(css), 0644); err != nil {
		return fmt.Errorf("写入 style.css 失败: %w", err)
	}

	// 复制静态资源
	if err := g.copyAssets(); err != nil {
		return fmt.Errorf("复制静态资源失败: %w", err)
	}

	return nil
}

// generateIndexHTML 生成主页面
func (g *Generator) generateIndexHTML(plan *ai.PPTPlan, pageCount int) string {
	// 生成页面缩略图列表
	thumbs := strings.Builder{}
	for i := 1; i <= pageCount; i++ {
		thumbs.WriteString(fmt.Sprintf(`<div class="slide-thumb" onclick="goToSlide(%d)">%d</div>`, i, i))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - PPT</title>
    <link rel="stylesheet" href="style.css">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: system-ui, -apple-system, sans-serif;
            background: #1a1a2e;
            color: #fff;
            overflow: hidden;
        }
        .container {
            display: flex;
            height: 100vh;
        }
        .main-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .slide-frame {
            width: 100%%;
            max-width: 1200px;
            aspect-ratio: 16/9;
            background: #000;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 20px 60px rgba(0,0,0,0.5);
        }
        .slide-frame iframe {
            width: 100%%;
            height: 100%%;
            border: none;
        }
        .controls {
            display: flex;
            gap: 16px;
            margin-top: 20px;
            align-items: center;
        }
        .controls button {
            padding: 12px 24px;
            background: #3b82f6;
            color: white;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            font-size: 16px;
        }
        .controls button:hover { background: #2563eb; }
        .page-indicator {
            color: #9ca3af;
            font-size: 16px;
        }
        .sidebar {
            width: 200px;
            background: #16213e;
            padding: 20px;
            overflow-y: auto;
        }
        .sidebar h3 {
            margin-bottom: 16px;
            color: #9ca3af;
            font-size: 14px;
            text-transform: uppercase;
        }
        .slide-thumb {
            width: 100%%;
            aspect-ratio: 16/9;
            background: #0f3460;
            border-radius: 4px;
            margin-bottom: 8px;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            font-size: 18px;
            font-weight: bold;
            transition: all 0.2s;
        }
        .slide-thumb:hover { background: #1e5f74; transform: scale(1.05); }
        .slide-thumb.active { background: #3b82f6; }
        .title-bar {
            position: absolute;
            top: 20px;
            left: 220px;
            right: 20px;
            color: #9ca3af;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="title-bar">%s</div>
    <div class="container">
        <div class="sidebar">
            <h3>幻灯片</h3>
            %s
        </div>
        <div class="main-content">
            <div class="slide-frame">
                <iframe id="slide-frame" src="page-1.html"></iframe>
            </div>
            <div class="controls">
                <button onclick="prevSlide()">上一页</button>
                <span class="page-indicator"><span id="current-page">1</span> / %d</span>
                <button onclick="nextSlide()">下一页</button>
                <button onclick="presentMode()">全屏演示</button>
            </div>
        </div>
    </div>
    <script>
        let currentSlide = 1;
        const totalSlides = %d;

        function updateFrame() {
            document.getElementById('slide-frame').src = 'page-' + currentSlide + '.html';
            document.getElementById('current-page').textContent = currentSlide;
            
            document.querySelectorAll('.slide-thumb').forEach((thumb, i) => {
                thumb.classList.toggle('active', i + 1 === currentSlide);
            });
        }

        function goToSlide(n) {
            currentSlide = n;
            updateFrame();
        }

        function nextSlide() {
            if (currentSlide < totalSlides) {
                currentSlide++;
                updateFrame();
            }
        }

        function prevSlide() {
            if (currentSlide > 1) {
                currentSlide--;
                updateFrame();
            }
        }

        function presentMode() {
            document.getElementById('slide-frame').requestFullscreen();
        }

        document.addEventListener('keydown', (e) => {
            if (e.key === 'ArrowRight' || e.key === ' ') nextSlide();
            if (e.key === 'ArrowLeft') prevSlide();
            if (e.key === 'Escape') document.exitFullscreen();
            if (e.key === 'F11') {
                e.preventDefault();
                presentMode();
            }
        });
    </script>
</body>
</html>`, plan.Title, plan.Title, thumbs.String(), pageCount, pageCount)
}

// wrapPageHTML 包装页面内容
func (g *Generator) wrapPageHTML(content string, pageNum, totalPages int) string {
	// 清理内容
	content = strings.TrimSpace(content)
	
	// 确保内容有 HTML 结构
	if !strings.Contains(content, "<") {
		content = fmt.Sprintf(`<div class="content-wrapper">
    <h1 class="page-title">Page %d</h1>
    <div class="content-body">
        <p>%s</p>
    </div>
</div>`, pageNum, content)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Page %d</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <link rel="stylesheet" href="style.css">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        html, body {
            width: 100%%;
            height: 100%%;
            font-family: system-ui, -apple-system, sans-serif;
        }
        body {
            background: %s;
            color: %s;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .ppt-page-root {
            width: 100%%;
            height: 100%%;
            max-width: 1600px;
            max-height: 900px;
            margin: 0 auto;
        }
        @media (max-width: 1600px) {
            .ppt-page-root { max-width: 100%%; }
        }
        @media (max-height: 900px) {
            .ppt-page-root { max-height: 100%%; }
        }
    </style>
</head>
<body>
    <div class="ppt-page-root p-8">
        %s
    </div>
    <script src="assets/ppt-runtime.js"></script>
    <script>
        // 页面初始化动画
        document.addEventListener('DOMContentLoaded', () => {
            if (typeof PPT !== 'undefined') {
                PPT.init();
            }
        });
    </script>
</body>
</html>`, pageNum, pageNum, 
		getBackgroundColor(g.style),
		getTextColor(g.style),
		content)
}

// generateCSS 生成 CSS
func (g *Generator) generateCSS(contract *ai.DesignContract) string {
	bgColor := "#FFFFFF"
	if len(contract.Palette) > 0 {
		bgColor = contract.Palette[0]
	}
	
	return fmt.Sprintf(`/* Oh My PPT - Generated Styles */
:root {
    --primary: %s;
    --secondary: %s;
    --accent: %s;
    --text: %s;
    --background: %s;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: system-ui, -apple-system, 'Segoe UI', sans-serif;
    background: var(--background);
    color: var(--text);
    line-height: 1.6;
}

/* 页面容器 */
.ppt-page-root {
    width: 100%%;
    height: 100%%;
    position: relative;
}

/* 标题样式 */
.page-title, h1 {
    font-size: 3rem;
    font-weight: 700;
    margin-bottom: 1.5rem;
    color: var(--primary);
}

/* 副标题 */
.subtitle, h2 {
    font-size: 2rem;
    font-weight: 600;
    margin-bottom: 1rem;
    color: var(--secondary);
}

/* 内容区域 */
.content-wrapper {
    max-width: 100%%;
}

.content-body {
    font-size: 1.25rem;
    line-height: 1.8;
}

/* 卡片 */
.card {
    background: white;
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

/* 网格布局 */
.grid-layout {
    display: grid;
    gap: 1.5rem;
}

.grid-2 { grid-template-columns: repeat(2, 1fr); }
.grid-3 { grid-template-columns: repeat(3, 1fr); }
.grid-4 { grid-template-columns: repeat(4, 1fr); }

/* 列表 */
ul, ol {
    padding-left: 1.5rem;
}

li {
    margin-bottom: 0.75rem;
}

/* 图表容器 */
.chart-container {
    position: relative;
    height: 300px;
    width: 100%%;
}

/* 页面指示器 */
.page-indicator {
    position: absolute;
    bottom: 1rem;
    right: 1rem;
    font-size: 0.875rem;
    color: rgba(0, 0, 0, 0.5);
}
`, 
		getPrimaryColor(g.style),
		getSecondaryColor(g.style),
		getAccentColor(g.style),
		getTextColor(g.style),
		bgColor)
}

// copyAssets 复制静态资源
func (g *Generator) copyAssets() error {
	assetsDir := filepath.Join(g.outputDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return err
	}

	// 写入运行时脚本
	runtimeJS := getPPRuntimeJS()
	if err := os.WriteFile(filepath.Join(assetsDir, "ppt-runtime.js"), []byte(runtimeJS), 0644); err != nil {
		return err
	}

	return nil
}

// 辅助函数
func getBackgroundColor(s *styles.Style) string {
	if len(s.Palette) > 0 {
		return s.Palette[0]
	}
	return "#FFFFFF"
}

func getPrimaryColor(s *styles.Style) string {
	if len(s.Palette) > 1 {
		return s.Palette[1]
	}
	return "#1F2937"
}

func getSecondaryColor(s *styles.Style) string {
	if len(s.Palette) > 2 {
		return s.Palette[2]
	}
	return "#3B82F6"
}

func getAccentColor(s *styles.Style) string {
	if len(s.Palette) > 3 {
		return s.Palette[3]
	}
	return "#10B981"
}

func getTextColor(s *styles.Style) string {
	// 深色风格用浅色文字
	darkKeywords := []string{"dark", "black", "night", "cyber"}
	for _, kw := range darkKeywords {
		if strings.Contains(strings.ToLower(s.ID), kw) {
			return "#FFFFFF"
		}
	}
	return "#1F2937"
}

// getPPRuntimeJS PPT 运行时脚本
func getPPRuntimeJS() string {
	return `// Oh My PPT Runtime v1.0
const PPT = {
    animations: [],
    charts: [],
    
    init() {
        console.log('PPT Runtime Initialized');
        this.setupKeyboardNavigation();
        this.setupAnimations();
    },
    
    setupKeyboardNavigation() {
        document.addEventListener('keydown', (e) => {
            if (e.key === 'ArrowRight' || e.key === ' ') {
                e.preventDefault();
            }
        });
    },
    
    setupAnimations() {
        // 自动发现并执行动画
        const animatedElements = document.querySelectorAll('[data-animate]');
        animatedElements.forEach(el => {
            const animation = el.dataset.animate;
            const delay = parseInt(el.dataset.delay || '0');
            setTimeout(() => {
                this.runAnimation(el, animation);
            }, delay);
        });
    },
    
    runAnimation(element, animationType) {
        const animations = {
            'fade-in': { opacity: [0, 1], duration: 500 },
            'slide-up': { transform: ['translateY(30px)', 'translateY(0)'], opacity: [0, 1], duration: 500 },
            'slide-down': { transform: ['translateY(-30px)', 'translateY(0)'], opacity: [0, 1], duration: 500 },
            'slide-left': { transform: ['translateX(30px)', 'translateX(0)'], opacity: [0, 1], duration: 500 },
            'slide-right': { transform: ['translateX(-30px)', 'translateX(0)'], opacity: [0, 1], duration: 500 },
            'scale-in': { transform: ['scale(0.8)', 'scale(1)'], opacity: [0, 1], duration: 400 },
        };
        
        const config = animations[animationType] || animations['fade-in'];
        element.style.transition = "all " + config.duration + "ms ease-out";
        
        if (config.opacity) element.style.opacity = config.opacity[1];
        if (config.transform) element.style.transform = config.transform[1];
    },
    
    animate(targets, params) {
        const elements = typeof targets === 'string' 
            ? document.querySelectorAll(targets) 
            : [targets];
        
        elements.forEach(el => {
            if (params.opacity !== undefined) el.style.opacity = params.opacity;
            if (params.transform !== undefined) el.style.transform = params.transform;
            if (params.duration) el.style.transition = "all " + params.duration + "ms ease-out";
        });
    },
    
    createTimeline(steps) {
        let delay = 0;
        steps.forEach(step => {
            setTimeout(() => {
                if (step.element) {
                    const el = typeof step.element === 'string' 
                        ? document.querySelector(step.element) 
                        : step.element;
                    if (el) this.runAnimation(el, step.animation || 'fade-in');
                }
            }, delay);
            delay += step.duration || 500;
        });
    },
    
    stagger(selector, interval = 100) {
        const elements = document.querySelectorAll(selector);
        elements.forEach((el, i) => {
            el.style.opacity = '0';
            setTimeout(() => {
                this.runAnimation(el, 'fade-in');
            }, i * interval);
        });
    },
    
    createChart(canvasOrSelector, config) {
        const canvas = typeof canvasOrSelector === 'string'
            ? document.querySelector(canvasOrSelector)
            : canvasOrSelector;
            
        if (!canvas) {
            console.error('Canvas not found');
            return null;
        }
        
        if (typeof Chart === 'undefined') {
            console.warn('Chart.js not loaded');
            return null;
        }
        
        const chart = new Chart(canvas, config);
        this.charts.push(chart);
        return chart;
    },
    
    updateChart(chart, newConfig) {
        if (chart) {
            chart.data = newConfig.data || chart.data;
            chart.options = newConfig.options || chart.options;
            chart.update();
        }
    },
    
    destroyChart(chart) {
        if (chart) {
            chart.destroy();
            this.charts = this.charts.filter(c => c !== chart);
        }
    },
    
    resizeCharts() {
        this.charts.forEach(chart => chart.resize());
    }
};

// 自动初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => PPT.init());
} else {
    PPT.init();
}
`
}
